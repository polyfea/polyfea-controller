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

package v1alpha1

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const TestConst = "test"

func TestAttribute_DeepCopy(t *testing.T) {
	attr := &Attribute{
		Name:  TestConst,
		Value: runtime.RawExtension{Raw: []byte(`{"key":"value"}`)},
	}
	copied := attr.DeepCopy()
	if copied == nil {
		t.Error("DeepCopy returned nil")
	}

	// Test nil
	var nilAttr *Attribute
	if nilAttr.DeepCopy() != nil {
		t.Error("DeepCopy of nil should return nil")
	}
}

func TestCacheRoute_DeepCopy(t *testing.T) {
	cr := &CacheRoute{
		Pattern:              ptr("pattern"),
		Destination:          ptr("dest"),
		Strategy:             ptr("cache-first"),
		MaxAgeSeconds:        ptr(int32(3600)),
		SyncRetentionMinutes: ptr(int32(60)),
		Method:               ptr("GET"),
		Statuses:             []int32{200, 201},
	}
	copied := cr.DeepCopy()
	if copied == nil {
		t.Error("DeepCopy returned nil")
	}

	if cr.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestDisplayRules_DeepCopy(t *testing.T) {
	dr := &DisplayRules{
		AllOf:  []Matcher{{Path: "/test"}},
		AnyOf:  []Matcher{{Role: "admin"}},
		NoneOf: []Matcher{{ContextName: "context"}},
	}
	copied := dr.DeepCopy()
	if copied == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestHeader_DeepCopy(t *testing.T) {
	h := &Header{Name: "X-Test", Value: "value"}
	if h.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestMatcher_DeepCopy(t *testing.T) {
	m := &Matcher{Path: "/test"}
	if m.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestMetaTag_DeepCopy(t *testing.T) {
	mt := &MetaTag{Name: "description", Content: TestConst}
	if mt.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestMicroFrontend_DeepCopy(t *testing.T) {
	mf := &MicroFrontend{
		ObjectMeta: metav1.ObjectMeta{Name: TestConst},
		Spec:       MicroFrontendSpec{},
		Status:     MicroFrontendStatus{},
	}
	copied := mf.DeepCopy()
	if copied == nil {
		t.Error("DeepCopy returned nil")
	}

	obj := mf.DeepCopyObject()
	if obj == nil {
		t.Error("DeepCopyObject returned nil")
	}
}

func TestMicroFrontendClass_DeepCopy(t *testing.T) {
	mfc := &MicroFrontendClass{
		ObjectMeta: metav1.ObjectMeta{Name: TestConst},
		Spec:       MicroFrontendClassSpec{},
		Status:     MicroFrontendClassStatus{},
	}
	copied := mfc.DeepCopy()
	if copied == nil {
		t.Error("DeepCopy returned nil")
	}

	obj := mfc.DeepCopyObject()
	if obj == nil {
		t.Error("DeepCopyObject returned nil")
	}
}

func TestMicroFrontendClassList_DeepCopy(t *testing.T) {
	list := &MicroFrontendClassList{
		Items: []MicroFrontendClass{{ObjectMeta: metav1.ObjectMeta{Name: TestConst}}},
	}
	copied := list.DeepCopy()
	if copied == nil {
		t.Error("DeepCopy returned nil")
	}

	obj := list.DeepCopyObject()
	if obj == nil {
		t.Error("DeepCopyObject returned nil")
	}
}

func TestMicroFrontendClassReference_DeepCopy(t *testing.T) {
	ref := &MicroFrontendClassReference{Name: TestConst, Accepted: true}
	if ref.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestMicroFrontendClassSpec_DeepCopy(t *testing.T) {
	spec := &MicroFrontendClassSpec{
		BaseUri:           ptr("https://example.com"),
		Title:             ptr(TestConst),
		NamespacePolicy:   &NamespacePolicy{From: NamespaceFromAll},
		ExtraMetaTags:     []MetaTag{{Name: TestConst, Content: "value"}},
		ExtraHeaders:      []Header{{Name: "X-Test", Value: "value"}},
		ProgressiveWebApp: &ProgressiveWebApp{},
	}
	if spec.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestMicroFrontendClassStatus_DeepCopy(t *testing.T) {
	status := &MicroFrontendClassStatus{
		Conditions: []metav1.Condition{{Type: "Ready", Status: metav1.ConditionTrue}},
	}
	if status.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestMicroFrontendList_DeepCopy(t *testing.T) {
	list := &MicroFrontendList{
		Items: []MicroFrontend{{ObjectMeta: metav1.ObjectMeta{Name: TestConst}}},
	}
	copied := list.DeepCopy()
	if copied == nil {
		t.Error("DeepCopy returned nil")
	}

	obj := list.DeepCopyObject()
	if obj == nil {
		t.Error("DeepCopyObject returned nil")
	}
}

func TestMicroFrontendSpec_DeepCopy(t *testing.T) {
	spec := &MicroFrontendSpec{
		Service:         &ServiceReference{Name: ptr(TestConst)},
		Proxy:           ptr(true),
		CacheControl:    ptr("no-cache"),
		ModulePath:      ptr("/module.js"),
		StaticResources: []StaticResources{{Kind: "script", Path: "/test.js"}},
		FrontendClass:   ptr("default"),
		DependsOn:       []string{"dep1"},
		CacheOptions:    &PWACache{},
		ImportMap:       &ImportMap{},
	}
	if spec.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestMicroFrontendStatus_DeepCopy(t *testing.T) {
	status := &MicroFrontendStatus{
		Conditions:         []metav1.Condition{{Type: "Ready", Status: metav1.ConditionTrue}},
		FrontendClassRef:   &MicroFrontendClassReference{Name: TestConst},
		ImportMapConflicts: []ImportMapConflict{{Specifier: TestConst}},
	}
	if status.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestNamespacePolicy_DeepCopy(t *testing.T) {
	np := &NamespacePolicy{
		From:       NamespaceFromNamespaces,
		Namespaces: []string{"ns1", "ns2"},
	}
	if np.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestObjectReference_DeepCopy(t *testing.T) {
	ref := &ObjectReference{Name: TestConst, Namespace: "default", Found: true}
	if ref.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestPWACache_DeepCopy(t *testing.T) {
	cache := &PWACache{
		PreCache:    []PreCacheEntry{{URL: ptr("/test")}},
		CacheRoutes: []CacheRoute{{Pattern: ptr("/*")}},
	}
	if cache.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestPWAIcon_DeepCopy(t *testing.T) {
	icon := &PWAIcon{
		Sizes:   ptr("192x192"),
		Src:     ptr("/icon.png"),
		Type:    ptr("image/png"),
		Purpose: ptr("any"),
	}
	if icon.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestPort_DeepCopy(t *testing.T) {
	port := &Port{Name: "http", Number: ptr(int32(8080))}
	if port.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestPreCacheEntry_DeepCopy(t *testing.T) {
	entry := &PreCacheEntry{URL: ptr("/test"), Revision: ptr("v1")}
	if entry.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestProgressiveWebApp_DeepCopy(t *testing.T) {
	pwa := &ProgressiveWebApp{
		WebAppManifest:             &WebAppManifest{},
		CacheOptions:               &PWACache{},
		PolyfeaSWReconcileInterval: ptr(int32(1800000)),
	}
	if pwa.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestServiceReference_DeepCopy(t *testing.T) {
	sr := &ServiceReference{
		Name:      ptr("service"),
		URI:       ptr("https://example.com"),
		Namespace: ptr("default"),
		Port:      ptr(int32(80)),
		Scheme:    ptr("http"),
		Domain:    ptr("svc.cluster.local"),
	}
	if sr.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestStaticResources_DeepCopy(t *testing.T) {
	sr := &StaticResources{
		Kind:       "script",
		Path:       "/test.js",
		Attributes: []Attribute{{Name: "type", Value: runtime.RawExtension{Raw: []byte(`"module"`)}}},
		Proxy:      ptr(true),
	}
	if sr.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestStyle_DeepCopy(t *testing.T) {
	s := &Style{Name: "color", Value: "red"}
	if s.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestWebAppManifest_DeepCopy(t *testing.T) {
	wam := &WebAppManifest{
		Name:            ptr("App"),
		Icons:           []PWAIcon{{Src: ptr("/icon.png")}},
		StartUrl:        ptr("/"),
		Display:         ptr("standalone"),
		DisplayOverride: []string{"window-controls-overlay"},
	}
	if wam.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestWebComponent_DeepCopy(t *testing.T) {
	wc := &WebComponent{
		ObjectMeta: metav1.ObjectMeta{Name: TestConst},
		Spec:       WebComponentSpec{},
		Status:     WebComponentStatus{},
	}
	copied := wc.DeepCopy()
	if copied == nil {
		t.Error("DeepCopy returned nil")
	}

	obj := wc.DeepCopyObject()
	if obj == nil {
		t.Error("DeepCopyObject returned nil")
	}
}

func TestWebComponentList_DeepCopy(t *testing.T) {
	list := &WebComponentList{
		Items: []WebComponent{{ObjectMeta: metav1.ObjectMeta{Name: TestConst}}},
	}
	copied := list.DeepCopy()
	if copied == nil {
		t.Error("DeepCopy returned nil")
	}

	obj := list.DeepCopyObject()
	if obj == nil {
		t.Error("DeepCopyObject returned nil")
	}
}

func TestWebComponentSpec_DeepCopy(t *testing.T) {
	spec := &WebComponentSpec{
		MicroFrontend: ptr("mf"),
		Element:       ptr("my-element"),
		Attributes:    []Attribute{{Name: TestConst, Value: runtime.RawExtension{Raw: []byte(`"value"`)}}},
		DisplayRules:  []DisplayRules{{AllOf: []Matcher{{Path: "/"}}}},
		Priority:      ptr(int32(10)),
		Style:         []Style{{Name: "color", Value: "red"}},
	}
	if spec.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

func TestWebComponentStatus_DeepCopy(t *testing.T) {
	status := &WebComponentStatus{
		Conditions:       []metav1.Condition{{Type: "Ready", Status: metav1.ConditionTrue}},
		MicroFrontendRef: &ObjectReference{Name: "mf", Found: true},
	}
	if status.DeepCopy() == nil {
		t.Error("DeepCopy returned nil")
	}
}

// Test DeepCopyInto methods
func TestDeepCopyInto(t *testing.T) {
	t.Run("Attribute", func(t *testing.T) {
		in := &Attribute{Name: TestConst, Value: runtime.RawExtension{Raw: []byte(`"value"`)}}
		out := &Attribute{}
		in.DeepCopyInto(out)
		if out.Name != in.Name {
			t.Error("DeepCopyInto failed")
		}
	})

	t.Run("CacheRoute", func(t *testing.T) {
		in := &CacheRoute{Pattern: ptrString(TestConst)}
		out := &CacheRoute{}
		in.DeepCopyInto(out)
		if out.Pattern == nil || *out.Pattern != TestConst {
			t.Error("DeepCopyInto failed")
		}
	})

	t.Run("DisplayRules", func(t *testing.T) {
		in := &DisplayRules{AllOf: []Matcher{{Path: "/test"}}}
		out := &DisplayRules{}
		in.DeepCopyInto(out)
		if len(out.AllOf) != 1 {
			t.Error("DeepCopyInto failed")
		}
	})

	t.Run("MicroFrontend", func(t *testing.T) {
		in := &MicroFrontend{ObjectMeta: metav1.ObjectMeta{Name: TestConst}}
		out := &MicroFrontend{}
		in.DeepCopyInto(out)
		if out.Name != TestConst {
			t.Error("DeepCopyInto failed")
		}
	})

	t.Run("MicroFrontendClass", func(t *testing.T) {
		in := &MicroFrontendClass{ObjectMeta: metav1.ObjectMeta{Name: TestConst}}
		out := &MicroFrontendClass{}
		in.DeepCopyInto(out)
		if out.Name != TestConst {
			t.Error("DeepCopyInto failed")
		}
	})

	t.Run("WebComponent", func(t *testing.T) {
		in := &WebComponent{ObjectMeta: metav1.ObjectMeta{Name: TestConst}}
		out := &WebComponent{}
		in.DeepCopyInto(out)
		if out.Name != TestConst {
			t.Error("DeepCopyInto failed")
		}
	})
}

// Test DeepCopyObject for nil cases
func TestDeepCopyObject_Nil(t *testing.T) {
	t.Run("MicroFrontend", func(t *testing.T) {
		var mf *MicroFrontend
		if mf.DeepCopyObject() != nil {
			t.Error("DeepCopyObject of nil should return nil")
		}
	})

	t.Run("MicroFrontendClass", func(t *testing.T) {
		var mfc *MicroFrontendClass
		if mfc.DeepCopyObject() != nil {
			t.Error("DeepCopyObject of nil should return nil")
		}
	})

	t.Run("WebComponent", func(t *testing.T) {
		var wc *WebComponent
		if wc.DeepCopyObject() != nil {
			t.Error("DeepCopyObject of nil should return nil")
		}
	})

	t.Run("MicroFrontendList", func(t *testing.T) {
		var list *MicroFrontendList
		if list.DeepCopyObject() != nil {
			t.Error("DeepCopyObject of nil should return nil")
		}
	})

	t.Run("MicroFrontendClassList", func(t *testing.T) {
		var list *MicroFrontendClassList
		if list.DeepCopyObject() != nil {
			t.Error("DeepCopyObject of nil should return nil")
		}
	})

	t.Run("WebComponentList", func(t *testing.T) {
		var list *WebComponentList
		if list.DeepCopyObject() != nil {
			t.Error("DeepCopyObject of nil should return nil")
		}
	})
}

func ptrString(s string) *string {
	return &s
}

func TestImportMap_DeepCopy(t *testing.T) {
	tests := []struct {
		name string
		in   *ImportMap
	}{
		{
			name: "nil",
			in:   nil,
		},
		{
			name: "empty",
			in:   &ImportMap{},
		},
		{
			name: "with imports",
			in: &ImportMap{
				Imports: map[string]string{
					"angular": "./bundle/angular.mjs",
					"react":   "https://cdn.example.com/react.js",
				},
			},
		},
		{
			name: "with scopes",
			in: &ImportMap{
				Scopes: map[string]map[string]string{
					"/legacy/": {
						"angular": "./old-angular.mjs",
					},
				},
			},
		},
		{
			name: "with imports and scopes",
			in: &ImportMap{
				Imports: map[string]string{
					"angular": "./bundle/angular.mjs",
				},
				Scopes: map[string]map[string]string{
					"/legacy/": {
						"angular": "./old-angular.mjs",
					},
				},
			},
		},
		{
			name: "with nil scope value",
			in: &ImportMap{
				Scopes: map[string]map[string]string{
					"/nil/": nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copied := tt.in.DeepCopy()

			if tt.in == nil {
				if copied != nil {
					t.Errorf("DeepCopy() of nil should return nil")
				}
				return
			}

			if copied == nil {
				t.Errorf("DeepCopy() returned nil for non-nil input")
			}
		})
	}
}

func TestImportMap_DeepCopyInto(t *testing.T) {
	original := &ImportMap{
		Imports: map[string]string{
			TestConst: "./test.js",
		},
		Scopes: map[string]map[string]string{
			"/app/": {
				"lodash": "./lodash.js",
			},
		},
	}

	out := &ImportMap{}
	original.DeepCopyInto(out)

	if out.Imports == nil {
		t.Errorf("DeepCopyInto() did not copy imports")
	}
	if out.Scopes == nil {
		t.Errorf("DeepCopyInto() did not copy scopes")
	}
}

func TestImportMapConflict_DeepCopy(t *testing.T) {
	tests := []struct {
		name string
		in   *ImportMapConflict
	}{
		{
			name: "nil",
			in:   nil,
		},
		{
			name: "empty",
			in:   &ImportMapConflict{},
		},
		{
			name: "with all fields",
			in: &ImportMapConflict{
				Specifier:      "angular",
				RequestedPath:  "./bundle/angular-v15.mjs",
				RegisteredPath: "./bundle/angular-v17.mjs",
				RegisteredBy:   "namespace/mf-name",
				Scope:          "/legacy/",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copied := tt.in.DeepCopy()

			if tt.in == nil {
				if copied != nil {
					t.Errorf("DeepCopy() of nil should return nil")
				}
				return
			}

			if copied == nil {
				t.Errorf("DeepCopy() returned nil for non-nil input")
			}
		})
	}
}

func TestImportMapConflict_DeepCopyInto(t *testing.T) {
	original := &ImportMapConflict{
		Specifier:      "vue",
		RequestedPath:  "./vue2.js",
		RegisteredPath: "./vue3.js",
		RegisteredBy:   "team/app",
		Scope:          "/app/",
	}

	out := &ImportMapConflict{}
	original.DeepCopyInto(out)

	if out.Specifier != original.Specifier {
		t.Errorf("DeepCopyInto() did not copy Specifier")
	}
	if out.Scope != original.Scope {
		t.Errorf("DeepCopyInto() did not copy Scope")
	}
}
