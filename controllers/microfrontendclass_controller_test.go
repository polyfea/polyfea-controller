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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	polyfeav1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

var _ = Describe("MicroFrontendClass controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		MicroFrontendClassName      = "test-microfrontendclass"
		MicroFrontendClassNamespace = "default"
		MicroFrontendClassFinalizer = "polyfea.github.io/finalizer"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating a MicroFrontendClass", func() {
		It("Should add and remove finallizer", func() {
			By("By creating a new MicroFrontendClass")
			ctx := context.Background()
			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"http://localhost:8080"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Succeed())

			microFrontendClassLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			createdMicroFrontendClass := &polyfeav1alpha1.MicroFrontendClass{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.BaseUri).Should(Equal(&[]string{"http://localhost:8080"}[0]))

			By("By checking the MicroFrontendClass has finalizer")
			Eventually(func() ([]string, error) {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				if err != nil {
					return []string{}, err
				}
				return createdMicroFrontendClass.ObjectMeta.GetFinalizers(), nil
			}, timeout, interval).Should(ContainElement(MicroFrontendClassFinalizer))

			By("By deleting the MicroFrontendClass")
			Expect(k8sClient.Delete(ctx, createdMicroFrontendClass)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())
		})

		It("Should not let create without base uri", func() {
			By("By creating a new MicroFrontendClass")
			ctx := context.Background()
			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Not(Succeed()))
		})

		It("Should not let create without title", func() {
			By("By creating a new MicroFrontendClass")
			ctx := context.Background()
			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					BaseUri:   &[]string{"http://localhost:8080"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Not(Succeed()))
		})

		It("Should let create with base uri and title and fill defaults", func() {
			By("By creating a new MicroFrontendClass")
			ctx := context.Background()
			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:   &[]string{"Test MicroFrontendClass"}[0],
					BaseUri: &[]string{"http://localhost:8080"}[0],
				},
			}

			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Succeed())

			microFrontendClassLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			createdMicroFrontendClass := &polyfeav1alpha1.MicroFrontendClass{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.BaseUri).Should(Equal(&[]string{"http://localhost:8080"}[0]))
			Expect(createdMicroFrontendClass.Spec.CspHeader).Should(Equal("default-src 'self'; font-src 'self'; script-src 'strict-dynamic' 'nonce-{NONCE_VALUE}'; worker-src 'self'; manifest-src 'self'; style-src 'self' 'strict-dynamic' 'nonce-{NONCE_VALUE}'; style-src-attr 'self' 'unsafe-inline';"))
			Expect(createdMicroFrontendClass.Spec.ExtraHeaders).Should(BeNil())
			Expect(createdMicroFrontendClass.Spec.UserRolesHeader).Should(Equal("x-auth-request-roles"))
			Expect(createdMicroFrontendClass.Spec.UserHeader).Should(Equal("x-auth-request-user"))

			Expect(k8sClient.Delete(ctx, createdMicroFrontendClass)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())
		})

		It("Should add and remove microfrontendclass from repository", func() {
			By("By creating a new MicroFrontendClass")
			ctx := context.Background()
			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"http://localhost:8080"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Succeed())

			microFrontendClassLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			createdMicroFrontendClass := &polyfeav1alpha1.MicroFrontendClass{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.BaseUri).Should(Equal(&[]string{"http://localhost:8080"}[0]))

			Eventually(func() []*polyfeav1alpha1.MicroFrontendClass {
				result, _ := microFrontendClassRepository.GetItems(func(mf *polyfeav1alpha1.MicroFrontendClass) bool {
					println("Checking microfrontend " + mf.Name)
					return mf.Name == MicroFrontendClassName
				})
				return result
			}, timeout, interval).Should(HaveLen(1))

			By("By deleting the MicroFrontend")
			Expect(k8sClient.Delete(ctx, createdMicroFrontendClass)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			Eventually(func() []*polyfeav1alpha1.MicroFrontendClass {
				result, _ := microFrontendClassRepository.GetItems(func(mf *polyfeav1alpha1.MicroFrontendClass) bool {
					println("Checking microfrontend " + mf.Name)
					return mf.Name == MicroFrontendClassName
				})
				return result
			}, timeout, interval).Should(HaveLen(0))
		})

		It("Should add and remove http route when routing parent refs are specified", func() {
			By("By creating a new MicroFrontendClass with routing parent refs")
			ctx := context.Background()

			operatorService := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "polyfea-webserver",
					Namespace: MicroFrontendClassNamespace,
					Labels:    map[string]string{OperatorServiceSelectorName: OperatorServiceSelectorValue},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name: PortName,
							Port: 80,
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, operatorService)).Should(Succeed())

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					Routing: &polyfeav1alpha1.Routing{
						ParentRefs: []gatewayv1.ParentReference{
							{
								Name: "abcd",
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Succeed())

			microFrontendClassLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			createdMicroFrontendClass := &polyfeav1alpha1.MicroFrontendClass{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.BaseUri).Should(Equal(&[]string{"/someother"}[0]))

			httpRouteLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			httpRoute := &gatewayv1.HTTPRoute{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, httpRouteLookupKey, httpRoute)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(httpRoute.Spec.Rules[0].Matches[0].Path.Value).Should(Equal(microFrontendClass.Spec.BaseUri))
			Expect(httpRoute.OwnerReferences[0].Kind).Should(Equal("MicroFrontendClass"))
			Expect(httpRoute.Spec.CommonRouteSpec.ParentRefs).Should(Equal(microFrontendClass.Spec.Routing.ParentRefs))

			By("By deleting the MicroFrontend")
			Expect(k8sClient.Delete(ctx, createdMicroFrontendClass)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			Eventually(func() []*polyfeav1alpha1.MicroFrontendClass {
				result, _ := microFrontendClassRepository.GetItems(func(mf *polyfeav1alpha1.MicroFrontendClass) bool {
					println("Checking microfrontend " + mf.Name)
					return mf.Name == MicroFrontendClassName
				})
				return result
			}, timeout, interval).Should(HaveLen(0))

			Expect(k8sClient.Delete(ctx, httpRoute)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, httpRouteLookupKey, httpRoute)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())
			Expect(k8sClient.Delete(ctx, operatorService)).Should(Succeed())
		})

		It("Should add and remove ingres when routing ingress class is specified", func() {
			By("By creating a new MicroFrontendClass with routing ingress class")
			ctx := context.Background()

			operatorService := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "polyfea-webserver",
					Namespace: MicroFrontendClassNamespace,
					Labels:    map[string]string{OperatorServiceSelectorName: OperatorServiceSelectorValue},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name: PortName,
							Port: 80,
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, operatorService)).Should(Succeed())

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					Routing: &polyfeav1alpha1.Routing{
						IngressClassName: &[]string{"nginx"}[0],
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Succeed())

			microFrontendClassLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			createdMicroFrontendClass := &polyfeav1alpha1.MicroFrontendClass{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.BaseUri).Should(Equal(&[]string{"/someother"}[0]))

			ingressLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			ingress := &networkingv1.Ingress{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, ingressLookupKey, ingress)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(ingress.Spec.IngressClassName).Should(Equal(microFrontendClass.Spec.Routing.IngressClassName))
			Expect(ingress.OwnerReferences[0].Kind).Should(Equal("MicroFrontendClass"))
			Expect(ingress.Spec.Rules[0].HTTP.Paths[0].Path).Should(Equal(*microFrontendClass.Spec.BaseUri))

			By("By deleting the MicroFrontend")
			Expect(k8sClient.Delete(ctx, createdMicroFrontendClass)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			Eventually(func() []*polyfeav1alpha1.MicroFrontendClass {
				result, _ := microFrontendClassRepository.GetItems(func(mf *polyfeav1alpha1.MicroFrontendClass) bool {
					println("Checking microfrontend " + mf.Name)
					return mf.Name == MicroFrontendClassName
				})
				return result
			}, timeout, interval).Should(HaveLen(0))

			Expect(k8sClient.Delete(ctx, ingress)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, ingressLookupKey, ingress)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())
			Expect(k8sClient.Delete(ctx, operatorService)).Should(Succeed())
		})

		It("Should not allow creating both ingress and httproute", func() {
			By("By creating a new MicroFrontendClass with routing ingress and parent refs")
			ctx := context.Background()

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					Routing: &polyfeav1alpha1.Routing{
						IngressClassName: &[]string{"nginx"}[0],
						ParentRefs: []gatewayv1.ParentReference{
							{
								Name: "abcd",
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Not(Succeed()))
		})

		It("Should remove HttpRoute if routing is not specified", func() {
			By("By creating a new MicroFrontendClass with routing parent refs")
			ctx := context.Background()

			operatorService := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "polyfea-webserver",
					Namespace: MicroFrontendClassNamespace,
					Labels:    map[string]string{OperatorServiceSelectorName: OperatorServiceSelectorValue},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name: PortName,
							Port: 80,
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, operatorService)).Should(Succeed())

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					Routing: &polyfeav1alpha1.Routing{
						ParentRefs: []gatewayv1.ParentReference{
							{
								Name: "abcd",
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Succeed())

			microFrontendClassLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			createdMicroFrontendClass := &polyfeav1alpha1.MicroFrontendClass{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.BaseUri).Should(Equal(&[]string{"/someother"}[0]))

			httpRouteLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			httpRoute := &gatewayv1.HTTPRoute{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, httpRouteLookupKey, httpRoute)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(httpRoute.Spec.Rules[0].Matches[0].Path.Value).Should(Equal(microFrontendClass.Spec.BaseUri))
			Expect(httpRoute.OwnerReferences[0].Kind).Should(Equal("MicroFrontendClass"))
			Expect(httpRoute.Spec.CommonRouteSpec.ParentRefs).Should(Equal(microFrontendClass.Spec.Routing.ParentRefs))

			By("By setting the routing to nil")
			patch := client.RawPatch(types.MergePatchType, []byte(`{"spec":{"routing":null}}`))
			Expect(k8sClient.Patch(ctx, microFrontendClass, patch)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.Routing).Should(BeNil())

			Eventually(func() bool {
				err := k8sClient.Get(ctx, httpRouteLookupKey, httpRoute)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			By("By deleting the MicroFrontend")
			Expect(k8sClient.Delete(ctx, createdMicroFrontendClass)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			Eventually(func() []*polyfeav1alpha1.MicroFrontendClass {
				result, _ := microFrontendClassRepository.GetItems(func(mf *polyfeav1alpha1.MicroFrontendClass) bool {
					println("Checking microfrontend " + mf.Name)
					return mf.Name == MicroFrontendClassName
				})
				return result
			}, timeout, interval).Should(HaveLen(0))

			Expect(k8sClient.Delete(ctx, operatorService)).Should(Succeed())
		})

		It("Should remove ingress if routing is not specified", func() {
			By("By creating a new MicroFrontendClass with routing ingress class")
			ctx := context.Background()

			operatorService := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "polyfea-webserver",
					Namespace: MicroFrontendClassNamespace,
					Labels:    map[string]string{OperatorServiceSelectorName: OperatorServiceSelectorValue},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name: PortName,
							Port: 80,
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, operatorService)).Should(Succeed())

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					Routing: &polyfeav1alpha1.Routing{
						IngressClassName: &[]string{"nginx"}[0],
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Succeed())

			microFrontendClassLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			createdMicroFrontendClass := &polyfeav1alpha1.MicroFrontendClass{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.BaseUri).Should(Equal(&[]string{"/someother"}[0]))

			ingressLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			ingress := &networkingv1.Ingress{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, ingressLookupKey, ingress)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(ingress.Spec.IngressClassName).Should(Equal(microFrontendClass.Spec.Routing.IngressClassName))
			Expect(ingress.OwnerReferences[0].Kind).Should(Equal("MicroFrontendClass"))
			Expect(ingress.Spec.Rules[0].HTTP.Paths[0].Path).Should(Equal(*microFrontendClass.Spec.BaseUri))

			By("By setting the routing to nil")
			patch := client.RawPatch(types.MergePatchType, []byte(`{"spec":{"routing":null}}`))
			Expect(k8sClient.Patch(ctx, microFrontendClass, patch)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.Routing).Should(BeNil())

			Eventually(func() bool {
				err := k8sClient.Get(ctx, ingressLookupKey, ingress)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			By("By deleting the MicroFrontend")
			Expect(k8sClient.Delete(ctx, createdMicroFrontendClass)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			Eventually(func() []*polyfeav1alpha1.MicroFrontendClass {
				result, _ := microFrontendClassRepository.GetItems(func(mf *polyfeav1alpha1.MicroFrontendClass) bool {
					println("Checking microfrontend " + mf.Name)
					return mf.Name == MicroFrontendClassName
				})
				return result
			}, timeout, interval).Should(HaveLen(0))

			Expect(k8sClient.Delete(ctx, operatorService)).Should(Succeed())
		})

		It("Should switch from route to ingress if routing is changed", func() {
			By("By creating a new MicroFrontendClass with routing parent refs")
			ctx := context.Background()

			operatorService := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "polyfea-webserver",
					Namespace: MicroFrontendClassNamespace,
					Labels:    map[string]string{OperatorServiceSelectorName: OperatorServiceSelectorValue},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name: PortName,
							Port: 80,
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, operatorService)).Should(Succeed())

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					Routing: &polyfeav1alpha1.Routing{
						ParentRefs: []gatewayv1.ParentReference{
							{
								Name: "abcd",
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Succeed())

			microFrontendClassLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			createdMicroFrontendClass := &polyfeav1alpha1.MicroFrontendClass{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.BaseUri).Should(Equal(&[]string{"/someother"}[0]))

			httpRouteLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			httpRoute := &gatewayv1.HTTPRoute{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, httpRouteLookupKey, httpRoute)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(httpRoute.Spec.Rules[0].Matches[0].Path.Value).Should(Equal(microFrontendClass.Spec.BaseUri))
			Expect(httpRoute.OwnerReferences[0].Kind).Should(Equal("MicroFrontendClass"))
			Expect(httpRoute.Spec.CommonRouteSpec.ParentRefs).Should(Equal(microFrontendClass.Spec.Routing.ParentRefs))

			By("By changing the routing to ingress")
			patch := client.RawPatch(types.MergePatchType, []byte(`{"spec":{"routing":{"parentRefs":null, "ingressClassName":"nginx"}}}`))
			Expect(k8sClient.Patch(ctx, microFrontendClass, patch)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.Routing.ParentRefs).Should(BeNil())

			Eventually(func() bool {
				err := k8sClient.Get(ctx, httpRouteLookupKey, httpRoute)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			ingressLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			ingress := &networkingv1.Ingress{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, ingressLookupKey, ingress)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			By("By deleting the MicroFrontend")
			Expect(k8sClient.Delete(ctx, createdMicroFrontendClass)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			Eventually(func() []*polyfeav1alpha1.MicroFrontendClass {
				result, _ := microFrontendClassRepository.GetItems(func(mf *polyfeav1alpha1.MicroFrontendClass) bool {
					println("Checking microfrontend " + mf.Name)
					return mf.Name == MicroFrontendClassName
				})
				return result
			}, timeout, interval).Should(HaveLen(0))

			Expect(k8sClient.Delete(ctx, operatorService)).Should(Succeed())
		})

		It("Should switch from ingress to route if routing is changed", func() {
			By("By creating a new MicroFrontendClass with routing ingress class")
			ctx := context.Background()

			operatorService := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "polyfea-webserver",
					Namespace: MicroFrontendClassNamespace,
					Labels:    map[string]string{OperatorServiceSelectorName: OperatorServiceSelectorValue},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name: PortName,
							Port: 80,
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, operatorService)).Should(Succeed())

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					Routing: &polyfeav1alpha1.Routing{
						IngressClassName: &[]string{"nginx"}[0],
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Succeed())

			microFrontendClassLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			createdMicroFrontendClass := &polyfeav1alpha1.MicroFrontendClass{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.BaseUri).Should(Equal(&[]string{"/someother"}[0]))

			ingressLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			ingress := &networkingv1.Ingress{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, ingressLookupKey, ingress)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(ingress.Spec.IngressClassName).Should(Equal(microFrontendClass.Spec.Routing.IngressClassName))
			Expect(ingress.OwnerReferences[0].Kind).Should(Equal("MicroFrontendClass"))
			Expect(ingress.Spec.Rules[0].HTTP.Paths[0].Path).Should(Equal(*microFrontendClass.Spec.BaseUri))

			By("By changing the routing to http route")
			patch := client.RawPatch(types.MergePatchType, []byte(`{"spec":{"routing":{"parentRefs":[{"name": "abcd"}], "ingressClassName":null}}}`))
			Expect(k8sClient.Patch(ctx, microFrontendClass, patch)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.Routing.IngressClassName).Should(BeNil())

			Eventually(func() bool {
				err := k8sClient.Get(ctx, ingressLookupKey, ingress)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			httpRouteLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			httpRoute := &gatewayv1.HTTPRoute{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, httpRouteLookupKey, httpRoute)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			By("By deleting the MicroFrontend")
			Expect(k8sClient.Delete(ctx, createdMicroFrontendClass)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			Eventually(func() []*polyfeav1alpha1.MicroFrontendClass {
				result, _ := microFrontendClassRepository.GetItems(func(mf *polyfeav1alpha1.MicroFrontendClass) bool {
					println("Checking microfrontend " + mf.Name)
					return mf.Name == MicroFrontendClassName
				})
				return result
			}, timeout, interval).Should(HaveLen(0))

			Expect(k8sClient.Delete(ctx, operatorService)).Should(Succeed())
		})

		It("Should apply changes for routing routing parent refs are changed", func() {
			By("By creating a new MicroFrontendClass with routing parent refs")
			ctx := context.Background()

			operatorService := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "polyfea-webserver",
					Namespace: MicroFrontendClassNamespace,
					Labels:    map[string]string{OperatorServiceSelectorName: OperatorServiceSelectorValue},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name: PortName,
							Port: 80,
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, operatorService)).Should(Succeed())

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					Routing: &polyfeav1alpha1.Routing{
						ParentRefs: []gatewayv1.ParentReference{
							{
								Name: "abcd",
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Succeed())

			microFrontendClassLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			createdMicroFrontendClass := &polyfeav1alpha1.MicroFrontendClass{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.BaseUri).Should(Equal(&[]string{"/someother"}[0]))

			httpRouteLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			httpRoute := &gatewayv1.HTTPRoute{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, httpRouteLookupKey, httpRoute)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(httpRoute.Spec.Rules[0].Matches[0].Path.Value).Should(Equal(microFrontendClass.Spec.BaseUri))
			Expect(httpRoute.OwnerReferences[0].Kind).Should(Equal("MicroFrontendClass"))
			Expect(httpRoute.Spec.CommonRouteSpec.ParentRefs).Should(Equal(microFrontendClass.Spec.Routing.ParentRefs))

			By("By changing the routing parent refs")
			patch := client.RawPatch(types.MergePatchType, []byte(`{"spec":{"routing":{"parentRefs":[{"name": "dcba"}], "ingressClassName":null}}}`))
			Expect(k8sClient.Patch(ctx, microFrontendClass, patch)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			Eventually(func() gatewayv1.ObjectName {
				err := k8sClient.Get(ctx, httpRouteLookupKey, httpRoute)

				if err != nil {
					return ""
				}

				return httpRoute.Spec.CommonRouteSpec.ParentRefs[0].Name

			}, timeout, interval).Should(Equal(gatewayv1.ObjectName("dcba")))

			By("By deleting the MicroFrontend")
			Expect(k8sClient.Delete(ctx, createdMicroFrontendClass)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			Eventually(func() []*polyfeav1alpha1.MicroFrontendClass {
				result, _ := microFrontendClassRepository.GetItems(func(mf *polyfeav1alpha1.MicroFrontendClass) bool {
					println("Checking microfrontend " + mf.Name)
					return mf.Name == MicroFrontendClassName
				})
				return result
			}, timeout, interval).Should(HaveLen(0))

			Expect(k8sClient.Delete(ctx, operatorService)).Should(Succeed())
		})

		It("Should apply changes for routing routing ingressclassname is changed", func() {
			By("By creating a new MicroFrontendClass with routing ingressclassname")
			ctx := context.Background()

			operatorService := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "polyfea-webserver",
					Namespace: MicroFrontendClassNamespace,
					Labels:    map[string]string{OperatorServiceSelectorName: OperatorServiceSelectorValue},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name: PortName,
							Port: 80,
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, operatorService)).Should(Succeed())

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					Routing: &polyfeav1alpha1.Routing{
						IngressClassName: &[]string{"nginx"}[0],
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Succeed())

			microFrontendClassLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			createdMicroFrontendClass := &polyfeav1alpha1.MicroFrontendClass{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.BaseUri).Should(Equal(&[]string{"/someother"}[0]))

			ingressLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			ingress := &networkingv1.Ingress{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, ingressLookupKey, ingress)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(ingress.Spec.IngressClassName).Should(Equal(microFrontendClass.Spec.Routing.IngressClassName))
			Expect(ingress.OwnerReferences[0].Kind).Should(Equal("MicroFrontendClass"))
			Expect(ingress.Spec.Rules[0].HTTP.Paths[0].Path).Should(Equal(*microFrontendClass.Spec.BaseUri))

			By("By changing the routing ingressclassname")
			patch := client.RawPatch(types.MergePatchType, []byte(`{"spec":{"routing":{"ingressClassName":"notnginx"}}}`))
			Expect(k8sClient.Patch(ctx, microFrontendClass, patch)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			Eventually(func() *string {
				err := k8sClient.Get(ctx, ingressLookupKey, ingress)

				if err != nil {
					return nil
				}

				return ingress.Spec.IngressClassName

			}, timeout, interval).Should(Equal(&[]string{"notnginx"}[0]))

			By("By deleting the MicroFrontend")
			Expect(k8sClient.Delete(ctx, createdMicroFrontendClass)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			Eventually(func() []*polyfeav1alpha1.MicroFrontendClass {
				result, _ := microFrontendClassRepository.GetItems(func(mf *polyfeav1alpha1.MicroFrontendClass) bool {
					println("Checking microfrontend " + mf.Name)
					return mf.Name == MicroFrontendClassName
				})
				return result
			}, timeout, interval).Should(HaveLen(0))

			Expect(k8sClient.Delete(ctx, operatorService)).Should(Succeed())
		})

		It("Should apply changes for routing service port is changed", func() {
			By("By creating a new MicroFrontendClass with routing parent refs")
			ctx := context.Background()

			operatorService := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "polyfea-webserver",
					Namespace: MicroFrontendClassNamespace,
					Labels:    map[string]string{OperatorServiceSelectorName: OperatorServiceSelectorValue},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name: PortName,
							Port: 80,
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, operatorService)).Should(Succeed())

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					Routing: &polyfeav1alpha1.Routing{
						ParentRefs: []gatewayv1.ParentReference{
							{
								Name: "abcd",
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Succeed())

			microFrontendClassLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			createdMicroFrontendClass := &polyfeav1alpha1.MicroFrontendClass{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.BaseUri).Should(Equal(&[]string{"/someother"}[0]))

			httpRouteLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			httpRoute := &gatewayv1.HTTPRoute{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, httpRouteLookupKey, httpRoute)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(httpRoute.Spec.Rules[0].Matches[0].Path.Value).Should(Equal(microFrontendClass.Spec.BaseUri))
			Expect(httpRoute.OwnerReferences[0].Kind).Should(Equal("MicroFrontendClass"))
			Expect(httpRoute.Spec.CommonRouteSpec.ParentRefs).Should(Equal(microFrontendClass.Spec.Routing.ParentRefs))

			By("By changing the port on service")
			operatorServiceLookupName := types.NamespacedName{Name: "polyfea-webserver", Namespace: MicroFrontendClassNamespace}
			patch := client.RawPatch(types.MergePatchType, []byte(`{"spec":{"ports":[{"name":"webserver", "port": 8080}]}}`))
			Expect(k8sClient.Patch(ctx, operatorService, patch)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, operatorServiceLookupName, operatorService)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			expectedPort := gatewayv1.PortNumber(8080)
			Eventually(func() *gatewayv1.PortNumber {
				err := k8sClient.Get(ctx, httpRouteLookupKey, httpRoute)

				if err != nil {
					return nil
				}

				return httpRoute.Spec.Rules[0].BackendRefs[0].Port

			}, timeout, interval).Should(Equal(&expectedPort))

			By("By deleting the MicroFrontend")
			Expect(k8sClient.Delete(ctx, createdMicroFrontendClass)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			Eventually(func() []*polyfeav1alpha1.MicroFrontendClass {
				result, _ := microFrontendClassRepository.GetItems(func(mf *polyfeav1alpha1.MicroFrontendClass) bool {
					println("Checking microfrontend " + mf.Name)
					return mf.Name == MicroFrontendClassName
				})
				return result
			}, timeout, interval).Should(HaveLen(0))

			Expect(k8sClient.Delete(ctx, operatorService)).Should(Succeed())
		})

		It("Should create RefGrant", func() {
			By("By creating a new MicroFrontendClass with routing parent refs")
			ctx := context.Background()

			operatorService := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "polyfea-webserver",
					Namespace: MicroFrontendClassNamespace,
					Labels:    map[string]string{OperatorServiceSelectorName: OperatorServiceSelectorValue},
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name: PortName,
							Port: 80,
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, operatorService)).Should(Succeed())

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					Routing: &polyfeav1alpha1.Routing{
						ParentRefs: []gatewayv1.ParentReference{
							{
								Name: "abcd",
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Succeed())

			microFrontendClassLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			createdMicroFrontendClass := &polyfeav1alpha1.MicroFrontendClass{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.BaseUri).Should(Equal(&[]string{"/someother"}[0]))

			refGrantLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			refGrant := &gatewayv1alpha2.ReferenceGrant{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, refGrantLookupKey, refGrant)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(refGrant.Namespace).Should(Equal(operatorService.Namespace))

			By("By deleting the MicroFrontend")
			Expect(k8sClient.Delete(ctx, createdMicroFrontendClass)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			Eventually(func() []*polyfeav1alpha1.MicroFrontendClass {
				result, _ := microFrontendClassRepository.GetItems(func(mf *polyfeav1alpha1.MicroFrontendClass) bool {
					println("Checking microfrontend " + mf.Name)
					return mf.Name == MicroFrontendClassName
				})
				return result
			}, timeout, interval).Should(HaveLen(0))

			Expect(k8sClient.Delete(ctx, operatorService)).Should(Succeed())
		})

		It("Should create pwa with default PolyfeaSWReconcileInterval when not specified", func() {
			By("By creating a new MicroFrontendClass with PWA")
			ctx := context.Background()

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					ProgressiveWebApp: &polyfeav1alpha1.ProgressiveWebApp{
						WebAppManifest: &polyfeav1alpha1.WebAppManifest{
							Name: &[]string{"Test"}[0],
							Icons: []polyfeav1alpha1.PWAIcon{
								{
									Src:   &[]string{"icon.png"}[0],
									Sizes: &[]string{"192x192"}[0],
									Type:  &[]string{"image/png"}[0],
								},
							},
							StartUrl: &[]string{"/"}[0],
							Display:  &[]string{"standalone"}[0],
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Succeed())

			microFrontendClassLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			createdMicroFrontendClass := &polyfeav1alpha1.MicroFrontendClass{}
			defaultPolyfeaSWReconcileInterval := &[]int32{1800000}[0]

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.BaseUri).Should(Equal(&[]string{"/someother"}[0]))
			Expect(createdMicroFrontendClass.Spec.ProgressiveWebApp.PolyfeaSWReconcileInterval).Should(Equal(defaultPolyfeaSWReconcileInterval))

			By("By deleting the MicroFrontend")
			Expect(k8sClient.Delete(ctx, createdMicroFrontendClass)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())
		})

		It("Should not create pwa without WebAppManifest", func() {
			By("By creating a new MicroFrontendClass with PWA")
			ctx := context.Background()

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader:   "X-User-Roles",
					UserHeader:        "X-User-Id",
					ProgressiveWebApp: &polyfeav1alpha1.ProgressiveWebApp{},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Not(Succeed()))
		})

		It("Should not create pwa without mandatory fields in web manifest", func() {
			By("By creating a new MicroFrontendClass with manifest without name")
			ctx := context.Background()

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					ProgressiveWebApp: &polyfeav1alpha1.ProgressiveWebApp{
						WebAppManifest: &polyfeav1alpha1.WebAppManifest{
							Icons: []polyfeav1alpha1.PWAIcon{
								{
									Src:   &[]string{"icon.png"}[0],
									Sizes: &[]string{"192x192"}[0],
									Type:  &[]string{"image/png"}[0],
								},
							},
							StartUrl: &[]string{"/"}[0],
							Display:  &[]string{"standalone"}[0],
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Not(Succeed()))
		})

		It("Should not create pwa without mandatory fields in web manifest", func() {
			By("By creating a new MicroFrontendClass with manifest without Icons")
			ctx := context.Background()

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					ProgressiveWebApp: &polyfeav1alpha1.ProgressiveWebApp{
						WebAppManifest: &polyfeav1alpha1.WebAppManifest{
							Name:     &[]string{"Test"}[0],
							StartUrl: &[]string{"/"}[0],
							Display:  &[]string{"standalone"}[0],
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Not(Succeed()))
		})

		It("Should not create pwa without mandatory fields in web manifest", func() {
			By("By creating a new MicroFrontendClass with manifest without starturl")
			ctx := context.Background()

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					ProgressiveWebApp: &polyfeav1alpha1.ProgressiveWebApp{
						WebAppManifest: &polyfeav1alpha1.WebAppManifest{
							Name: &[]string{"Test"}[0],
							Icons: []polyfeav1alpha1.PWAIcon{
								{
									Src:   &[]string{"icon.png"}[0],
									Sizes: &[]string{"192x192"}[0],
									Type:  &[]string{"image/png"}[0],
								},
							},
							Display: &[]string{"standalone"}[0],
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Not(Succeed()))
		})

		It("Should not create pwa without mandatory fields in web manifest", func() {
			By("By creating a new MicroFrontendClass with manifest without display")
			ctx := context.Background()

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					ProgressiveWebApp: &polyfeav1alpha1.ProgressiveWebApp{
						WebAppManifest: &polyfeav1alpha1.WebAppManifest{
							Name: &[]string{"Test"}[0],
							Icons: []polyfeav1alpha1.PWAIcon{
								{
									Src:   &[]string{"icon.png"}[0],
									Sizes: &[]string{"192x192"}[0],
									Type:  &[]string{"image/png"}[0],
								},
							},
							StartUrl: &[]string{"/"}[0],
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Not(Succeed()))
		})

		It("Should not create pwa without mandatory fields in icon", func() {
			By("By creating a new MicroFrontendClass with PWA without mandatory sizes in icon")
			ctx := context.Background()

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					ProgressiveWebApp: &polyfeav1alpha1.ProgressiveWebApp{
						WebAppManifest: &polyfeav1alpha1.WebAppManifest{
							Name: &[]string{"Test"}[0],
							Icons: []polyfeav1alpha1.PWAIcon{
								{
									Src:  &[]string{"icon.png"}[0],
									Type: &[]string{"image/png"}[0],
								},
							},
							StartUrl: &[]string{"/"}[0],
							Display:  &[]string{"standalone"}[0],
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Not(Succeed()))
		})

		It("Should not create pwa without mandatory fields in icon", func() {
			By("By creating a new MicroFrontendClass with PWA without mandatory src in icon")
			ctx := context.Background()

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					ProgressiveWebApp: &polyfeav1alpha1.ProgressiveWebApp{
						WebAppManifest: &polyfeav1alpha1.WebAppManifest{
							Name: &[]string{"Test"}[0],
							Icons: []polyfeav1alpha1.PWAIcon{
								{
									Sizes: &[]string{"192x192"}[0],
									Type:  &[]string{"image/png"}[0],
								},
							},
							StartUrl: &[]string{"/"}[0],
							Display:  &[]string{"standalone"}[0],
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Not(Succeed()))
		})

		It("Should not create pwa without mandatory fields in icon", func() {
			By("By creating a new MicroFrontendClass with PWA without mandatory type in icon")
			ctx := context.Background()

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					ProgressiveWebApp: &polyfeav1alpha1.ProgressiveWebApp{
						WebAppManifest: &polyfeav1alpha1.WebAppManifest{
							Name: &[]string{"Test"}[0],
							Icons: []polyfeav1alpha1.PWAIcon{
								{
									Sizes: &[]string{"192x192"}[0],
									Src:   &[]string{"icon.png"}[0],
								},
							},
							StartUrl: &[]string{"/"}[0],
							Display:  &[]string{"standalone"}[0],
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Not(Succeed()))
		})

		It("Should not create pwa without mandatory fields in pre cache", func() {
			By("By creating a new MicroFrontendClass with PWA without mandatory url in pre cache")
			ctx := context.Background()

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					ProgressiveWebApp: &polyfeav1alpha1.ProgressiveWebApp{
						WebAppManifest: &polyfeav1alpha1.WebAppManifest{
							Name: &[]string{"Test"}[0],
							Icons: []polyfeav1alpha1.PWAIcon{
								{
									Type:  &[]string{"image/png"}[0],
									Sizes: &[]string{"192x192"}[0],
									Src:   &[]string{"icon.png"}[0],
								},
							},
							StartUrl: &[]string{"/"}[0],
							Display:  &[]string{"standalone"}[0],
						},
						CacheOptions: &polyfeav1alpha1.PWACache{
							PreCache: []polyfeav1alpha1.PreCacheEntry{
								{
									Revision: &[]string{"1"}[0],
								},
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Not(Succeed()))
		})

		It("Should should create pwa without optional fields in pre cache", func() {
			By("By creating a new MicroFrontendClass with PWA without optional revision in pre cache")
			ctx := context.Background()

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					ProgressiveWebApp: &polyfeav1alpha1.ProgressiveWebApp{
						WebAppManifest: &polyfeav1alpha1.WebAppManifest{
							Name: &[]string{"Test"}[0],
							Icons: []polyfeav1alpha1.PWAIcon{
								{
									Type:  &[]string{"image/png"}[0],
									Sizes: &[]string{"192x192"}[0],
									Src:   &[]string{"icon.png"}[0],
								},
							},
							StartUrl: &[]string{"/"}[0],
							Display:  &[]string{"standalone"}[0],
						},
						CacheOptions: &polyfeav1alpha1.PWACache{
							PreCache: []polyfeav1alpha1.PreCacheEntry{
								{
									URL: &[]string{"https://example.com"}[0],
								},
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Succeed())

			microFrontendClassLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			createdMicroFrontendClass := &polyfeav1alpha1.MicroFrontendClass{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.BaseUri).Should(Equal(&[]string{"/someother"}[0]))
			Expect(createdMicroFrontendClass.Spec.ProgressiveWebApp.CacheOptions.PreCache[0].Revision).Should(BeNil())

			By("By deleting the MicroFrontend")
			Expect(k8sClient.Delete(ctx, createdMicroFrontendClass)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())
		})

		It("Should should create pwa without optional fields in runtime cache", func() {
			By("By creating a new MicroFrontendClass with PWA without optional fields in runtime cache")
			ctx := context.Background()

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					ProgressiveWebApp: &polyfeav1alpha1.ProgressiveWebApp{
						WebAppManifest: &polyfeav1alpha1.WebAppManifest{
							Name: &[]string{"Test"}[0],
							Icons: []polyfeav1alpha1.PWAIcon{
								{
									Type:  &[]string{"image/png"}[0],
									Sizes: &[]string{"192x192"}[0],
									Src:   &[]string{"icon.png"}[0],
								},
							},
							StartUrl: &[]string{"/"}[0],
							Display:  &[]string{"standalone"}[0],
						},
						CacheOptions: &polyfeav1alpha1.PWACache{
							CacheRoutes: []polyfeav1alpha1.CacheRoute{
								{
									Pattern: &[]string{"/"}[0],
								},
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Succeed())

			microFrontendClassLookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
			createdMicroFrontendClass := &polyfeav1alpha1.MicroFrontendClass{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontendClass.Spec.BaseUri).Should(Equal(&[]string{"/someother"}[0]))
			Expect(createdMicroFrontendClass.Spec.ProgressiveWebApp.CacheOptions.CacheRoutes[0].Strategy).Should(Equal(&[]string{"cache-first"}[0]))
			Expect(createdMicroFrontendClass.Spec.ProgressiveWebApp.CacheOptions.CacheRoutes[0].Method).Should(Equal(&[]string{"GET"}[0]))
			Expect(createdMicroFrontendClass.Spec.ProgressiveWebApp.CacheOptions.CacheRoutes[0].Statuses).Should(Equal([]int32{0, 200, 201, 202, 204}))

			By("By deleting the MicroFrontend")
			Expect(k8sClient.Delete(ctx, createdMicroFrontendClass)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendClassLookupKey, createdMicroFrontendClass)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())
		})

		It("Should should not create pwa without mandatory fields in runtime cache", func() {
			By("By creating a new MicroFrontendClass with PWA without pattern field in runtime cache")
			ctx := context.Background()

			microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontendClass",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendClassName,
					Namespace: MicroFrontendClassNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendClassSpec{
					Title:     &[]string{"Test MicroFrontendClass"}[0],
					BaseUri:   &[]string{"/someother"}[0],
					CspHeader: "default-src 'self';",
					ExtraHeaders: []polyfeav1alpha1.Header{
						{
							Name:  "X-Frame-Options",
							Value: "DENY",
						},
					},
					UserRolesHeader: "X-User-Roles",
					UserHeader:      "X-User-Id",
					ProgressiveWebApp: &polyfeav1alpha1.ProgressiveWebApp{
						WebAppManifest: &polyfeav1alpha1.WebAppManifest{
							Name: &[]string{"Test"}[0],
							Icons: []polyfeav1alpha1.PWAIcon{
								{
									Type:  &[]string{"image/png"}[0],
									Sizes: &[]string{"192x192"}[0],
									Src:   &[]string{"icon.png"}[0],
								},
							},
							StartUrl: &[]string{"/"}[0],
							Display:  &[]string{"standalone"}[0],
						},
						CacheOptions: &polyfeav1alpha1.PWACache{
							CacheRoutes: []polyfeav1alpha1.CacheRoute{
								{
									Strategy: &[]string{"network-first"}[0],
								},
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontendClass)).Should(Not(Succeed()))
		})
	})
})
