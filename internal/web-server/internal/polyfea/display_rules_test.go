package polyfea

import (
	"net/http/httptest"
	"testing"

	"github.com/polyfea/polyfea-controller/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestConvertAttributesWithNonJSONValue(t *testing.T) {
	attrs := []v1alpha1.Attribute{
		{Name: "valid", Value: runtime.RawExtension{Raw: []byte(`"hello"`)}},
		{Name: "non-json", Value: runtime.RawExtension{Raw: []byte(`not-json`)}},
	}
	result := convertAttributes(attrs)
	if (*result)["valid"] != "hello" {
		t.Errorf("expected 'hello', got %q", (*result)["valid"])
	}
	if (*result)["non-json"] != "not-json" {
		t.Errorf("expected raw fallback 'not-json', got %q", (*result)["non-json"])
	}
}

func TestConvertStyles(t *testing.T) {
	styles := []v1alpha1.Style{
		{Name: "color", Value: "red"},
		{Name: "font-size", Value: "12px"},
	}
	result := convertStyles(styles)
	if (*result)["color"] != "red" {
		t.Errorf("expected 'red', got %q", (*result)["color"])
	}
	if (*result)["font-size"] != "12px" {
		t.Errorf("expected '12px', got %q", (*result)["font-size"])
	}
}

func TestAddExtraHeaders(t *testing.T) {
	headers := []v1alpha1.Header{
		{Name: "X-Custom", Value: "value1"},
		{Name: "X-Other", Value: "value2"},
	}
	w := httptest.NewRecorder()
	addExtraHeaders(w, headers)
	if w.Header().Get("X-Custom") != "value1" {
		t.Errorf("expected X-Custom=value1, got %q", w.Header().Get("X-Custom"))
	}
	if w.Header().Get("X-Other") != "value2" {
		t.Errorf("expected X-Other=value2, got %q", w.Header().Get("X-Other"))
	}
}

func TestSelectMatchingWebComponentsAllOf(t *testing.T) {
	wc := &v1alpha1.WebComponent{
		Spec: v1alpha1.WebComponentSpec{
			DisplayRules: []v1alpha1.DisplayRules{
				{AllOf: []v1alpha1.Matcher{{Path: "test.*", ContextName: "ctx"}}},
			},
		},
	}

	if !selectMatchingWebComponents(wc, "ctx", "test-path", nil) {
		t.Error("expected match")
	}
	if selectMatchingWebComponents(wc, "other", "test-path", nil) {
		t.Error("expected no match for wrong context name")
	}
}

// TestSelectMatchingWebComponentsNestedOperators verifies that allOf, anyOf, and noneOf
// can be used recursively inside another allOf entry to express combinatorial rules
// without listing every combination explicitly.
//
// The rule below reads:
//
//	"show if (context is navigation OR applications)
//	       AND (role is admin OR datascience)
//	       AND NOT path matches ^/secret$"
func TestSelectMatchingWebComponentsNestedOperators(t *testing.T) {
	wc := &v1alpha1.WebComponent{
		Spec: v1alpha1.WebComponentSpec{
			DisplayRules: []v1alpha1.DisplayRules{
				{
					AllOf: []v1alpha1.Matcher{
						{
							AnyOf: []v1alpha1.Matcher{
								{ContextName: "navigation"},
								{ContextName: "applications"},
							},
						},
						{
							AnyOf: []v1alpha1.Matcher{
								{Role: "admin"},
								{Role: "datascience"},
							},
						},
						{
							NoneOf: []v1alpha1.Matcher{
								{Path: "^/secret$"},
							},
						},
					},
				},
			},
		},
	}

	cases := []struct {
		desc    string
		context string
		path    string
		roles   []string
		want    bool
	}{
		{"navigation + admin", "navigation", "/page", []string{"admin"}, true},
		{"applications + datascience", "applications", "/page", []string{"datascience"}, true},
		{"navigation + datascience", "navigation", "/page", []string{"datascience"}, true},
		{"wrong context", "other", "/page", []string{"admin"}, false},
		{"wrong role", "navigation", "/page", []string{"other"}, false},
		{"excluded path", "navigation", "/secret", []string{"admin"}, false},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := selectMatchingWebComponents(wc, tc.context, tc.path, tc.roles)
			if got != tc.want {
				t.Errorf("selectMatchingWebComponents(%q, %q, %v) = %v, want %v",
					tc.context, tc.path, tc.roles, got, tc.want)
			}
		})
	}
}
