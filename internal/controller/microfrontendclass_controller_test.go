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

	"k8s.io/apimachinery/pkg/api/errors" // Correct import for `errors.IsNotFound`

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
)

const (
	MicroFrontendClassName      = "test-microfrontendclass"
	MicroFrontendClassNamespace = "default"
	MicroFrontendClassFinalizer = "polyfea.github.io/finalizer"
)

func setupMicroFrontendClass(title, baseUri *string) *v1alpha1.MicroFrontendClass {
	return &v1alpha1.MicroFrontendClass{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "polyfea.github.io/v1alpha1",
			Kind:       "MicroFrontendClass",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      MicroFrontendClassName,
			Namespace: MicroFrontendClassNamespace,
		},
		Spec: v1alpha1.MicroFrontendClassSpec{
			Title:   title,
			BaseUri: baseUri,
		},
	}
}

func ensureMicroFrontendClassDeleted(ctx context.Context) {
	existingMicroFrontendClass := &v1alpha1.MicroFrontendClass{}
	err := k8sClient.Get(ctx, types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}, existingMicroFrontendClass)
	if err == nil {
		// Check if already being deleted
		if existingMicroFrontendClass.DeletionTimestamp == nil {
			// Not yet marked for deletion, delete it now
			Expect(k8sClient.Delete(ctx, existingMicroFrontendClass)).Should(Succeed())
		}
		// Wait for the resource to be fully deleted (whether we just deleted it or it was already being deleted)
		Eventually(func() bool {
			return errors.IsNotFound(k8sClient.Get(ctx, types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}, existingMicroFrontendClass))
		}, time.Second*10, time.Millisecond*250).Should(BeTrue())
		// Give extra time for repository cleanup
		time.Sleep(time.Millisecond * 500)
	}
}

var _ = Describe("MicroFrontendClass Controller", func() {
	const (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When reconciling a resource", func() {
		Context("Validation and defaults", func() {
			DescribeTable("Validation scenarios",
				func(title, baseUri *string, shouldSucceed bool) {
					testCtx := context.Background()

					ensureMicroFrontendClassDeleted(testCtx)

					mfc := setupMicroFrontendClass(title, baseUri)
					if shouldSucceed {
						Expect(k8sClient.Create(testCtx, mfc)).Should(Succeed())
					} else {
						Expect(k8sClient.Create(testCtx, mfc)).Should(Not(Succeed()))
					}
				},
				Entry("missing baseUri", ptr("Test"), nil, false),
				Entry("missing title", nil, ptr("/base"), false),
				Entry("valid MicroFrontendClass", ptr("Test"), ptr("/base"), true),
			)

			It("Should fill defaults when only required fields are set", func() {
				testCtx := context.Background()

				ensureMicroFrontendClassDeleted(testCtx)

				mfc := setupMicroFrontendClass(ptr("Test"), ptr("/base"))
				Expect(k8sClient.Create(testCtx, mfc)).Should(Succeed())

				lookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
				created := &v1alpha1.MicroFrontendClass{}
				Eventually(func() bool {
					return k8sClient.Get(testCtx, lookupKey, created) == nil
				}, timeout, interval).Should(BeTrue())
				Expect(created.Spec.CspHeader).ShouldNot(BeEmpty())
				Expect(created.Spec.UserRolesHeader).ShouldNot(BeEmpty())
				Expect(created.Spec.UserHeader).ShouldNot(BeEmpty())

				Expect(k8sClient.Delete(testCtx, created)).Should(Succeed())
				Eventually(func() bool {
					return errors.IsNotFound(k8sClient.Get(testCtx, lookupKey, created))
				}, timeout, interval).Should(BeTrue())
			})
		})

		Context("Repository add/remove", func() {
			It("Should add and remove MicroFrontendClass from repository", func() {
				testCtx := context.Background()

				ensureMicroFrontendClassDeleted(testCtx)

				mfc := setupMicroFrontendClass(ptr("Test"), ptr("/base"))
				Expect(k8sClient.Create(testCtx, mfc)).Should(Succeed())

				lookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
				created := &v1alpha1.MicroFrontendClass{}
				Eventually(func() bool {
					return k8sClient.Get(testCtx, lookupKey, created) == nil
				}, timeout, interval).Should(BeTrue())

				Eventually(func() []*v1alpha1.MicroFrontendClass {
					result, _ := microFrontendClassRepository.List(func(m *v1alpha1.MicroFrontendClass) bool {
						return m.Name == MicroFrontendClassName
					})
					return result
				}, timeout, interval).Should(HaveLen(1))

				Expect(k8sClient.Delete(testCtx, created)).Should(Succeed())
				Eventually(func() []*v1alpha1.MicroFrontendClass {
					result, _ := microFrontendClassRepository.List(func(m *v1alpha1.MicroFrontendClass) bool {
						return m.Name == MicroFrontendClassName
					})
					return result
				}, timeout, interval).Should(BeEmpty())
			})
		})

		Context("Manifest validation", func() {
			DescribeTable("Manifest scenarios",
				func(manifest *v1alpha1.WebAppManifest, shouldSucceed bool) {
					testCtx := context.Background()

					ensureMicroFrontendClassDeleted(testCtx)

					mfc := setupMicroFrontendClass(ptr("Test"), ptr("/base"))
					mfc.Spec.ProgressiveWebApp = &v1alpha1.ProgressiveWebApp{WebAppManifest: manifest}
					if shouldSucceed {
						Expect(k8sClient.Create(testCtx, mfc)).Should(Succeed())
					} else {
						Expect(k8sClient.Create(testCtx, mfc)).Should(Not(Succeed()))
					}
				},
				Entry("missing name", &v1alpha1.WebAppManifest{
					Icons:    []v1alpha1.PWAIcon{{Src: ptr("icon.png"), Sizes: ptr("192x192"), Type: ptr("image/png")}},
					StartUrl: ptr("/"), Display: ptr("standalone"),
				}, false),
				Entry("missing icons", &v1alpha1.WebAppManifest{
					Name: ptr("Test"), StartUrl: ptr("/"), Display: ptr("standalone"),
				}, false),
				Entry("missing startUrl", &v1alpha1.WebAppManifest{
					Name: ptr("Test"), Icons: []v1alpha1.PWAIcon{{Src: ptr("icon.png"), Sizes: ptr("192x192"), Type: ptr("image/png")}},
					Display: ptr("standalone"),
				}, false),
				Entry("missing display", &v1alpha1.WebAppManifest{
					Name: ptr("Test"), Icons: []v1alpha1.PWAIcon{{Src: ptr("icon.png"), Sizes: ptr("192x192"), Type: ptr("image/png")}},
					StartUrl: ptr("/"),
				}, false),
				Entry("valid manifest", &v1alpha1.WebAppManifest{
					Name: ptr("Test"), Icons: []v1alpha1.PWAIcon{{Src: ptr("icon.png"), Sizes: ptr("192x192"), Type: ptr("image/png")}},
					StartUrl: ptr("/"), Display: ptr("standalone"),
				}, true),
			)
		})

		Context("Namespace Policy", func() {
			It("Should default to 'All' namespace policy when not specified", func() {
				testCtx := context.Background()
				ensureMicroFrontendClassDeleted(testCtx)

				mfc := setupMicroFrontendClass(ptr("Test"), ptr("/base"))
				// No namespace policy specified
				Expect(k8sClient.Create(testCtx, mfc)).Should(Succeed())

				lookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
				created := &v1alpha1.MicroFrontendClass{}
				Eventually(func() bool {
					return k8sClient.Get(testCtx, lookupKey, created) == nil
				}, timeout, interval).Should(BeTrue())

				// Should allow all namespaces (default behavior)
				Expect(created.IsNamespaceAllowed("default")).Should(BeTrue())
				Expect(created.IsNamespaceAllowed("other-namespace")).Should(BeTrue())
				Expect(created.IsNamespaceAllowed("production")).Should(BeTrue())

				Expect(k8sClient.Delete(testCtx, created)).Should(Succeed())
			})

			It("Should respect 'All' namespace policy", func() {
				testCtx := context.Background()
				ensureMicroFrontendClassDeleted(testCtx)

				mfc := setupMicroFrontendClass(ptr("Test All Policy"), ptr("/all"))
				mfc.Spec.NamespacePolicy = &v1alpha1.NamespacePolicy{
					From: v1alpha1.NamespaceFromAll,
				}
				Expect(k8sClient.Create(testCtx, mfc)).Should(Succeed())

				lookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
				created := &v1alpha1.MicroFrontendClass{}
				Eventually(func() bool {
					return k8sClient.Get(testCtx, lookupKey, created) == nil
				}, timeout, interval).Should(BeTrue())

				// Should allow all namespaces
				Expect(created.IsNamespaceAllowed("default")).Should(BeTrue())
				Expect(created.IsNamespaceAllowed("team-a")).Should(BeTrue())
				Expect(created.IsNamespaceAllowed("team-b")).Should(BeTrue())
				Expect(created.IsNamespaceAllowed("production")).Should(BeTrue())
				Expect(created.IsNamespaceAllowed("any-namespace")).Should(BeTrue())

				// Verify status is updated
				Eventually(func() string {
					err := k8sClient.Get(testCtx, lookupKey, created)
					if err != nil {
						return ""
					}
					return created.Status.Phase
				}, timeout, interval).Should(Equal(v1alpha1.MicroFrontendClassPhaseReady))

				Expect(k8sClient.Delete(testCtx, created)).Should(Succeed())
			})

			It("Should respect 'Same' namespace policy", func() {
				testCtx := context.Background()
				ensureMicroFrontendClassDeleted(testCtx)

				mfc := setupMicroFrontendClass(ptr("Test Same Policy"), ptr("/same"))
				mfc.Spec.NamespacePolicy = &v1alpha1.NamespacePolicy{
					From: v1alpha1.NamespaceFromSame,
				}
				Expect(k8sClient.Create(testCtx, mfc)).Should(Succeed())

				lookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
				created := &v1alpha1.MicroFrontendClass{}
				Eventually(func() error {
					return k8sClient.Get(testCtx, lookupKey, created)
				}, timeout, interval).Should(Succeed())

				// Should only allow the same namespace (default)
				Expect(created.IsNamespaceAllowed("default")).Should(BeTrue())

				// Should reject other namespaces
				Expect(created.IsNamespaceAllowed("team-a")).Should(BeFalse())
				Expect(created.IsNamespaceAllowed("team-b")).Should(BeFalse())
				Expect(created.IsNamespaceAllowed("production")).Should(BeFalse())
				Expect(created.IsNamespaceAllowed("other-namespace")).Should(BeFalse())

				// Verify status is updated
				Eventually(func() string {
					err := k8sClient.Get(testCtx, lookupKey, created)
					if err != nil {
						return ""
					}
					return created.Status.Phase
				}, timeout, interval).Should(Equal(v1alpha1.MicroFrontendClassPhaseReady))

				Expect(k8sClient.Delete(testCtx, created)).Should(Succeed())
			})

			It("Should respect 'FromNamespaces' policy with specific namespaces", func() {
				testCtx := context.Background()
				ensureMicroFrontendClassDeleted(testCtx)

				mfc := setupMicroFrontendClass(ptr("Test FromNamespaces Policy"), ptr("/fromnamespaces"))
				mfc.Spec.NamespacePolicy = &v1alpha1.NamespacePolicy{
					From: v1alpha1.NamespaceFromNamespaces,
					Namespaces: []string{
						"team-a",
						"team-b",
						"production",
					},
				}
				Expect(k8sClient.Create(testCtx, mfc)).Should(Succeed())

				lookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
				created := &v1alpha1.MicroFrontendClass{}
				Eventually(func() error {
					return k8sClient.Get(testCtx, lookupKey, created)
				}, timeout, interval).Should(Succeed())

				// Should allow specified namespaces
				Expect(created.IsNamespaceAllowed("team-a")).Should(BeTrue())
				Expect(created.IsNamespaceAllowed("team-b")).Should(BeTrue())
				Expect(created.IsNamespaceAllowed("production")).Should(BeTrue())

				// Should reject other namespaces
				Expect(created.IsNamespaceAllowed("default")).Should(BeFalse())
				Expect(created.IsNamespaceAllowed("staging")).Should(BeFalse())
				Expect(created.IsNamespaceAllowed("development")).Should(BeFalse())
				Expect(created.IsNamespaceAllowed("other-namespace")).Should(BeFalse())

				// Verify status is updated
				Eventually(func() string {
					err := k8sClient.Get(testCtx, lookupKey, created)
					if err != nil {
						return ""
					}
					return created.Status.Phase
				}, timeout, interval).Should(Equal(v1alpha1.MicroFrontendClassPhaseReady))

				Expect(k8sClient.Delete(testCtx, created)).Should(Succeed())
			})

			It("Should handle empty namespace list in 'FromNamespaces' policy", func() {
				testCtx := context.Background()
				ensureMicroFrontendClassDeleted(testCtx)

				mfc := setupMicroFrontendClass(ptr("Test Empty Namespaces"), ptr("/empty"))
				mfc.Spec.NamespacePolicy = &v1alpha1.NamespacePolicy{
					From:       v1alpha1.NamespaceFromNamespaces,
					Namespaces: []string{}, // Empty list
				}
				Expect(k8sClient.Create(testCtx, mfc)).Should(Succeed())

				lookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
				created := &v1alpha1.MicroFrontendClass{}
				Eventually(func() error {
					return k8sClient.Get(testCtx, lookupKey, created)
				}, timeout, interval).Should(Succeed())

				// Should reject all namespaces when list is empty
				Expect(created.IsNamespaceAllowed("default")).Should(BeFalse())
				Expect(created.IsNamespaceAllowed("team-a")).Should(BeFalse())
				Expect(created.IsNamespaceAllowed("production")).Should(BeFalse())

				Expect(k8sClient.Delete(testCtx, created)).Should(Succeed())
			})

			It("Should handle 'FromNamespaces' with single namespace", func() {
				testCtx := context.Background()
				ensureMicroFrontendClassDeleted(testCtx)

				mfc := setupMicroFrontendClass(ptr("Test Single Namespace"), ptr("/single"))
				mfc.Spec.NamespacePolicy = &v1alpha1.NamespacePolicy{
					From:       v1alpha1.NamespaceFromNamespaces,
					Namespaces: []string{"production"},
				}
				Expect(k8sClient.Create(testCtx, mfc)).Should(Succeed())

				lookupKey := types.NamespacedName{Name: MicroFrontendClassName, Namespace: MicroFrontendClassNamespace}
				created := &v1alpha1.MicroFrontendClass{}
				Eventually(func() error {
					return k8sClient.Get(testCtx, lookupKey, created)
				}, timeout, interval).Should(Succeed())

				// Should only allow the single specified namespace
				Expect(created.IsNamespaceAllowed("production")).Should(BeTrue())

				// Should reject all other namespaces
				Expect(created.IsNamespaceAllowed("default")).Should(BeFalse())
				Expect(created.IsNamespaceAllowed("team-a")).Should(BeFalse())
				Expect(created.IsNamespaceAllowed("staging")).Should(BeFalse())

				Expect(k8sClient.Delete(testCtx, created)).Should(Succeed())
			})
		})
	})
})
