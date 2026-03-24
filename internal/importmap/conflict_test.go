package importmap

import (
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	polyfeav1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
)

func ptr[T any](v T) *T { return &v }

func makeMF(namespace, name, frontendClass string, createdAt time.Time, imports map[string]string, scopes map[string]map[string]string, accepted bool) polyfeav1alpha1.MicroFrontend {
	mf := polyfeav1alpha1.MicroFrontend{
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			Namespace:         namespace,
			CreationTimestamp: metav1.NewTime(createdAt),
		},
		Spec: polyfeav1alpha1.MicroFrontendSpec{
			FrontendClass: ptr(frontendClass),
		},
	}
	if imports != nil || scopes != nil {
		mf.Spec.ImportMap = &polyfeav1alpha1.ImportMap{
			Imports: imports,
			Scopes:  scopes,
		}
	}
	if accepted {
		mf.Status.FrontendClassRef = &polyfeav1alpha1.MicroFrontendClassReference{Accepted: true}
	}
	return mf
}

var t0 = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

func TestDetectConflicts(t *testing.T) {
	tests := []struct {
		name              string
		mf                polyfeav1alpha1.MicroFrontend
		others            []polyfeav1alpha1.MicroFrontend
		frontendClassName string
		wantConflicts     int
	}{
		{
			name:              "no conflicts when import maps don't overlap",
			mf:                makeMF("ns", "mf-a", "cls", t0.Add(time.Hour), map[string]string{"react": "./react.js"}, nil, true),
			others:            []polyfeav1alpha1.MicroFrontend{makeMF("ns", "mf-b", "cls", t0, map[string]string{"vue": "./vue.js"}, nil, true)},
			frontendClassName: "cls",
			wantConflicts:     0,
		},
		{
			name:              "conflict when same specifier maps to different path",
			mf:                makeMF("ns", "mf-a", "cls", t0.Add(time.Hour), map[string]string{"react": "./react-v17.js"}, nil, true),
			others:            []polyfeav1alpha1.MicroFrontend{makeMF("ns", "mf-b", "cls", t0, map[string]string{"react": "./react-v18.js"}, nil, true)},
			frontendClassName: "cls",
			wantConflicts:     1,
		},
		{
			name:              "no conflict when same specifier maps to same path",
			mf:                makeMF("ns", "mf-a", "cls", t0.Add(time.Hour), map[string]string{"react": "./react.js"}, nil, true),
			others:            []polyfeav1alpha1.MicroFrontend{makeMF("ns", "mf-b", "cls", t0, map[string]string{"react": "./react.js"}, nil, true)},
			frontendClassName: "cls",
			wantConflicts:     0,
		},
		{
			name:              "first-registered wins: older MF has no conflict",
			mf:                makeMF("ns", "mf-a", "cls", t0, map[string]string{"react": "./react-v17.js"}, nil, true),
			others:            []polyfeav1alpha1.MicroFrontend{makeMF("ns", "mf-b", "cls", t0.Add(time.Hour), map[string]string{"react": "./react-v18.js"}, nil, true)},
			frontendClassName: "cls",
			wantConflicts:     0,
		},
		{
			name:              "skips self",
			mf:                makeMF("ns", "mf-a", "cls", t0, map[string]string{"react": "./react.js"}, nil, true),
			others:            []polyfeav1alpha1.MicroFrontend{makeMF("ns", "mf-a", "cls", t0, map[string]string{"react": "./other.js"}, nil, true)},
			frontendClassName: "cls",
			wantConflicts:     0,
		},
		{
			name:              "skips MFs in different FrontendClass",
			mf:                makeMF("ns", "mf-a", "cls-a", t0.Add(time.Hour), map[string]string{"react": "./react-v17.js"}, nil, true),
			others:            []polyfeav1alpha1.MicroFrontend{makeMF("ns", "mf-b", "cls-b", t0, map[string]string{"react": "./react-v18.js"}, nil, true)},
			frontendClassName: "cls-a",
			wantConflicts:     0,
		},
		{
			name:              "skips unaccepted MFs",
			mf:                makeMF("ns", "mf-a", "cls", t0.Add(time.Hour), map[string]string{"react": "./react-v17.js"}, nil, true),
			others:            []polyfeav1alpha1.MicroFrontend{makeMF("ns", "mf-b", "cls", t0, map[string]string{"react": "./react-v18.js"}, nil, false)},
			frontendClassName: "cls",
			wantConflicts:     0,
		},
		{
			name:              "skips MFs with no import map",
			mf:                makeMF("ns", "mf-a", "cls", t0.Add(time.Hour), map[string]string{"react": "./react.js"}, nil, true),
			others:            []polyfeav1alpha1.MicroFrontend{makeMF("ns", "mf-b", "cls", t0, nil, nil, true)},
			frontendClassName: "cls",
			wantConflicts:     0,
		},
		{
			name: "scoped import conflicts detected",
			mf: makeMF("ns", "mf-a", "cls", t0.Add(time.Hour), nil,
				map[string]map[string]string{"/app/": {"lodash": "./lodash-v3.js"}}, true),
			others: []polyfeav1alpha1.MicroFrontend{
				makeMF("ns", "mf-b", "cls", t0, nil,
					map[string]map[string]string{"/app/": {"lodash": "./lodash-v4.js"}}, true),
			},
			frontendClassName: "cls",
			wantConflicts:     1,
		},
		{
			name: "mixed top-level and scoped conflicts",
			mf: makeMF("ns", "mf-a", "cls", t0.Add(time.Hour),
				map[string]string{"react": "./react-old.js"},
				map[string]map[string]string{"/app/": {"vue": "./vue-old.js"}}, true),
			others: []polyfeav1alpha1.MicroFrontend{
				makeMF("ns", "mf-b", "cls", t0,
					map[string]string{"react": "./react-new.js"},
					map[string]map[string]string{"/app/": {"vue": "./vue-new.js"}}, true),
			},
			frontendClassName: "cls",
			wantConflicts:     2,
		},
		{
			name: "uses default FrontendClassName when spec is nil",
			mf: func() polyfeav1alpha1.MicroFrontend {
				mf := makeMF("ns", "mf-a", "", t0.Add(time.Hour), map[string]string{"react": "./react-old.js"}, nil, true)
				mf.Spec.FrontendClass = nil
				return mf
			}(),
			others: []polyfeav1alpha1.MicroFrontend{
				func() polyfeav1alpha1.MicroFrontend {
					mf := makeMF("ns", "mf-b", "", t0, map[string]string{"react": "./react-new.js"}, nil, true)
					mf.Spec.FrontendClass = nil
					return mf
				}(),
			},
			frontendClassName: DefaultFrontendClassName,
			wantConflicts:     1,
		},
		{
			name: "uses default FrontendClassName when spec is empty string",
			mf:   makeMF("ns", "mf-a", "", t0.Add(time.Hour), map[string]string{"react": "./react-old.js"}, nil, true),
			others: []polyfeav1alpha1.MicroFrontend{
				makeMF("ns", "mf-b", "", t0, map[string]string{"react": "./react-new.js"}, nil, true),
			},
			frontendClassName: DefaultFrontendClassName,
			wantConflicts:     1,
		},
		{
			name:              "no others means no conflicts",
			mf:                makeMF("ns", "mf-a", "cls", t0, map[string]string{"react": "./react.js"}, nil, true),
			others:            nil,
			frontendClassName: "cls",
			wantConflicts:     0,
		},
		{
			name:              "nil import map on target means no conflicts",
			mf:                makeMF("ns", "mf-a", "cls", t0, nil, nil, true),
			others:            []polyfeav1alpha1.MicroFrontend{makeMF("ns", "mf-b", "cls", t0, map[string]string{"react": "./react.js"}, nil, true)},
			frontendClassName: "cls",
			wantConflicts:     0,
		},
		{
			name: "multiple others with conflicts from different MFs",
			mf:   makeMF("ns", "mf-c", "cls", t0.Add(2*time.Hour), map[string]string{"react": "./react-old.js", "vue": "./vue-old.js"}, nil, true),
			others: []polyfeav1alpha1.MicroFrontend{
				makeMF("ns", "mf-a", "cls", t0, map[string]string{"react": "./react-new.js"}, nil, true),
				makeMF("ns", "mf-b", "cls", t0.Add(time.Hour), map[string]string{"vue": "./vue-new.js"}, nil, true),
			},
			frontendClassName: "cls",
			wantConflicts:     2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conflicts := DetectConflicts(&tt.mf, tt.others, tt.frontendClassName)
			if len(conflicts) != tt.wantConflicts {
				t.Errorf("DetectConflicts() returned %d conflicts, want %d. Conflicts: %+v", len(conflicts), tt.wantConflicts, conflicts)
			}
		})
	}
}

func TestDetectConflictsDetails(t *testing.T) {
	mf := makeMF("ns", "mf-a", "cls", t0.Add(time.Hour), map[string]string{"react": "./react-v17.js"}, nil, true)
	others := []polyfeav1alpha1.MicroFrontend{
		makeMF("other-ns", "mf-b", "cls", t0, map[string]string{"react": "./react-v18.js"}, nil, true),
	}

	conflicts := DetectConflicts(&mf, others, "cls")
	if len(conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(conflicts))
	}

	c := conflicts[0]
	if c.Specifier != "react" {
		t.Errorf("Specifier = %q, want %q", c.Specifier, "react")
	}
	if c.RequestedPath != "./react-v17.js" {
		t.Errorf("RequestedPath = %q, want %q", c.RequestedPath, "./react-v17.js")
	}
	if c.RegisteredPath != "./react-v18.js" {
		t.Errorf("RegisteredPath = %q, want %q", c.RegisteredPath, "./react-v18.js")
	}
	if c.RegisteredBy != "other-ns/mf-b" {
		t.Errorf("RegisteredBy = %q, want %q", c.RegisteredBy, "other-ns/mf-b")
	}
	if c.Scope != "" {
		t.Errorf("Scope = %q, want empty", c.Scope)
	}
}

func TestDetectConflictsScopedDetails(t *testing.T) {
	mf := makeMF("ns", "mf-a", "cls", t0.Add(time.Hour), nil,
		map[string]map[string]string{"/app/": {"lodash": "./lodash-v3.js"}}, true)
	others := []polyfeav1alpha1.MicroFrontend{
		makeMF("ns", "mf-b", "cls", t0, nil,
			map[string]map[string]string{"/app/": {"lodash": "./lodash-v4.js"}}, true),
	}

	conflicts := DetectConflicts(&mf, others, "cls")
	if len(conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(conflicts))
	}

	c := conflicts[0]
	if c.Scope != "/app/" {
		t.Errorf("Scope = %q, want %q", c.Scope, "/app/")
	}
	if c.Specifier != "lodash" {
		t.Errorf("Specifier = %q, want %q", c.Specifier, "lodash")
	}
}

func TestConflictsEqual(t *testing.T) {
	tests := []struct {
		name string
		a, b []polyfeav1alpha1.ImportMapConflict
		want bool
	}{
		{
			name: "both empty",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			name: "both empty slices",
			a:    []polyfeav1alpha1.ImportMapConflict{},
			b:    []polyfeav1alpha1.ImportMapConflict{},
			want: true,
		},
		{
			name: "equal single element",
			a: []polyfeav1alpha1.ImportMapConflict{
				{Specifier: "react", RequestedPath: "./a.js", RegisteredPath: "./b.js", RegisteredBy: "ns/mf"},
			},
			b: []polyfeav1alpha1.ImportMapConflict{
				{Specifier: "react", RequestedPath: "./a.js", RegisteredPath: "./b.js", RegisteredBy: "ns/mf"},
			},
			want: true,
		},
		{
			name: "equal different order",
			a: []polyfeav1alpha1.ImportMapConflict{
				{Specifier: "react", RequestedPath: "./a.js", RegisteredPath: "./b.js", RegisteredBy: "ns/mf1"},
				{Specifier: "vue", RequestedPath: "./c.js", RegisteredPath: "./d.js", RegisteredBy: "ns/mf2"},
			},
			b: []polyfeav1alpha1.ImportMapConflict{
				{Specifier: "vue", RequestedPath: "./c.js", RegisteredPath: "./d.js", RegisteredBy: "ns/mf2"},
				{Specifier: "react", RequestedPath: "./a.js", RegisteredPath: "./b.js", RegisteredBy: "ns/mf1"},
			},
			want: true,
		},
		{
			name: "unequal different lengths",
			a: []polyfeav1alpha1.ImportMapConflict{
				{Specifier: "react", RequestedPath: "./a.js", RegisteredPath: "./b.js", RegisteredBy: "ns/mf"},
			},
			b:    nil,
			want: false,
		},
		{
			name: "unequal same keys different RequestedPath",
			a: []polyfeav1alpha1.ImportMapConflict{
				{Specifier: "react", RequestedPath: "./a.js", RegisteredPath: "./b.js", RegisteredBy: "ns/mf"},
			},
			b: []polyfeav1alpha1.ImportMapConflict{
				{Specifier: "react", RequestedPath: "./DIFFERENT.js", RegisteredPath: "./b.js", RegisteredBy: "ns/mf"},
			},
			want: false,
		},
		{
			name: "unequal same keys different RegisteredBy",
			a: []polyfeav1alpha1.ImportMapConflict{
				{Specifier: "react", RequestedPath: "./a.js", RegisteredPath: "./b.js", RegisteredBy: "ns/mf1"},
			},
			b: []polyfeav1alpha1.ImportMapConflict{
				{Specifier: "react", RequestedPath: "./a.js", RegisteredPath: "./b.js", RegisteredBy: "ns/mf2"},
			},
			want: false,
		},
		{
			name: "equal with scoped conflicts",
			a: []polyfeav1alpha1.ImportMapConflict{
				{Specifier: "lodash", Scope: "/app/", RequestedPath: "./a.js", RegisteredPath: "./b.js", RegisteredBy: "ns/mf"},
			},
			b: []polyfeav1alpha1.ImportMapConflict{
				{Specifier: "lodash", Scope: "/app/", RequestedPath: "./a.js", RegisteredPath: "./b.js", RegisteredBy: "ns/mf"},
			},
			want: true,
		},
		{
			name: "unequal different scopes",
			a: []polyfeav1alpha1.ImportMapConflict{
				{Specifier: "lodash", Scope: "/app/", RequestedPath: "./a.js", RegisteredPath: "./b.js", RegisteredBy: "ns/mf"},
			},
			b: []polyfeav1alpha1.ImportMapConflict{
				{Specifier: "lodash", Scope: "/other/", RequestedPath: "./a.js", RegisteredPath: "./b.js", RegisteredBy: "ns/mf"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConflictsEqual(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("ConflictsEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
