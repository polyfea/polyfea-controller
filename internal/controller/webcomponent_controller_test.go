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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	polyfeav1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
)

const (
	WebComponentName      = "test-webcomponent"
	WebComponentNamespace = "default"
	WebComponentFinalizer = "polyfea.github.io/finalizer"
)

// Ensure the WebComponent resource does not already exist
func ensureWebComponentDeleted(ctx context.Context, timeout time.Duration, interval time.Duration) {
	existingWebComponent := &polyfeav1alpha1.WebComponent{}
	err := k8sClient.Get(ctx, types.NamespacedName{Name: WebComponentName, Namespace: WebComponentNamespace}, existingWebComponent)
	if err == nil {
		// Check if already being deleted
		if existingWebComponent.DeletionTimestamp == nil {
			// Not yet marked for deletion, delete it now
			Expect(k8sClient.Delete(ctx, existingWebComponent)).Should(Succeed())
		}
		// Wait for the resource to be fully deleted (whether we just deleted it or it was already being deleted)
		Eventually(func() bool {
			return errors.IsNotFound(k8sClient.Get(ctx, types.NamespacedName{Name: WebComponentName, Namespace: WebComponentNamespace}, existingWebComponent))
		}, timeout, interval).Should(BeTrue())
		// Give extra time for repository cleanup
		time.Sleep(time.Millisecond * 500)
	}
}

var _ = Describe("WebComponent Controller", func() {
	const (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When reconciling a resource", func() {
		DescribeTable("Validation scenarios",
			func(spec polyfeav1alpha1.WebComponentSpec, shouldSucceed bool) {
				testCtx := context.Background()
				ensureWebComponentDeleted(testCtx, timeout, interval)
				webComponent := &polyfeav1alpha1.WebComponent{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "polyfea.github.io/v1alpha1",
						Kind:       "WebComponent",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      WebComponentName,
						Namespace: WebComponentNamespace,
					},
					Spec: spec,
				}
				if shouldSucceed {
					Expect(k8sClient.Create(testCtx, webComponent)).Should(Succeed())
				} else {
					Expect(k8sClient.Create(testCtx, webComponent)).Should(Not(Succeed()))
				}
			},
			Entry("missing element", polyfeav1alpha1.WebComponentSpec{
				MicroFrontend: &[]string{"test-microfrontend"}[0],
				Attributes: []polyfeav1alpha1.Attribute{
					{Name: "label", Value: runtime.RawExtension{Raw: []byte(`"My Menu Item"`)}},
				},
				DisplayRules: []polyfeav1alpha1.DisplayRules{{
					AllOf: []polyfeav1alpha1.Matcher{{Path: "pathname"}, {ContextName: "user"}},
				}},
				Priority: &[]int32{0}[0],
				Style:    []polyfeav1alpha1.Style{{Name: "color", Value: "red"}},
			}, false),
			Entry("missing display rules", polyfeav1alpha1.WebComponentSpec{
				Element:       &[]string{"my-menu-item"}[0],
				MicroFrontend: &[]string{"test-microfrontend"}[0],
				Attributes: []polyfeav1alpha1.Attribute{
					{Name: "label", Value: runtime.RawExtension{Raw: []byte(`"My Menu Item"`)}},
				},
				Priority: &[]int32{0}[0],
				Style:    []polyfeav1alpha1.Style{{Name: "color", Value: "red"}},
			}, false),
			Entry("valid WebComponent", polyfeav1alpha1.WebComponentSpec{
				Element: &[]string{"my-menu-item"}[0],
				DisplayRules: []polyfeav1alpha1.DisplayRules{{
					AllOf: []polyfeav1alpha1.Matcher{{Path: "pathname"}, {ContextName: "user"}},
				}},
			}, true),
		)

		It("Should create with defaults when optional fields are missing", func() {
			testCtx := context.Background()
			ensureWebComponentDeleted(testCtx, timeout, interval)

			webComponent := &polyfeav1alpha1.WebComponent{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "WebComponent",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      WebComponentName,
					Namespace: WebComponentNamespace,
				},
				Spec: polyfeav1alpha1.WebComponentSpec{
					Element: &[]string{"my-menu-item"}[0],
					DisplayRules: []polyfeav1alpha1.DisplayRules{{
						AllOf: []polyfeav1alpha1.Matcher{{Path: "pathname"}, {ContextName: "user"}},
					}},
				},
			}
			Expect(k8sClient.Create(testCtx, webComponent)).Should(Succeed())

			webComponentLookupKey := types.NamespacedName{Name: WebComponentName, Namespace: WebComponentNamespace}
			createdWebComponent := &polyfeav1alpha1.WebComponent{}

			Eventually(func() bool {
				return k8sClient.Get(testCtx, webComponentLookupKey, createdWebComponent) == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdWebComponent.Spec.MicroFrontend).Should(BeNil())
			Expect(createdWebComponent.Spec.Element).Should(Equal(&[]string{"my-menu-item"}[0]))
			Expect(createdWebComponent.Spec.Priority).Should(Equal(&[]int32{0}[0]))
			Expect(createdWebComponent.Spec.Attributes).Should(BeNil())
			Expect(createdWebComponent.Spec.Style).Should(BeNil())

			Expect(k8sClient.Delete(testCtx, createdWebComponent)).Should(Succeed())
			Eventually(func() bool {
				return errors.IsNotFound(k8sClient.Get(testCtx, webComponentLookupKey, createdWebComponent))
			}, timeout, interval).Should(BeTrue())
		})

		It("Should add and remove finalizer", func() {
			testCtx := context.Background()
			ensureWebComponentDeleted(testCtx, timeout, interval)
			webComponent := &polyfeav1alpha1.WebComponent{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "WebComponent",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      WebComponentName,
					Namespace: WebComponentNamespace,
				},
				Spec: polyfeav1alpha1.WebComponentSpec{
					MicroFrontend: &[]string{"test-microfrontend"}[0],
					Element:       &[]string{"my-menu-item"}[0],
					Attributes: []polyfeav1alpha1.Attribute{
						{Name: "label", Value: runtime.RawExtension{Raw: []byte(`"My Menu Item"`)}},
					},
					DisplayRules: []polyfeav1alpha1.DisplayRules{{
						AllOf: []polyfeav1alpha1.Matcher{{Path: "pathname"}, {ContextName: "user"}},
					}},
					Priority: &[]int32{0}[0],
					Style:    []polyfeav1alpha1.Style{{Name: "color", Value: "red"}},
				},
			}
			Expect(k8sClient.Create(testCtx, webComponent)).Should(Succeed())

			webComponentLookupKey := types.NamespacedName{Name: WebComponentName, Namespace: WebComponentNamespace}
			createdWebComponent := &polyfeav1alpha1.WebComponent{}

			Eventually(func() bool {
				return k8sClient.Get(testCtx, webComponentLookupKey, createdWebComponent) == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdWebComponent.Spec.Element).Should(Equal(&[]string{"my-menu-item"}[0]))

			Eventually(func() []*polyfeav1alpha1.WebComponent {
				result, _ := webComponentRepository.List(func(mf *polyfeav1alpha1.WebComponent) bool {
					return mf.Name == WebComponentName
				})
				return result
			}, timeout, interval).Should(HaveLen(1))

			Expect(k8sClient.Delete(testCtx, createdWebComponent)).Should(Succeed())
			Eventually(func() bool {
				return errors.IsNotFound(k8sClient.Get(testCtx, webComponentLookupKey, createdWebComponent))
			}, timeout, interval).Should(BeTrue())

			Eventually(func() []*polyfeav1alpha1.WebComponent {
				result, _ := webComponentRepository.List(func(mf *polyfeav1alpha1.WebComponent) bool {
					return mf.Name == WebComponentName
				})
				return result
			}, timeout, interval).Should(BeEmpty())
		})
	})
})
