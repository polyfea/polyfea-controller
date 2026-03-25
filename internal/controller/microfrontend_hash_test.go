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
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
)

const (
	testModulePath   = "/app.js"
	testETagV1       = `"v1"`
	testETagV2       = `"v2"`
	testExistingHash = "existinghash1"
)

// newTestReconciler creates a minimal MicroFrontendReconciler sufficient for hash tests.
func newTestReconciler() *MicroFrontendReconciler {
	return &MicroFrontendReconciler{}
}

// newMFWithURL creates a MicroFrontend whose ResolvedServiceURL points to the given server.
func newMFWithURL(serviceURL string) *v1alpha1.MicroFrontend {
	modulePath := testModulePath
	return &v1alpha1.MicroFrontend{
		ObjectMeta: metav1.ObjectMeta{Name: "test-mf", Namespace: "default"},
		Spec: v1alpha1.MicroFrontendSpec{
			ModulePath: &modulePath,
		},
		Status: v1alpha1.MicroFrontendStatus{
			ResolvedServiceURL: serviceURL,
		},
	}
}

// expectedHash returns the first 12 hex characters of SHA256(content).
func expectedHash(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])[:12]
}

func TestResolveModuleHash_SetsHashOnFirstFetch(t *testing.T) {
	body := []byte("console.log('hello');")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", testETagV1)
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	r := newTestReconciler()
	mf := newMFWithURL(srv.URL)

	changed := r.resolveModuleHash(context.Background(), mf)

	if !changed {
		t.Fatal("expected changed=true on first fetch")
	}
	want := expectedHash(body)
	if mf.Status.ModuleHash != want {
		t.Errorf("ModuleHash = %q, want %q", mf.Status.ModuleHash, want)
	}
	if mf.Status.ModuleETag != testETagV1 {
		t.Errorf("ModuleETag = %q, want %q", mf.Status.ModuleETag, testETagV1)
	}
}

func TestResolveModuleHash_NoChangeWhenContentUnchanged(t *testing.T) {
	body := []byte("console.log('hello');")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", testETagV1)
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	r := newTestReconciler()
	mf := newMFWithURL(srv.URL)
	// Pre-populate with the same hash as the server would produce
	mf.Status.ModuleHash = expectedHash(body)
	mf.Status.ModuleETag = testETagV1

	changed := r.resolveModuleHash(context.Background(), mf)

	if changed {
		t.Fatal("expected changed=false when hash and ETag are already current")
	}
}

func TestResolveModuleHash_UpdatesHashWhenContentChanges(t *testing.T) {
	oldBody := []byte("console.log('old');")
	newBody := []byte("console.log('new');")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", testETagV2)
		_, _ = w.Write(newBody)
	}))
	defer srv.Close()

	r := newTestReconciler()
	mf := newMFWithURL(srv.URL)
	mf.Status.ModuleHash = expectedHash(oldBody)
	mf.Status.ModuleETag = testETagV1 // Different ETag → no 304, full response

	changed := r.resolveModuleHash(context.Background(), mf)

	if !changed {
		t.Fatal("expected changed=true when content changed")
	}
	want := expectedHash(newBody)
	if mf.Status.ModuleHash != want {
		t.Errorf("ModuleHash = %q, want %q", mf.Status.ModuleHash, want)
	}
	if mf.Status.ModuleETag != testETagV2 {
		t.Errorf("ModuleETag = %q, want %q", mf.Status.ModuleETag, testETagV2)
	}
}

func TestResolveModuleHash_ETagHit_Returns304(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("If-None-Match") == testETagV1 {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("ETag", testETagV1)
		_, _ = w.Write([]byte("content"))
	}))
	defer srv.Close()

	r := newTestReconciler()
	mf := newMFWithURL(srv.URL)
	mf.Status.ModuleHash = testExistingHash
	mf.Status.ModuleETag = testETagV1

	changed := r.resolveModuleHash(context.Background(), mf)

	if changed {
		t.Fatal("expected changed=false on 304 Not Modified")
	}
	// Hash must be preserved unchanged
	if mf.Status.ModuleHash != testExistingHash {
		t.Errorf("ModuleHash changed unexpectedly to %q", mf.Status.ModuleHash)
	}
}

func TestResolveModuleHash_KeepsExistingHashOnFetchFailure(t *testing.T) {
	// Point to a port with no listener
	r := newTestReconciler()
	mf := newMFWithURL("http://127.0.0.1:1")
	mf.Status.ModuleHash = testExistingHash
	mf.Status.ModuleETag = testETagV1

	changed := r.resolveModuleHash(context.Background(), mf)

	if changed {
		t.Fatal("expected changed=false on fetch failure")
	}
	if mf.Status.ModuleHash != testExistingHash {
		t.Errorf("ModuleHash changed on failure: got %q", mf.Status.ModuleHash)
	}
}

func TestResolveModuleHash_KeepsExistingHashOnNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	r := newTestReconciler()
	mf := newMFWithURL(srv.URL)
	mf.Status.ModuleHash = testExistingHash

	changed := r.resolveModuleHash(context.Background(), mf)

	if changed {
		t.Fatal("expected changed=false on non-200 response")
	}
	if mf.Status.ModuleHash != testExistingHash {
		t.Errorf("ModuleHash changed on error response: got %q", mf.Status.ModuleHash)
	}
}

func TestResolveModuleHash_SkipsWhenNoServiceURL(t *testing.T) {
	r := newTestReconciler()
	mf := newMFWithURL("")

	changed := r.resolveModuleHash(context.Background(), mf)

	if changed {
		t.Fatal("expected changed=false when ResolvedServiceURL is empty")
	}
}

func TestResolveModuleHash_SkipsWhenNoModulePath(t *testing.T) {
	r := newTestReconciler()
	mf := &v1alpha1.MicroFrontend{
		Status: v1alpha1.MicroFrontendStatus{ResolvedServiceURL: "http://example.com"},
	}

	changed := r.resolveModuleHash(context.Background(), mf)

	if changed {
		t.Fatal("expected changed=false when ModulePath is nil")
	}
}

func TestResolveModuleHash_HashIs12HexChars(t *testing.T) {
	body := []byte("module content")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	r := newTestReconciler()
	mf := newMFWithURL(srv.URL)

	r.resolveModuleHash(context.Background(), mf)

	if len(mf.Status.ModuleHash) != 12 {
		t.Errorf("ModuleHash length = %d, want 12; got %q", len(mf.Status.ModuleHash), mf.Status.ModuleHash)
	}
	for _, c := range mf.Status.ModuleHash {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			t.Errorf("ModuleHash contains non-hex char %q in %q", c, mf.Status.ModuleHash)
		}
	}
}
