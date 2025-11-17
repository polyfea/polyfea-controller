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

package controller

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	polyfeav1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
)

const (
	MicroFrontendName      = "test-microfrontend"
	MicroFrontendNamespace = "default"
	MicroFrontendFinalizer = "polyfea.github.io/finalizer"
)

func setupMicroFrontend(name string, service, modulePath *string, proxy *bool, frontendClass *string, staticResources []polyfeav1alpha1.StaticResources, cacheOptions *polyfeav1alpha1.PWACache) *polyfeav1alpha1.MicroFrontend {
	return &polyfeav1alpha1.MicroFrontend{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "polyfea.github.io/v1alpha1",
			Kind:       "MicroFrontend",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: MicroFrontendNamespace,
		},
		Spec: polyfeav1alpha1.MicroFrontendSpec{
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
	existingMicroFrontend := &polyfeav1alpha1.MicroFrontend{}
	err := k8sClient.Get(ctx, types.NamespacedName{Name: MicroFrontendName, Namespace: MicroFrontendNamespace}, existingMicroFrontend)
	if err == nil {
		// Delete the existing resource
		Expect(k8sClient.Delete(ctx, existingMicroFrontend)).Should(Succeed())
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
		const resourceName = "test-microfrontend"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default",
		}
		microfrontend := &polyfeav1alpha1.MicroFrontend{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind MicroFrontend")
			err := k8sClient.Get(ctx, typeNamespacedName, microfrontend)
			if err != nil && errors.IsNotFound(err) {
				service := "http://test-service.test-namespace.svc.cluster.local"
				modulePath := "module.jsm"
				proxy := true
				resource := &polyfeav1alpha1.MicroFrontend{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: polyfeav1alpha1.MicroFrontendSpec{
						Service:    &service,
						ModulePath: &modulePath,
						Proxy:      &proxy,
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &polyfeav1alpha1.MicroFrontend{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance MicroFrontend")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})

		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &MicroFrontendReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		Context("When creating a MicroFrontend", func() {
			It("Should add and remove finalizer", func() {
				By("Creating a new MicroFrontend")
				testCtx := context.Background()
				proxy := true

				ensureMicroFrontendDeleted(testCtx, timeout, interval)

				microFrontend := setupMicroFrontend(
					MicroFrontendName,
					ptr("http://test-service.test-namespace.svc.cluster.local"),
					ptr("module.jsm"),
					&proxy,
					ptr("test-microfrontendclass"),
					[]polyfeav1alpha1.StaticResources{{Path: "static", Kind: "script"}},
					nil,
				)
				Expect(k8sClient.Create(testCtx, microFrontend)).Should(Succeed())

				microFrontendLookupKey := types.NamespacedName{Name: MicroFrontendName, Namespace: MicroFrontendNamespace}
				createdMicroFrontend := &polyfeav1alpha1.MicroFrontend{}

				Eventually(func() bool {
					return k8sClient.Get(testCtx, microFrontendLookupKey, createdMicroFrontend) == nil
				}, timeout, interval).Should(BeTrue())
				Expect(createdMicroFrontend.Spec.Service).Should(Equal(ptr("http://test-service.test-namespace.svc.cluster.local")))

				By("Checking the MicroFrontend has finalizer")
				// Wait for the finalizer to be added
				Eventually(func() []string {
					updatedMicroFrontend := &polyfeav1alpha1.MicroFrontend{}
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
				func(service, modulePath *string, shouldSucceed bool) {
					testCtx := context.Background()
					proxy := true

					ensureMicroFrontendDeleted(testCtx, timeout, interval)

					microFrontend := setupMicroFrontend(
						MicroFrontendName,
						service,
						modulePath,
						&proxy,
						ptr("test-microfrontendclass"),
						[]polyfeav1alpha1.StaticResources{{Path: "static", Kind: "script"}},
						nil,
					)
					if shouldSucceed {
						Expect(k8sClient.Create(testCtx, microFrontend)).Should(Succeed())
					} else {
						Expect(k8sClient.Create(testCtx, microFrontend)).Should(Not(Succeed()))
					}
				},
				Entry("missing service", nil, ptr("module.jsm"), false),
				Entry("missing modulePath", ptr("http://test-service.test-namespace.svc.cluster.local"), nil, false),
				Entry("valid MicroFrontend", ptr("http://test-service.test-namespace.svc.cluster.local"), ptr("module.jsm"), true),
			)

			It("Should create with defaults if optional fields are not specified", func() {
				By("Creating a new MicroFrontend with only required fields")
				testCtx := context.Background()

				ensureMicroFrontendDeleted(testCtx, timeout, interval)

				microFrontend := setupMicroFrontend(
					MicroFrontendName,
					ptr("http://test-service.test-namespace.svc.cluster.local"),
					ptr("module.jsm"),
					nil,
					nil,
					nil,
					nil,
				)
				Expect(k8sClient.Create(testCtx, microFrontend)).Should(Succeed())

				microFrontendLookupKey := types.NamespacedName{Name: MicroFrontendName, Namespace: MicroFrontendNamespace}
				createdMicroFrontend := &polyfeav1alpha1.MicroFrontend{}

				Eventually(func() bool {
					return k8sClient.Get(testCtx, microFrontendLookupKey, createdMicroFrontend) == nil
				}, timeout, interval).Should(BeTrue())
				Expect(createdMicroFrontend.Spec.Service).Should(Equal(ptr("http://test-service.test-namespace.svc.cluster.local")))
				Expect(*createdMicroFrontend.Spec.Proxy).Should(Equal(true))
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
		})
	})
})
