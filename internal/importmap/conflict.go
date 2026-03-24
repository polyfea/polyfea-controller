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

package importmap

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	polyfeav1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
)

const DefaultFrontendClassName = "polyfea-controller-default"

// entry tracks an import map entry's registration details.
type entry struct {
	path              string
	namespace         string
	name              string
	creationTimestamp metav1.Time
}

// DetectConflicts finds import map conflicts between the given MicroFrontend and other
// accepted MicroFrontends in the same FrontendClass. Uses first-registered-wins based
// on creation timestamp.
func DetectConflicts(
	mf *polyfeav1alpha1.MicroFrontend,
	mfList []polyfeav1alpha1.MicroFrontend,
	frontendClassName string,
) []polyfeav1alpha1.ImportMapConflict {
	existingImports, existingScopes := buildExistingImportMaps(mf, frontendClassName, mfList)
	return findConflicts(mf, existingImports, existingScopes)
}

// ConflictsEqual compares two slices of ImportMapConflict for equality.
func ConflictsEqual(a, b []polyfeav1alpha1.ImportMapConflict) bool {
	if len(a) != len(b) {
		return false
	}

	aMap := make(map[string]polyfeav1alpha1.ImportMapConflict)
	for _, conflict := range a {
		key := conflict.Scope + ":" + conflict.Specifier
		aMap[key] = conflict
	}

	for _, conflict := range b {
		key := conflict.Scope + ":" + conflict.Specifier
		existing, ok := aMap[key]
		if !ok ||
			existing.RequestedPath != conflict.RequestedPath ||
			existing.RegisteredPath != conflict.RegisteredPath ||
			existing.RegisteredBy != conflict.RegisteredBy {
			return false
		}
	}

	return true
}

// buildExistingImportMaps builds maps of existing import entries from other MicroFrontends.
func buildExistingImportMaps(
	mf *polyfeav1alpha1.MicroFrontend,
	frontendClassName string,
	mfList []polyfeav1alpha1.MicroFrontend,
) (map[string]*entry, map[string]map[string]*entry) {
	existingImports := make(map[string]*entry)
	existingScopes := make(map[string]map[string]*entry)

	for i := range mfList {
		other := &mfList[i]

		if !shouldProcess(mf, other, frontendClassName) {
			continue
		}

		processEntries(other, existingImports, existingScopes)
	}

	return existingImports, existingScopes
}

// shouldProcess checks if a MicroFrontend should be processed for import map conflicts.
func shouldProcess(
	mf *polyfeav1alpha1.MicroFrontend,
	other *polyfeav1alpha1.MicroFrontend,
	frontendClassName string,
) bool {
	// Skip self
	if other.Namespace == mf.Namespace && other.Name == mf.Name {
		return false
	}

	// Skip if not same frontend class
	otherClassName := DefaultFrontendClassName
	if other.Spec.FrontendClass != nil && *other.Spec.FrontendClass != "" {
		otherClassName = *other.Spec.FrontendClass
	}
	if otherClassName != frontendClassName {
		return false
	}

	// Skip if not accepted
	if other.Status.FrontendClassRef == nil || !other.Status.FrontendClassRef.Accepted {
		return false
	}

	// Skip if no import map
	if other.Spec.ImportMap == nil {
		return false
	}

	return true
}

// processEntries processes import map entries from a MicroFrontend.
func processEntries(
	mf *polyfeav1alpha1.MicroFrontend,
	existingImports map[string]*entry,
	existingScopes map[string]map[string]*entry,
) {
	// Process top-level imports
	for specifier, path := range mf.Spec.ImportMap.Imports {
		registerEntry(specifier, path, mf, existingImports)
	}

	// Process scoped imports
	for scope, scopeImports := range mf.Spec.ImportMap.Scopes {
		if existingScopes[scope] == nil {
			existingScopes[scope] = make(map[string]*entry)
		}
		for specifier, path := range scopeImports {
			registerEntry(specifier, path, mf, existingScopes[scope])
		}
	}
}

// registerEntry registers an import map entry, keeping the oldest one (first-registered wins).
func registerEntry(
	specifier string,
	path string,
	mf *polyfeav1alpha1.MicroFrontend,
	entries map[string]*entry,
) {
	e := &entry{
		path:              path,
		namespace:         mf.Namespace,
		name:              mf.Name,
		creationTimestamp: mf.CreationTimestamp,
	}

	if existing, ok := entries[specifier]; ok {
		// Keep the older one (first-registered)
		if mf.CreationTimestamp.Before(&existing.creationTimestamp) {
			entries[specifier] = e
		}
	} else {
		entries[specifier] = e
	}
}

// findConflicts finds conflicts between requested and existing import map entries.
func findConflicts(
	mf *polyfeav1alpha1.MicroFrontend,
	existingImports map[string]*entry,
	existingScopes map[string]map[string]*entry,
) []polyfeav1alpha1.ImportMapConflict {
	if mf.Spec.ImportMap == nil {
		return nil
	}

	var conflicts []polyfeav1alpha1.ImportMapConflict

	// Check top-level imports
	conflicts = checkConflicts(mf, mf.Spec.ImportMap.Imports, existingImports, "", conflicts)

	// Check scoped imports
	for scope, scopeImports := range mf.Spec.ImportMap.Scopes {
		if existingScopeMap, ok := existingScopes[scope]; ok {
			conflicts = checkConflicts(mf, scopeImports, existingScopeMap, scope, conflicts)
		}
	}

	return conflicts
}

// checkConflicts checks for conflicts in a set of imports.
func checkConflicts(
	mf *polyfeav1alpha1.MicroFrontend,
	imports map[string]string,
	existingEntries map[string]*entry,
	scope string,
	conflicts []polyfeav1alpha1.ImportMapConflict,
) []polyfeav1alpha1.ImportMapConflict {
	for specifier, requestedPath := range imports {
		if existing, ok := existingEntries[specifier]; ok {
			// Check if this MicroFrontend is older (registered first)
			if !mf.CreationTimestamp.Before(&existing.creationTimestamp) && requestedPath != existing.path {
				// Conflict: another MicroFrontend registered this specifier first with different path
				conflicts = append(conflicts, polyfeav1alpha1.ImportMapConflict{
					Specifier:      specifier,
					RequestedPath:  requestedPath,
					RegisteredPath: existing.path,
					RegisteredBy:   existing.namespace + "/" + existing.name,
					Scope:          scope,
				})
			}
		}
	}
	return conflicts
}
