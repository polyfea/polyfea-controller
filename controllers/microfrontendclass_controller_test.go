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
					BaseUri:   "http://localhost:8080",
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
			Expect(createdMicroFrontendClass.Spec.BaseUri).Should(Equal("http://localhost:8080"))

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
	})
})
