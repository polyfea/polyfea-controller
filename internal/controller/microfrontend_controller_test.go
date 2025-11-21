/*
Copyright 2025.

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

package controller

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	v1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
)

const (
	MicroFrontendName      = "test-microfrontend"
	MicroFrontendNamespace = "default"
	MicroFrontendFinalizer = "polyfea.github.io/finalizer"
)

func setupMicroFrontend(name string, service *v1alpha1.ServiceReference, modulePath *string, proxy *bool, frontendClass *string, staticResources []v1alpha1.StaticResources, cacheOptions *v1alpha1.PWACache) *v1alpha1.MicroFrontend {
	return &v1alpha1.MicroFrontend{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "polyfea.github.io/v1alpha1",
			Kind:       "MicroFrontend",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: MicroFrontendNamespace,
		},
		Spec: v1alpha1.MicroFrontendSpec{
			Service:         service,
			Proxy:           proxy,
			CacheStrategy:   "none",
			ModulePath:      modulePath,
			StaticResources: staticResources,
			FrontendClass:   frontendClass,
			CacheOptions:    cacheOptions,
		},
	}
}

func ensureMicroFrontendDeleted(ctx context.Context, timeout time.Duration, interval time.Duration) {
	existingMicroFrontend := &v1alpha1.MicroFrontend{}
	err := k8sClient.Get(ctx, types.NamespacedName{Name: MicroFrontendName, Namespace: MicroFrontendNamespace}, existingMicroFrontend)
	if err == nil {
		// Delete the existing resource
		Expect(k8sClient.Delete(ctx, existingMicroFrontend)).To(Succeed())

		// Wait for the resource to be fully deleted
		Eventually(func() bool {
			return errors.IsNotFound(k8sClient.Get(ctx, types.NamespacedName{Name: MicroFrontendName, Namespace: MicroFrontendNamespace}, existingMicroFrontend))
		}, timeout, interval).Should(BeTrue())
	}
}

var _ = Describe("MicroFrontend Controller", func() {
	const (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When reconciling a resource", func() {
		Context("When creating a MicroFrontend", func() {
			It("Should add and remove finalizer", func() {
				By("Creating a new MicroFrontend")
				testCtx := context.Background()
				proxy := true

				ensureMicroFrontendDeleted(testCtx, timeout, interval)

				microFrontend := setupMicroFrontend(
					MicroFrontendName,
					&v1alpha1.ServiceReference{
						Name:      ptr("test-service"),
						Namespace: ptr("test-namespace"),
						Port:      ptr(int32(80)),
						Scheme:    ptr("http"),
					},
					ptr("module.jsm"),
					&proxy,
					ptr("test-microfrontendclass"),
					[]v1alpha1.StaticResources{{Path: "static", Kind: "script"}},
					nil,
				)
				Expect(k8sClient.Create(testCtx, microFrontend)).To(Succeed())

				microFrontendLookupKey := types.NamespacedName{Name: MicroFrontendName, Namespace: MicroFrontendNamespace}

				createdMicroFrontend := &v1alpha1.MicroFrontend{}

				Eventually(func() bool {
					return k8sClient.Get(testCtx, microFrontendLookupKey, createdMicroFrontend) == nil
				}, timeout, interval).Should(BeTrue())
				Expect(createdMicroFrontend.Spec.Service.Name).Should(Equal(ptr("test-service")))
				Expect(createdMicroFrontend.Spec.Service.Namespace).Should(Equal(ptr("test-namespace")))

				By("Checking the MicroFrontend has finalizer")
				// Wait for the finalizer to be added
				Eventually(func() []string {
					updatedMicroFrontend := &v1alpha1.MicroFrontend{}
					err := k8sClient.Get(testCtx, types.NamespacedName{Name: MicroFrontendName, Namespace: MicroFrontendNamespace}, updatedMicroFrontend)
					if err != nil {
						return nil
					}
					return updatedMicroFrontend.ObjectMeta.Finalizers
				}, timeout, interval).Should(ContainElement(MicroFrontendFinalizer))

				By("Deleting the MicroFrontend")
				Expect(k8sClient.Delete(testCtx, createdMicroFrontend)).Should(Succeed())
				Eventually(func() bool {
					return errors.IsNotFound(k8sClient.Get(testCtx, microFrontendLookupKey, createdMicroFrontend))
				}, timeout, interval).Should(BeTrue())
			})

			DescribeTable("Validation scenarios",
				func(service *v1alpha1.ServiceReference, modulePath *string, shouldSucceed bool) {
					testCtx := context.Background()
					proxy := true

					ensureMicroFrontendDeleted(testCtx, timeout, interval)

					microFrontend := setupMicroFrontend(
						MicroFrontendName,
						service,
						modulePath,
						&proxy,
						ptr("test-microfrontendclass"),
						[]v1alpha1.StaticResources{{Path: "static", Kind: "script"}},
						nil,
					)
					if shouldSucceed {
						Expect(k8sClient.Create(testCtx, microFrontend)).Should(Succeed())
					} else {
						Expect(k8sClient.Create(testCtx, microFrontend)).Should(Not(Succeed()))
					}
				},
				Entry("missing service", nil, ptr("module.jsm"), false),
				Entry("missing modulePath", &v1alpha1.ServiceReference{Name: ptr("test-service"), Namespace: ptr("test-namespace"), Port: ptr(int32(80)), Scheme: ptr("http")}, nil, false),
				Entry("valid MicroFrontend with in-cluster service", &v1alpha1.ServiceReference{Name: ptr("test-service"), Namespace: ptr("test-namespace"), Port: ptr(int32(80)), Scheme: ptr("http")}, ptr("module.jsm"), true),
				Entry("valid MicroFrontend with external URI", &v1alpha1.ServiceReference{URI: ptr("https://cdn.example.com")}, ptr("module.jsm"), true),
				Entry("valid MicroFrontend with HTTP external URI", &v1alpha1.ServiceReference{URI: ptr("http://external-service.com")}, ptr("module.jsm"), true),
				Entry("invalid: both name and URI specified", &v1alpha1.ServiceReference{Name: ptr("test-service"), URI: ptr("https://cdn.example.com")}, ptr("module.jsm"), false),
				Entry("invalid: neither name nor URI specified", &v1alpha1.ServiceReference{Namespace: ptr("test-namespace"), Port: ptr(int32(80))}, ptr("module.jsm"), false),
				Entry("invalid: empty name and empty URI", &v1alpha1.ServiceReference{Name: ptr(""), URI: ptr("")}, ptr("module.jsm"), false),
			)

			It("Should create with defaults if optional fields are not specified", func() {
				By("Creating a new MicroFrontend with only required fields")
				testCtx := context.Background()

				ensureMicroFrontendDeleted(testCtx, timeout, interval)

				microFrontend := setupMicroFrontend(
					MicroFrontendName,
					&v1alpha1.ServiceReference{
						Name:      ptr("test-service"),
						Namespace: ptr("test-namespace"),
						Port:      ptr(int32(80)),
						Scheme:    ptr("http"),
					},
					ptr("module.jsm"),
					nil,
					nil,
					nil,
					nil,
				)
				Expect(k8sClient.Create(testCtx, microFrontend)).Should(Succeed())

				microFrontendLookupKey := types.NamespacedName{Name: MicroFrontendName, Namespace: MicroFrontendNamespace}
				createdMicroFrontend := &v1alpha1.MicroFrontend{}

				Eventually(func() bool {
					return k8sClient.Get(testCtx, microFrontendLookupKey, createdMicroFrontend) == nil
				}, timeout, interval).Should(BeTrue())
				Expect(createdMicroFrontend.Spec.Service.Name).Should(Equal(ptr("test-service")))
				Expect(createdMicroFrontend.Spec.Service.Namespace).Should(Equal(ptr("test-namespace")))
				Expect(*createdMicroFrontend.Spec.Proxy).Should(BeTrue())
				Expect(createdMicroFrontend.Spec.CacheStrategy).Should(Equal("none"))
				Expect(createdMicroFrontend.Spec.FrontendClass).Should(Equal(ptr("polyfea-controller-default")))
				Expect(createdMicroFrontend.Spec.CacheControl).Should(BeNil())
				Expect(createdMicroFrontend.Spec.StaticResources).Should(BeNil())
				Expect(createdMicroFrontend.Spec.DependsOn).Should(BeNil())

				Expect(k8sClient.Delete(testCtx, createdMicroFrontend)).Should(Succeed())
				Eventually(func() bool {
					return errors.IsNotFound(k8sClient.Get(testCtx, microFrontendLookupKey, createdMicroFrontend))
				}, timeout, interval).Should(BeTrue())
			})

			It("Should create MicroFrontend with external URI service", func() {
				By("Creating a MicroFrontend with external URI")
				testCtx := context.Background()
				proxy := false

				ensureMicroFrontendDeleted(testCtx, timeout, interval)

				microFrontend := setupMicroFrontend(
					MicroFrontendName,
					&v1alpha1.ServiceReference{
						URI: ptr("https://cdn.example.com"),
					},
					ptr("modules/app.js"),
					&proxy,
					ptr("test-microfrontendclass"),
					[]v1alpha1.StaticResources{{Path: "styles/main.css", Kind: "stylesheet"}},
					nil,
				)
				Expect(k8sClient.Create(testCtx, microFrontend)).To(Succeed())

				microFrontendLookupKey := types.NamespacedName{Name: MicroFrontendName, Namespace: MicroFrontendNamespace}
				createdMicroFrontend := &v1alpha1.MicroFrontend{}

				Eventually(func() bool {
					return k8sClient.Get(testCtx, microFrontendLookupKey, createdMicroFrontend) == nil
				}, timeout, interval).Should(BeTrue())

				By("Verifying the URI is set correctly")
				Expect(createdMicroFrontend.Spec.Service.URI).Should(Equal(ptr("https://cdn.example.com")))
				Expect(createdMicroFrontend.Spec.Service.Name).Should(BeNil())
				Expect(createdMicroFrontend.Spec.Service.Namespace).Should(BeNil())
				Expect(*createdMicroFrontend.Spec.Proxy).Should(BeFalse())

				Expect(k8sClient.Delete(testCtx, createdMicroFrontend)).Should(Succeed())
				Eventually(func() bool {
					return errors.IsNotFound(k8sClient.Get(testCtx, microFrontendLookupKey, createdMicroFrontend))
				}, timeout, interval).Should(BeTrue())
			})
		})
	})
})
