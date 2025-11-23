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

func setupMicroFrontend(service *v1alpha1.ServiceReference, modulePath *string, proxy *bool, frontendClass *string, staticResources []v1alpha1.StaticResources) *v1alpha1.MicroFrontend {
	return &v1alpha1.MicroFrontend{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "polyfea.github.io/v1alpha1",
			Kind:       "MicroFrontend",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MicroFrontendName,
			Namespace: MicroFrontendNamespace,
		},
		Spec: v1alpha1.MicroFrontendSpec{
			Service:         service,
			Proxy:           proxy,
			CacheStrategy:   "none",
			ModulePath:      modulePath,
			StaticResources: staticResources,
			FrontendClass:   frontendClass,
			CacheOptions:    nil,
		},
	}
}

func ensureMicroFrontendDeleted(ctx context.Context) {
	existingMicroFrontend := &v1alpha1.MicroFrontend{}
	err := k8sClient.Get(ctx, types.NamespacedName{Name: MicroFrontendName, Namespace: MicroFrontendNamespace}, existingMicroFrontend)
	if err == nil {
		// Check if already being deleted
		if existingMicroFrontend.DeletionTimestamp == nil {
			// Not yet marked for deletion, delete it now
			Expect(k8sClient.Delete(ctx, existingMicroFrontend)).To(Succeed())
		}
		// Wait for the resource to be fully deleted (whether we just deleted it or it was already being deleted)
		Eventually(func() bool {
			return errors.IsNotFound(k8sClient.Get(ctx, types.NamespacedName{Name: MicroFrontendName, Namespace: MicroFrontendNamespace}, existingMicroFrontend))
		}, time.Second*10, time.Millisecond*250).Should(BeTrue())
		// Give extra time for repository cleanup
		time.Sleep(time.Millisecond * 500)
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

				ensureMicroFrontendDeleted(testCtx)

				microFrontend := setupMicroFrontend(
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

					ensureMicroFrontendDeleted(testCtx)

					microFrontend := setupMicroFrontend(
						service,
						modulePath,
						&proxy,
						ptr("test-microfrontendclass"),
						[]v1alpha1.StaticResources{{Path: "static", Kind: "script"}},
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

				ensureMicroFrontendDeleted(testCtx)
				microFrontend := setupMicroFrontend(
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
				Expect(createdMicroFrontend.Spec.FrontendClass).Should(Equal(ptr(DefaultFrontendClassName)))
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

				ensureMicroFrontendDeleted(testCtx)

				microFrontend := setupMicroFrontend(
					&v1alpha1.ServiceReference{
						URI: ptr("https://cdn.example.com"),
					},
					ptr("modules/app.js"),
					&proxy,
					ptr("test-microfrontendclass"),
					[]v1alpha1.StaticResources{{Path: "styles/main.css", Kind: "stylesheet"}},
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

			It("Should reconcile MicroFrontends when MicroFrontendClass namespace policy changes", func() {
				By("Creating a MicroFrontendClass with 'All' namespace policy")
				testCtx := context.Background()

				ensureMicroFrontendDeleted(testCtx)
				ensureMicroFrontendClassDeleted(testCtx)

				// Create MicroFrontendClass
				mfcName := "test-microfrontendclass"
				mfc := &v1alpha1.MicroFrontendClass{
					ObjectMeta: metav1.ObjectMeta{
						Name:      mfcName,
						Namespace: MicroFrontendNamespace,
					},
					Spec: v1alpha1.MicroFrontendClassSpec{
						BaseUri: ptr("https://test.example.com"),
						Title:   ptr("Test Frontend"),
						NamespacePolicy: &v1alpha1.NamespacePolicy{
							From: v1alpha1.NamespaceFromAll,
						},
					},
				}
				Expect(k8sClient.Create(testCtx, mfc)).To(Succeed()) // Create MicroFrontend referencing this class
				proxy := true
				microFrontend := setupMicroFrontend(
					&v1alpha1.ServiceReference{
						URI: ptr("https://example.com"),
					},
					ptr("app.js"),
					&proxy,
					&mfcName,
					nil,
				)
				Expect(k8sClient.Create(testCtx, microFrontend)).To(Succeed())

				microFrontendLookupKey := types.NamespacedName{Name: MicroFrontendName, Namespace: MicroFrontendNamespace}
				createdMicroFrontend := &v1alpha1.MicroFrontend{}

				By("Verifying MicroFrontend is accepted with 'All' policy")
				Eventually(func() bool {
					err := k8sClient.Get(testCtx, microFrontendLookupKey, createdMicroFrontend)
					if err != nil {
						return false
					}
					return createdMicroFrontend.Status.Phase == v1alpha1.MicroFrontendPhaseReady &&
						createdMicroFrontend.Status.FrontendClassRef != nil &&
						createdMicroFrontend.Status.FrontendClassRef.Accepted
				}, timeout, interval).Should(BeTrue())

				By("Updating MicroFrontendClass to 'Same' namespace policy")
				mfcLookupKey := types.NamespacedName{Name: mfcName, Namespace: MicroFrontendNamespace}
				updatedMfc := &v1alpha1.MicroFrontendClass{}
				Eventually(func() error {
					// Re-fetch to get latest resourceVersion
					if err := k8sClient.Get(testCtx, mfcLookupKey, updatedMfc); err != nil {
						return err
					}
					updatedMfc.Spec.NamespacePolicy = &v1alpha1.NamespacePolicy{
						From: v1alpha1.NamespaceFromSame,
					}
					return k8sClient.Update(testCtx, updatedMfc)
				}, timeout, interval).Should(Succeed())

				By("Verifying MicroFrontend is still accepted (same namespace)")
				Eventually(func() bool {
					err := k8sClient.Get(testCtx, microFrontendLookupKey, createdMicroFrontend)
					if err != nil {
						return false
					}
					// Status should be updated due to the watch triggering reconciliation
					return createdMicroFrontend.Status.Phase == v1alpha1.MicroFrontendPhaseReady &&
						createdMicroFrontend.Status.FrontendClassRef != nil &&
						createdMicroFrontend.Status.FrontendClassRef.Accepted
				}, timeout, interval).Should(BeTrue())

				By("Updating MicroFrontendClass to specific namespaces (excluding default)")
				Eventually(func() error {
					// Re-fetch to get latest resourceVersion
					if err := k8sClient.Get(testCtx, mfcLookupKey, updatedMfc); err != nil {
						return err
					}
					updatedMfc.Spec.NamespacePolicy = &v1alpha1.NamespacePolicy{
						From:       v1alpha1.NamespaceFromNamespaces,
						Namespaces: []string{"other-namespace"},
					}
					return k8sClient.Update(testCtx, updatedMfc)
				}, timeout, interval).Should(Succeed())

				By("Verifying MicroFrontend is now rejected")
				Eventually(func() bool {
					err := k8sClient.Get(testCtx, microFrontendLookupKey, createdMicroFrontend)
					if err != nil {
						return false
					}
					// Watch should trigger reconciliation and MicroFrontend should be rejected
					return createdMicroFrontend.Status.Phase == v1alpha1.MicroFrontendPhaseRejected &&
						createdMicroFrontend.Status.FrontendClassRef != nil &&
						!createdMicroFrontend.Status.FrontendClassRef.Accepted &&
						createdMicroFrontend.Status.RejectionReason != ""
				}, timeout, interval).Should(BeTrue())

				By("Cleaning up")
				Expect(k8sClient.Delete(testCtx, createdMicroFrontend)).Should(Succeed())
				Expect(k8sClient.Delete(testCtx, updatedMfc)).Should(Succeed())
			})

			It("Should remove rejected MicroFrontends from repository", func() {
				By("Creating a MicroFrontendClass with specific namespaces policy")
				testCtx := context.Background()

				ensureMicroFrontendDeleted(testCtx)
				ensureMicroFrontendClassDeleted(testCtx)

				// Create MicroFrontendClass that only allows 'allowed-namespace'
				mfcName := "test-microfrontendclass"
				mfc := &v1alpha1.MicroFrontendClass{
					ObjectMeta: metav1.ObjectMeta{
						Name:      mfcName,
						Namespace: MicroFrontendNamespace,
					},
					Spec: v1alpha1.MicroFrontendClassSpec{
						BaseUri: ptr("https://test.example.com"),
						Title:   ptr("Test Frontend"),
						NamespacePolicy: &v1alpha1.NamespacePolicy{
							From:       v1alpha1.NamespaceFromNamespaces,
							Namespaces: []string{"allowed-namespace"},
						},
					},
				}
				Expect(k8sClient.Create(testCtx, mfc)).To(Succeed())

				By("Creating a MicroFrontend in 'default' namespace (not allowed)")
				proxy := true
				microFrontend := setupMicroFrontend(
					&v1alpha1.ServiceReference{
						URI: ptr("https://example.com"),
					},
					ptr("app.js"),
					&proxy,
					&mfcName,
					nil,
				)
				Expect(k8sClient.Create(testCtx, microFrontend)).To(Succeed())

				microFrontendLookupKey := types.NamespacedName{Name: MicroFrontendName, Namespace: MicroFrontendNamespace}
				createdMicroFrontend := &v1alpha1.MicroFrontend{}

				By("Verifying MicroFrontend is rejected due to namespace policy")
				Eventually(func() bool {
					err := k8sClient.Get(testCtx, microFrontendLookupKey, createdMicroFrontend)
					if err != nil {
						return false
					}
					return createdMicroFrontend.Status.Phase == v1alpha1.MicroFrontendPhaseRejected &&
						createdMicroFrontend.Status.FrontendClassRef != nil &&
						!createdMicroFrontend.Status.FrontendClassRef.Accepted
				}, timeout, interval).Should(BeTrue())

				By("Verifying rejected MicroFrontend is not in repository")
				// Try to get the MicroFrontend from repository
				repoMf, _ := microFrontendRepository.Get(createdMicroFrontend)
				Expect(repoMf).Should(BeNil()) // Should be nil since it's not in repository

				By("Updating MicroFrontendClass to allow 'default' namespace")
				mfcLookupKey := types.NamespacedName{Name: mfcName, Namespace: MicroFrontendNamespace}
				updatedMfc := &v1alpha1.MicroFrontendClass{}
				Eventually(func() error {
					if err := k8sClient.Get(testCtx, mfcLookupKey, updatedMfc); err != nil {
						return err
					}
					updatedMfc.Spec.NamespacePolicy = &v1alpha1.NamespacePolicy{
						From:       v1alpha1.NamespaceFromNamespaces,
						Namespaces: []string{"allowed-namespace", "default"},
					}
					return k8sClient.Update(testCtx, updatedMfc)
				}, timeout, interval).Should(Succeed())

				By("Verifying MicroFrontend is now accepted and added to repository")
				Eventually(func() bool {
					err := k8sClient.Get(testCtx, microFrontendLookupKey, createdMicroFrontend)
					if err != nil {
						return false
					}
					return createdMicroFrontend.Status.Phase == v1alpha1.MicroFrontendPhaseReady &&
						createdMicroFrontend.Status.FrontendClassRef != nil &&
						createdMicroFrontend.Status.FrontendClassRef.Accepted
				}, timeout, interval).Should(BeTrue())

				By("Verifying accepted MicroFrontend is now in repository")
				Eventually(func() bool {
					repoMf, err := microFrontendRepository.Get(createdMicroFrontend)
					return err == nil && repoMf != nil
				}, timeout, interval).Should(BeTrue())

				By("Cleaning up")
				Expect(k8sClient.Delete(testCtx, createdMicroFrontend)).Should(Succeed())
				Expect(k8sClient.Delete(testCtx, updatedMfc)).Should(Succeed())
			})
		})
	})
})
