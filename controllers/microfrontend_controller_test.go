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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	polyfeav1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	MicroFrontendName      = "test-microfrontend"
	MicroFrontendNamespace = "default"
	MicroFrontendFinalizer = "polyfea.github.io/finalizer"
)

func createMicroFrontend(name string, service, modulePath *string, proxy *bool, frontendClass *string, staticResources []polyfeav1alpha1.StaticResources, cacheOptions *polyfeav1alpha1.PWACache) *polyfeav1alpha1.MicroFrontend {
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

var _ = Describe("Microfrontend controller", func() {

	Context("When creating a MicroFrontend", func() {
		It("Should add and remove finalizer", func() {
			By("Creating a new MicroFrontend")
			ctx := context.Background()
			proxy := true
			microFrontend := createMicroFrontend(
				MicroFrontendName,
				ptr("http://test-service.test-namespace.svc.cluster.local"),
				ptr("module.jsm"),
				&proxy,
				ptr("test-microfrontendclass"),
				[]polyfeav1alpha1.StaticResources{{Path: "static", Kind: "script"}},
				nil,
			)
			Expect(k8sClient.Create(ctx, microFrontend)).Should(Succeed())

			microFrontendLookupKey := types.NamespacedName{Name: MicroFrontendName, Namespace: MicroFrontendNamespace}
			createdMicroFrontend := &polyfeav1alpha1.MicroFrontend{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendLookupKey, createdMicroFrontend)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontend.Spec.Service).Should(Equal(ptr("http://test-service.test-namespace.svc.cluster.local")))

			By("Checking the MicroFrontend has finalizer")
			Eventually(func() ([]string, error) {
				err := k8sClient.Get(ctx, microFrontendLookupKey, createdMicroFrontend)
				if err != nil {
					return []string{}, err
				}
				return createdMicroFrontend.ObjectMeta.GetFinalizers(), nil
			}, timeout, interval).Should(ContainElement(MicroFrontendFinalizer))

			By("Deleting the MicroFrontend")
			Expect(k8sClient.Delete(ctx, createdMicroFrontend)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendLookupKey, createdMicroFrontend)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())
		})

		DescribeTable("Should not create if required fields are missing",
			func(service, modulePath *string) {
				ctx := context.Background()
				proxy := true
				microFrontend := createMicroFrontend(
					MicroFrontendName,
					service,
					modulePath,
					&proxy,
					ptr("test-microfrontendclass"),
					[]polyfeav1alpha1.StaticResources{{Path: "static", Kind: "script"}},
					nil,
				)
				Expect(k8sClient.Create(ctx, microFrontend)).Should(Not(Succeed()))
			},
			Entry("missing service", nil, ptr("module.jsm")),
			Entry("missing modulePath", ptr("http://test-service.test-namespace.svc.cluster.local"), nil),
		)

		It("Should create with defaults if optional fields are not specified", func() {
			By("Creating a new MicroFrontend with only required fields")
			ctx := context.Background()
			microFrontend := createMicroFrontend(
				MicroFrontendName,
				ptr("http://test-service.test-namespace.svc.cluster.local"),
				ptr("module.jsm"),
				nil,
				nil,
				nil,
				nil,
			)
			Expect(k8sClient.Create(ctx, microFrontend)).Should(Succeed())

			microFrontendLookupKey := types.NamespacedName{Name: MicroFrontendName, Namespace: MicroFrontendNamespace}
			createdMicroFrontend := &polyfeav1alpha1.MicroFrontend{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendLookupKey, createdMicroFrontend)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontend.Spec.Service).Should(Equal(ptr("http://test-service.test-namespace.svc.cluster.local")))
			Expect(*createdMicroFrontend.Spec.Proxy).Should(Equal(true))
			Expect(createdMicroFrontend.Spec.CacheStrategy).Should(Equal("none"))
			Expect(createdMicroFrontend.Spec.FrontendClass).Should(Equal(ptr("polyfea-controller-default")))
			Expect(createdMicroFrontend.Spec.CacheControl).Should(BeNil())
			Expect(createdMicroFrontend.Spec.StaticResources).Should(BeNil())
			Expect(createdMicroFrontend.Spec.DependsOn).Should(BeNil())

			Expect(k8sClient.Delete(ctx, createdMicroFrontend)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendLookupKey, createdMicroFrontend)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())
		})

		It("Should add and remove the microfrontend from repository", func() {
			By("Creating a new MicroFrontend")
			ctx := context.Background()
			proxy := true
			microFrontend := createMicroFrontend(
				MicroFrontendName,
				ptr("http://test-service.test-namespace.svc.cluster.local"),
				ptr("module.jsm"),
				&proxy,
				ptr("test-microfrontendclass"),
				[]polyfeav1alpha1.StaticResources{{Path: "static", Kind: "script"}},
				nil,
			)
			Expect(k8sClient.Create(ctx, microFrontend)).Should(Succeed())

			microFrontendLookupKey := types.NamespacedName{Name: MicroFrontendName, Namespace: MicroFrontendNamespace}
			createdMicroFrontend := &polyfeav1alpha1.MicroFrontend{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendLookupKey, createdMicroFrontend)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Expect(createdMicroFrontend.Spec.Service).Should(Equal(ptr("http://test-service.test-namespace.svc.cluster.local")))

			Eventually(func() []*polyfeav1alpha1.MicroFrontend {
				result, _ := microFrontendRepository.List(func(mf *polyfeav1alpha1.MicroFrontend) bool {
					return mf.Name == MicroFrontendName
				})
				return result
			}, timeout, interval).Should(HaveLen(1))

			By("Deleting the MicroFrontend")
			Expect(k8sClient.Delete(ctx, createdMicroFrontend)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, microFrontendLookupKey, createdMicroFrontend)
				return errors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			Eventually(func() []*polyfeav1alpha1.MicroFrontend {
				result, _ := microFrontendRepository.List(func(mf *polyfeav1alpha1.MicroFrontend) bool {
					return mf.Name == MicroFrontendName
				})
				return result
			}, timeout, interval).Should(HaveLen(0))
		})

		DescribeTable("PWA PreCache validation",
			func(preCache []polyfeav1alpha1.PreCacheEntry, shouldSucceed bool) {
				ctx := context.Background()
				proxy := true
				microFrontend := createMicroFrontend(
					MicroFrontendName,
					ptr("http://test-service.test-namespace.svc.cluster.local"),
					ptr("module.jsm"),
					&proxy,
					ptr("test-microfrontendclass"),
					[]polyfeav1alpha1.StaticResources{{Path: "static", Kind: "script"}},
					&polyfeav1alpha1.PWACache{PreCache: preCache},
				)
				if shouldSucceed {
					Expect(k8sClient.Create(ctx, microFrontend)).Should(Succeed())
				} else {
					Expect(k8sClient.Create(ctx, microFrontend)).Should(Not(Succeed()))
				}
			},
			Entry("missing URL in PreCache", []polyfeav1alpha1.PreCacheEntry{{Revision: ptr("1")}}, false),
			Entry("valid PreCache", []polyfeav1alpha1.PreCacheEntry{{URL: ptr("1")}}, true),
		)

		DescribeTable("PWA RuntimeCache validation",
			func(cacheRoutes []polyfeav1alpha1.CacheRoute, shouldSucceed bool) {
				ctx := context.Background()
				proxy := true
				microFrontend := createMicroFrontend(
					MicroFrontendName,
					ptr("http://test-service.test-namespace.svc.cluster.local"),
					ptr("module.jsm"),
					&proxy,
					ptr("test-microfrontendclass"),
					[]polyfeav1alpha1.StaticResources{{Path: "static", Kind: "script"}},
					&polyfeav1alpha1.PWACache{CacheRoutes: cacheRoutes},
				)
				if shouldSucceed {
					Expect(k8sClient.Create(ctx, microFrontend)).Should(Succeed())
				} else {
					Expect(k8sClient.Create(ctx, microFrontend)).Should(Not(Succeed()))
				}
			},
			Entry("missing Pattern in CacheRoute", []polyfeav1alpha1.CacheRoute{{Strategy: ptr("network-first")}}, false),
			Entry("valid CacheRoute", []polyfeav1alpha1.CacheRoute{{Pattern: ptr("1")}}, true),
		)
	})
})
