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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("WebComponent controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		WebComponentName      = "test-webcomponent"
		WebComponentNamespace = "default"
		WebComponentFinalizer = "polyfea.github.io/finalizer"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating a WebComponent", func() {
		It("Should add and remove finallizer", func() {
			By("By creating a new WebComponent")
			ctx := context.Background()
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
						{
							Name:  "label",
							Value: runtime.RawExtension{Raw: []byte(`"My Menu Item"`)},
						},
					},
					DisplayRules: &polyfeav1alpha1.DisplayRules{
						AllOf: []polyfeav1alpha1.Matcher{
							{
								Path: "pathname",
							},
							{
								ContextName: "user",
							},
						},
					},
					Priority: &[]int32{0}[0],
					Style:    &[]string{"color: red;"}[0],
				},
			}
			Expect(k8sClient.Create(ctx, webComponent)).Should(Succeed())

			webComponentLookupKey := types.NamespacedName{Name: WebComponentName, Namespace: WebComponentNamespace}
			createdWebComponent := &polyfeav1alpha1.WebComponent{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, webComponentLookupKey, createdWebComponent)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdWebComponent.Spec.Element).Should(Equal(&[]string{"my-menu-item"}[0]))

			By("By checking the WebComponent has finalizer")
			Eventually(func() ([]string, error) {
				err := k8sClient.Get(ctx, webComponentLookupKey, createdWebComponent)
				if err != nil {
					return []string{}, err
				}
				return createdWebComponent.ObjectMeta.GetFinalizers(), nil
			}, timeout, interval).Should(ContainElement(WebComponentFinalizer))

			By("By deleting the WebComponent")
			Expect(k8sClient.Delete(ctx, createdWebComponent)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, webComponentLookupKey, createdWebComponent)
				return err != nil
			}, timeout, interval).Should(BeTrue())
		})

		It("Should not create if required parameter is missing", func() {
			By("By creating a new WebComponent without microfrontend")
			ctx := context.Background()
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
					Attributes: []polyfeav1alpha1.Attribute{
						{
							Name:  "label",
							Value: runtime.RawExtension{Raw: []byte(`"My Menu Item"`)},
						},
					},
					DisplayRules: &polyfeav1alpha1.DisplayRules{
						AllOf: []polyfeav1alpha1.Matcher{
							{
								Path: "pathname",
							},
							{
								ContextName: "user",
							},
						},
					},
					Priority: &[]int32{0}[0],
					Style:    &[]string{"color: red;"}[0],
				},
			}
			Expect(k8sClient.Create(ctx, webComponent)).Should(Not(Succeed()))

			By("By creating a new WebComponent without element")
			ctx = context.Background()
			webComponent = &polyfeav1alpha1.WebComponent{
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
					Attributes: []polyfeav1alpha1.Attribute{
						{
							Name:  "label",
							Value: runtime.RawExtension{Raw: []byte(`"My Menu Item"`)},
						},
					},
					DisplayRules: &polyfeav1alpha1.DisplayRules{
						AllOf: []polyfeav1alpha1.Matcher{
							{
								Path: "pathname",
							},
							{
								ContextName: "user",
							},
						},
					},
					Priority: &[]int32{0}[0],
					Style:    &[]string{"color: red;"}[0],
				},
			}
			Expect(k8sClient.Create(ctx, webComponent)).Should(Not(Succeed()))
		})

		It("Should create with defaults when optional fields are missing", func() {
			By("By creating a new WebComponent without optional fields")
			ctx := context.Background()
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
				},
			}
			Expect(k8sClient.Create(ctx, webComponent)).Should(Succeed())

			webComponentLookupKey := types.NamespacedName{Name: WebComponentName, Namespace: WebComponentNamespace}
			createdWebComponent := &polyfeav1alpha1.WebComponent{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, webComponentLookupKey, createdWebComponent)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdWebComponent.Spec.Element).Should(Equal(&[]string{"my-menu-item"}[0]))
			Expect(createdWebComponent.Spec.Priority).Should(Equal(&[]int32{0}[0]))
			Expect(createdWebComponent.Spec.Attributes).Should(BeNil())
			Expect(createdWebComponent.Spec.DisplayRules).Should(BeNil())
			Expect(createdWebComponent.Spec.Style).Should(BeNil())
		})
	})
})
