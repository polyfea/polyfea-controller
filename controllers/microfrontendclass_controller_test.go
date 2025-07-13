/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors" // Correct import for `errors.IsNotFound`

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	polyfeav1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	MicroFrontendClassName      = "test-microfrontendclass"
	MicroFrontendClassNamespace = "default"
	MicroFrontendClassFinalizer = "polyfea.github.io/finalizer"
)

func setupMicroFrontendClass(title, baseUri *string, routing *polyfeav1alpha1.Routing) *polyfeav1alpha1.MicroFrontendClass {
	return &polyfeav1alpha1.MicroFrontendClass{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "polyfea.github.io/v1alpha1",
			Kind:       "MicroFrontendClass",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MicroFrontendClassName,
			Namespace: MicroFrontendClassNamespace,
		},
		Spec: polyfeav1alpha1.MicroFrontendClassSpec{
			Title:   title,
			BaseUri: baseUri,
			Routing: routing,
		},
	}
}

var _ = Describe("MicroFrontendClass controller", func() {
	const (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("Validation and defaults", func() {
		DescribeTable("Validation scenarios",
			func(title, baseUri *string, shouldSucceed bool) {
				ctx := context.Background()
				mfc := setupMicroFrontendClass(title, baseUri, nil)
				if shouldSucceed {
					Expect(k8sClient.Create(ctx, mfc)).Should(Succeed())
				} else {
					Expect(k8sClient.Create(ctx, mfc)).Should(Not(Succeed()))
				}
			},
			Entry("missing baseUri", ptr("Test"), nil, false),
			Entry("missing title", nil, ptr("/base"), false),
			Entry("valid MicroFrontendClass", ptr("Test"), ptr("/base"), true),
		)

		It("Should fill defaults when only required fields are set", func() {
			ctx := context.Background()
			mfc := setupMicroFrontendClass(ptr("Test"), ptr("/base"), nil)
			Expect(k8sClient.Create(ctx, mfc)).Should(Succeed())

			lookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			created := &polyfeav1alpha1.MicroFrontendClass{}
			Eventually(func() bool {
				return k8sClient.Get(ctx, lookupKey, created) == nil
			}, timeout, interval).Should(BeTrue())
			Expect(created.Spec.CspHeader).ShouldNot(BeEmpty())
			Expect(created.Spec.UserRolesHeader).ShouldNot(BeEmpty())
			Expect(created.Spec.UserHeader).ShouldNot(BeEmpty())

			Expect(k8sClient.Delete(ctx, created)).Should(Succeed())
			Eventually(func() bool {
				return errors.IsNotFound(k8sClient.Get(ctx, lookupKey, created))
			}, timeout, interval).Should(BeTrue())
		})
	})

	Context("Repository add/remove", func() {
		It("Should add and remove MicroFrontendClass from repository", func() {
			ctx := context.Background()
			mfc := setupMicroFrontendClass(ptr("Test"), ptr("/base"), nil)
			Expect(k8sClient.Create(ctx, mfc)).Should(Succeed())

			lookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			created := &polyfeav1alpha1.MicroFrontendClass{}
			Eventually(func() bool {
				return k8sClient.Get(ctx, lookupKey, created) == nil
			}, timeout, interval).Should(BeTrue())

			Eventually(func() []*polyfeav1alpha1.MicroFrontendClass {
				result, _ := microFrontendClassRepository.List(func(m *polyfeav1alpha1.MicroFrontendClass) bool {
					return m.Name == MicroFrontendClassName
				})
				return result
			}, timeout, interval).Should(HaveLen(1))

			Expect(k8sClient.Delete(ctx, created)).Should(Succeed())
			Eventually(func() []*polyfeav1alpha1.MicroFrontendClass {
				result, _ := microFrontendClassRepository.List(func(m *polyfeav1alpha1.MicroFrontendClass) bool {
					return m.Name == MicroFrontendClassName
				})
				return result
			}, timeout, interval).Should(HaveLen(0))
		})
	})

	Context("Routing resource creation and switching", func() {
		type routingCase struct {
			routing       *polyfeav1alpha1.Routing
			expectIngress bool
			expectRoute   bool
		}

		DescribeTable("Routing scenarios",
			func(rc routingCase) {
				ctx := context.Background()
				operatorService := &corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "polyfea-webserver",
						Namespace: MicroFrontendClassNamespace,
						Labels:    map[string]string{OperatorServiceSelectorName: OperatorServiceSelectorValue},
					},
					Spec: corev1.ServiceSpec{
						Ports: []corev1.ServicePort{{Name: PortName, Port: 80}},
					},
				}
				Expect(k8sClient.Create(ctx, operatorService)).Should(Succeed())

				mfc := setupMicroFrontendClass(ptr("Test"), ptr("/base"), rc.routing)
				Expect(k8sClient.Create(ctx, mfc)).Should(Succeed())

				lookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
				created := &polyfeav1alpha1.MicroFrontendClass{}
				Eventually(func() bool {
					return k8sClient.Get(ctx, lookupKey, created) == nil
				}, timeout, interval).Should(BeTrue())

				if rc.expectRoute {
					httpRoute := &gatewayv1.HTTPRoute{}
					Eventually(func() bool {
						return k8sClient.Get(ctx, lookupKey, httpRoute) == nil
					}, timeout, interval).Should(BeTrue())
					Expect(httpRoute.OwnerReferences[0].Kind).Should(Equal("MicroFrontendClass"))
				}
				if rc.expectIngress {
					ingress := &networkingv1.Ingress{}
					Eventually(func() bool {
						return k8sClient.Get(ctx, lookupKey, ingress) == nil
					}, timeout, interval).Should(BeTrue())
					Expect(ingress.OwnerReferences[0].Kind).Should(Equal("MicroFrontendClass"))
				}

				Expect(k8sClient.Delete(ctx, created)).Should(Succeed())
				Expect(k8sClient.Delete(ctx, operatorService)).Should(Succeed())
			},
			Entry("parentRefs only", routingCase{routing: &polyfeav1alpha1.Routing{ParentRefs: []gatewayv1.ParentReference{{Name: "abcd"}}}, expectRoute: true, expectIngress: false}),
			Entry("ingressClassName only", routingCase{routing: &polyfeav1alpha1.Routing{IngressClassName: ptr("nginx")}, expectRoute: false, expectIngress: true}),
		)
	})

	Context("Manifest validation", func() {
		DescribeTable("Manifest scenarios",
			func(manifest *polyfeav1alpha1.WebAppManifest, shouldSucceed bool) {
				ctx := context.Background()
				mfc := setupMicroFrontendClass(ptr("Test"), ptr("/base"), nil)
				mfc.Spec.ProgressiveWebApp = &polyfeav1alpha1.ProgressiveWebApp{WebAppManifest: manifest}
				if shouldSucceed {
					Expect(k8sClient.Create(ctx, mfc)).Should(Succeed())
				} else {
					Expect(k8sClient.Create(ctx, mfc)).Should(Not(Succeed()))
				}
			},
			Entry("missing name", &polyfeav1alpha1.WebAppManifest{
				Icons:    []polyfeav1alpha1.PWAIcon{{Src: ptr("icon.png"), Sizes: ptr("192x192"), Type: ptr("image/png")}},
				StartUrl: ptr("/"), Display: ptr("standalone"),
			}, false),
			Entry("missing icons", &polyfeav1alpha1.WebAppManifest{
				Name: ptr("Test"), StartUrl: ptr("/"), Display: ptr("standalone"),
			}, false),
			Entry("missing startUrl", &polyfeav1alpha1.WebAppManifest{
				Name: ptr("Test"), Icons: []polyfeav1alpha1.PWAIcon{{Src: ptr("icon.png"), Sizes: ptr("192x192"), Type: ptr("image/png")}},
				Display: ptr("standalone"),
			}, false),
			Entry("missing display", &polyfeav1alpha1.WebAppManifest{
				Name: ptr("Test"), Icons: []polyfeav1alpha1.PWAIcon{{Src: ptr("icon.png"), Sizes: ptr("192x192"), Type: ptr("image/png")}},
				StartUrl: ptr("/"),
			}, false),
			Entry("valid manifest", &polyfeav1alpha1.WebAppManifest{
				Name: ptr("Test"), Icons: []polyfeav1alpha1.PWAIcon{{Src: ptr("icon.png"), Sizes: ptr("192x192"), Type: ptr("image/png")}},
				StartUrl: ptr("/"), Display: ptr("standalone"),
			}, true),
		)
	})
})
