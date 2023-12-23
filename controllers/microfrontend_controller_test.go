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
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Microfrontend controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		MicroFrontendName      = "test-microfrontend"
		MicroFrontendNamespace = "default"
		MicroFrontendFinalizer = "polyfea.github.io/finalizer"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating a MicroFrontend", func() {
		It("Should add and remove finallizer", func() {
			By("By creating a new MicroFrontend")
			ctx := context.Background()

			portNumber := int32(8080)
			preload := true

			microFrontend := &polyfeav1alpha1.MicroFrontend{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontend",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendName,
					Namespace: MicroFrontendNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendSpec{
					Service: &polyfeav1alpha1.ServiceReference{
						Name: "test-service",
						Port: &polyfeav1alpha1.Port{
							Number: &portNumber,
						},
					},
					Proxy:         &preload,
					CacheStrategy: "none",
					ModulePath:    &[]string{"module.jsm"}[0],
					StaticPaths:   []string{"static"},
					FrontendClass: &[]string{"test-microfrontendclass"}[0],
				},
			}
			Expect(k8sClient.Create(ctx, microFrontend)).Should(Succeed())

			microFrontendLookupKey := types.NamespacedName{Name: MicroFrontendName, Namespace: MicroFrontendNamespace}
			createdMicroFrontend := &polyfeav1alpha1.MicroFrontend{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendLookupKey, createdMicroFrontend)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontend.Spec.Service.Name).Should(Equal("test-service"))

			By("By checking the MicroFrontend has finalizer")
			Eventually(func() ([]string, error) {
				err := k8sClient.Get(ctx, microFrontendLookupKey, createdMicroFrontend)
				if err != nil {
					return []string{}, err
				}
				return createdMicroFrontend.ObjectMeta.GetFinalizers(), nil
			}, timeout, interval).Should(ContainElement(MicroFrontendFinalizer))

			By("By deleting the MicroFrontend")
			Expect(k8sClient.Delete(ctx, createdMicroFrontend)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendLookupKey, createdMicroFrontend)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())
		})

		It("Should not create if required fields are missing", func() {
			By("By creating a new MicroFrontend without service")
			ctx := context.Background()

			preload := true

			microFrontend := &polyfeav1alpha1.MicroFrontend{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontend",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendName,
					Namespace: MicroFrontendNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendSpec{
					Proxy:         &preload,
					CacheStrategy: "none",
					ModulePath:    &[]string{"module.jsm"}[0],
					StaticPaths:   []string{"static"},
					FrontendClass: &[]string{"test-microfrontendclass"}[0],
				},
			}
			Expect(k8sClient.Create(ctx, microFrontend)).Should(Not(Succeed()))

			By("By creating a new MicroFrontend without module path")
			ctx = context.Background()

			preload = true
			portNumber := int32(8080)

			microFrontend = &polyfeav1alpha1.MicroFrontend{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontend",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendName,
					Namespace: MicroFrontendNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendSpec{
					Service: &polyfeav1alpha1.ServiceReference{
						Name: "test-service",
						Port: &polyfeav1alpha1.Port{
							Number: &portNumber,
						},
					},
					Proxy:         &preload,
					CacheStrategy: "none",
					StaticPaths:   []string{"static"},
					FrontendClass: &[]string{"test-microfrontendclass"}[0],
				},
			}
			Expect(k8sClient.Create(ctx, microFrontend)).Should(Not(Succeed()))

			By("By creating a new MicroFrontend without frontend class")
			ctx = context.Background()

			preload = true
			portNumber = int32(8080)

			microFrontend = &polyfeav1alpha1.MicroFrontend{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontend",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendName,
					Namespace: MicroFrontendNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendSpec{
					Service: &polyfeav1alpha1.ServiceReference{
						Name: "test-service",
						Port: &polyfeav1alpha1.Port{
							Number: &portNumber,
						},
					},
					Proxy:         &preload,
					CacheStrategy: "none",
					ModulePath:    &[]string{"module.jsm"}[0],
					StaticPaths:   []string{"static"},
				},
			}
			Expect(k8sClient.Create(ctx, microFrontend)).Should(Not(Succeed()))
		})

		It("Should create with defaults if optional fields are not specified", func() {
			By("By creating a new MicroFrontend with only required fields")
			ctx := context.Background()

			portNumber := int32(8080)

			microFrontend := &polyfeav1alpha1.MicroFrontend{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "polyfea.github.io/v1alpha1",
					Kind:       "MicroFrontend",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MicroFrontendName,
					Namespace: MicroFrontendNamespace,
				},
				Spec: polyfeav1alpha1.MicroFrontendSpec{
					Service: &polyfeav1alpha1.ServiceReference{
						Name: "test-service",
						Port: &polyfeav1alpha1.Port{
							Number: &portNumber,
						},
					},
					ModulePath:    &[]string{"module.jsm"}[0],
					FrontendClass: &[]string{"test-microfrontendclass"}[0],
				},
			}
			Expect(k8sClient.Create(ctx, microFrontend)).Should(Succeed())

			microFrontendLookupKey := types.NamespacedName{Name: MicroFrontendName, Namespace: MicroFrontendNamespace}
			createdMicroFrontend := &polyfeav1alpha1.MicroFrontend{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendLookupKey, createdMicroFrontend)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontend.Spec.Service.Name).Should(Equal("test-service"))
			Expect(*createdMicroFrontend.Spec.Proxy).Should(Equal(true))
			Expect(createdMicroFrontend.Spec.CacheStrategy).Should(Equal("none"))
			Expect(createdMicroFrontend.Spec.CacheControl).Should(BeNil())
			Expect(createdMicroFrontend.Spec.StaticPaths).Should(BeNil())
			Expect(createdMicroFrontend.Spec.PreloadPaths).Should(BeNil())
			Expect(createdMicroFrontend.Spec.DependsOn).Should(BeNil())
		})
	})

})
