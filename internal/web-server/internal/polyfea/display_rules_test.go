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

// TestSelectMatchingWebComponentsWildCombinations exercises deep nesting across all three
// operators (allOf / anyOf / noneOf) at multiple levels, plus multiple DisplayRules entries
// (OR semantics between them), and scalar matchers alongside nested operators at the same level.
//
// Rule set (two DisplayRules, OR between them):
//
// Rule 1: allOf
//   - anyOf: context=dashboard OR (context=reports AND role=analyst)  <- nested allOf inside anyOf
//   - noneOf: anyOf: path=^/draft OR role=guest                       <- noneOf containing anyOf
//   - path=^/app                                                       <- scalar alongside nested ops
//
// Reads: "in (dashboard OR (reports AND analyst)), not on /draft and not a guest, path starts /app"
//
// Rule 2: anyOf
//   - role=superadmin
//   - allOf: context=admin-panel AND noneOf: path=^/admin/secret       <- double nesting
//
// Reads: "superadmin anywhere, OR admin-panel but not the secret path"
func TestSelectMatchingWebComponentsWildCombinations(t *testing.T) {
	wc := &v1alpha1.WebComponent{
		Spec: v1alpha1.WebComponentSpec{
			DisplayRules: []v1alpha1.DisplayRules{
				// Rule 1
				{
					AllOf: []v1alpha1.Matcher{
						{
							// anyOf: dashboard  OR  (reports AND analyst)
							AnyOf: []v1alpha1.Matcher{
								{ContextName: "dashboard"},
								{
									// nested allOf inside anyOf
									AllOf: []v1alpha1.Matcher{
										{ContextName: "reports"},
										{Role: "analyst"},
									},
								},
							},
						},
						{
							// noneOf: (path /draft OR role guest)
							NoneOf: []v1alpha1.Matcher{
								{
									AnyOf: []v1alpha1.Matcher{
										{Path: "^/draft"},
										{Role: "guest"},
									},
								},
							},
						},
						// scalar alongside nested ops
						{Path: "^/app"},
					},
				},
				// Rule 2
				{
					AnyOf: []v1alpha1.Matcher{
						{Role: "superadmin"},
						{
							AllOf: []v1alpha1.Matcher{
								{ContextName: "admin-panel"},
								{
									NoneOf: []v1alpha1.Matcher{
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

	cases := []struct {
		desc    string
		context string
		path    string
		roles   []string
		want    bool
	}{
		// Rule 1 matches
		{"dashboard analyst /app", "dashboard", "/app/page", []string{"analyst"}, true},
		{"reports + analyst /app", "reports", "/app/data", []string{"analyst"}, true},
		// Rule 1 misses — wrong context
		{"other context", "other", "/app/page", []string{"analyst"}, false},
		// Rule 1 misses — reports without analyst role
		{"reports without analyst", "reports", "/app/data", []string{"viewer"}, false},
		// Rule 1 misses — excluded path (noneOf block)
		{"dashboard /draft blocked", "dashboard", "/draft/item", []string{"analyst"}, false},
		// Rule 1 misses — guest role blocked by noneOf
		{"dashboard guest blocked", "dashboard", "/app/page", []string{"guest"}, false},
		// Rule 1 misses — path doesn't match ^/app scalar
		{"dashboard wrong path", "dashboard", "/home", []string{"analyst"}, false},
		// Rule 2 matches — superadmin anywhere
		{"superadmin any context", "anything", "/wherever", []string{"superadmin"}, true},
		// Rule 2 matches — admin-panel non-secret path
		{"admin-panel safe path", "admin-panel", "/admin/users", []string{"viewer"}, true},
		// Rule 2 misses — admin-panel but secret path
		{"admin-panel secret blocked", "admin-panel", "/admin/secret/key", []string{"viewer"}, false},
		// Neither rule — no role, wrong context, wrong path
		{"no match at all", "public", "/home", []string{}, false},
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
