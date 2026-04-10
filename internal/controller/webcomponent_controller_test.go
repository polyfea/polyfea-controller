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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	polyfeav1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
)

const (
	WebComponentName      = "test-webcomponent"
	WebComponentNamespace = "default"
)

func ensureWebComponentDeleted(ctx context.Context) {
	ensureResourceDeleted(ctx, &polyfeav1alpha1.WebComponent{}, types.NamespacedName{Name: WebComponentName, Namespace: WebComponentNamespace})
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
				ensureWebComponentDeleted(testCtx)
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
			ensureWebComponentDeleted(testCtx)

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
			ensureWebComponentDeleted(testCtx)
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

		It("Should handle MicroFrontend in a different namespace", func() {
			testCtx := context.Background()
			ensureWebComponentDeleted(testCtx)

			// Create a namespace for the MicroFrontend
			mfNamespace := "microfrontend-ns"
			namespaceObj := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: mfNamespace,
				},
			}
			_ = k8sClient.Create(testCtx, namespaceObj)

			// Create MicroFrontend in a different namespace
			microFrontendName := "test-microfrontend-cross-ns"
			microFrontend := &polyfeav1alpha1.MicroFrontend{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontend",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      microFrontendName,
					Namespace: mfNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendSpec{
					Service: &polyfeav1alpha1.ServiceReference{
						Name:      &[]string{"test-service"}[0],
						Namespace: &mfNamespace,
						Port:      &[]int32{80}[0],
						Scheme:    &[]string{"http"}[0],
					},
					ModulePath: &[]string{"/module"}[0],
				},
			}
			Expect(k8sClient.Create(testCtx, microFrontend)).Should(Succeed())

			// Create WebComponent in default namespace that references MicroFrontend in different namespace
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
					MicroFrontend: &microFrontendName,
					Element:       &[]string{"my-menu-item"}[0],
					DisplayRules: []polyfeav1alpha1.DisplayRules{{
						AllOf: []polyfeav1alpha1.Matcher{{Path: "pathname"}, {ContextName: "user"}},
					}},
				},
			}
			Expect(k8sClient.Create(testCtx, webComponent)).Should(Succeed())

			webComponentLookupKey := types.NamespacedName{Name: WebComponentName, Namespace: WebComponentNamespace}
			createdWebComponent := &polyfeav1alpha1.WebComponent{}

			// Eventually the WebComponent should be created
			Eventually(func() bool {
				return k8sClient.Get(testCtx, webComponentLookupKey, createdWebComponent) == nil
			}, timeout, interval).Should(BeTrue())

			// WebComponent should have MicroFrontendRef set with the correct namespace and Found=true
			Eventually(func() bool {
				err := k8sClient.Get(testCtx, webComponentLookupKey, createdWebComponent)
				if err != nil {
					return false
				}
				return createdWebComponent.Status.MicroFrontendRef != nil &&
					createdWebComponent.Status.MicroFrontendRef.Name == microFrontendName &&
					createdWebComponent.Status.MicroFrontendRef.Namespace == mfNamespace &&
					createdWebComponent.Status.MicroFrontendRef.Found == true
			}, timeout, interval).Should(BeTrue())

			// Verify the condition indicates MicroFrontend was found and resolved
			Eventually(func() bool {
				err := k8sClient.Get(testCtx, webComponentLookupKey, createdWebComponent)
				if err != nil {
					return false
				}
				condition := polyfeav1alpha1.GetCondition(createdWebComponent.Status.Conditions,
					polyfeav1alpha1.ConditionTypeMicroFrontendResolved)
				return condition != nil && condition.Status == metav1.ConditionTrue &&
					condition.Reason == polyfeav1alpha1.ReasonSuccessful
			}, timeout, interval).Should(BeTrue())

			// Cleanup
			Expect(k8sClient.Delete(testCtx, createdWebComponent)).Should(Succeed())
			Expect(k8sClient.Delete(testCtx, microFrontend)).Should(Succeed())
			Expect(k8sClient.Delete(testCtx, namespaceObj)).Should(Succeed())
		})

		It("Should accept and round-trip nested matcher operators via the Kubernetes API", func() {
			// This test verifies that the CRD schema (with x-kubernetes-preserve-unknown-fields
			// on the recursive Matcher fields) allows deeply nested allOf/anyOf/noneOf structures
			// to be written to and read back from the API server without any fields being pruned.
			testCtx := context.Background()
			ensureWebComponentDeleted(testCtx)

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
					Element: &[]string{"my-element"}[0],
					DisplayRules: []polyfeav1alpha1.DisplayRules{
						{
							// Top-level allOf containing nested anyOf and noneOf — mirrors the
							// feature request example: show for (nav OR apps) AND (admin OR ds)
							// AND NOT /secret.
							AllOf: []polyfeav1alpha1.Matcher{
								{
									AnyOf: []polyfeav1alpha1.Matcher{
										{ContextName: "navigation"},
										{ContextName: "applications"},
									},
								},
								{
									AnyOf: []polyfeav1alpha1.Matcher{
										{Role: "admin"},
										{Role: "datascience"},
									},
								},
								{
									NoneOf: []polyfeav1alpha1.Matcher{
										{Path: "^/secret$"},
									},
								},
							},
						},
						{
							// Second rule: deeply nested allOf inside anyOf inside noneOf,
							// exercising three levels and all three operators in one rule.
							AnyOf: []polyfeav1alpha1.Matcher{
								{Role: "superadmin"},
								{
									AllOf: []polyfeav1alpha1.Matcher{
										{ContextName: "admin-panel"},
										{
											NoneOf: []polyfeav1alpha1.Matcher{
												{Path: "^/admin/secret"},
											},
										},
									},
								},
							},
						},
					},
				},
			}

			Expect(k8sClient.Create(testCtx, webComponent)).Should(Succeed())

			lookupKey := types.NamespacedName{Name: WebComponentName, Namespace: WebComponentNamespace}
			fetched := &polyfeav1alpha1.WebComponent{}
			Eventually(func() bool {
				return k8sClient.Get(testCtx, lookupKey, fetched) == nil
			}, timeout, interval).Should(BeTrue())

			// Rule 1: allOf with three nested matchers
			Expect(fetched.Spec.DisplayRules).Should(HaveLen(2))
			rule1 := fetched.Spec.DisplayRules[0]
			Expect(rule1.AllOf).Should(HaveLen(3))

			// First allOf entry: anyOf with two context matchers
			Expect(rule1.AllOf[0].AnyOf).Should(HaveLen(2))
			Expect(rule1.AllOf[0].AnyOf[0].ContextName).Should(Equal("navigation"))
			Expect(rule1.AllOf[0].AnyOf[1].ContextName).Should(Equal("applications"))

			// Second allOf entry: anyOf with two role matchers
			Expect(rule1.AllOf[1].AnyOf).Should(HaveLen(2))
			Expect(rule1.AllOf[1].AnyOf[0].Role).Should(Equal("admin"))
			Expect(rule1.AllOf[1].AnyOf[1].Role).Should(Equal("datascience"))

			// Third allOf entry: noneOf with path matcher
			Expect(rule1.AllOf[2].NoneOf).Should(HaveLen(1))
			Expect(rule1.AllOf[2].NoneOf[0].Path).Should(Equal("^/secret$"))

			// Rule 2: anyOf with scalar and nested allOf+noneOf (three levels deep)
			rule2 := fetched.Spec.DisplayRules[1]
			Expect(rule2.AnyOf).Should(HaveLen(2))
			Expect(rule2.AnyOf[0].Role).Should(Equal("superadmin"))
			Expect(rule2.AnyOf[1].AllOf).Should(HaveLen(2))
			Expect(rule2.AnyOf[1].AllOf[0].ContextName).Should(Equal("admin-panel"))
			Expect(rule2.AnyOf[1].AllOf[1].NoneOf).Should(HaveLen(1))
			Expect(rule2.AnyOf[1].AllOf[1].NoneOf[0].Path).Should(Equal("^/admin/secret"))

			Expect(k8sClient.Delete(testCtx, fetched)).Should(Succeed())
			Eventually(func() bool {
				return errors.IsNotFound(k8sClient.Get(testCtx, lookupKey, fetched))
			}, timeout, interval).Should(BeTrue())
		})
	})
})
