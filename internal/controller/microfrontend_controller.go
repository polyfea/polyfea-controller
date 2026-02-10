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
	"time"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	polyfeav1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/internal/repository"
)

// MicroFrontendReconciler reconciles a MicroFrontend object.
type MicroFrontendReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Repository repository.Repository[*polyfeav1alpha1.MicroFrontend]
}

// +kubebuilder:rbac:groups=polyfea.github.io,resources=microfrontends,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=polyfea.github.io,resources=microfrontends/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=polyfea.github.io,resources=microfrontends/finalizers,verbs=update

// Reconcile moves the current state of the cluster closer to the desired state.
func (r *MicroFrontendReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	const finalizerName = "polyfea.github.io/finalizer"
	logger := log.FromContext(ctx)

	mf := &polyfeav1alpha1.MicroFrontend{}
	if err := r.Get(ctx, req.NamespacedName, mf); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("MicroFrontend resource not found; assuming it was deleted.")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get MicroFrontend")
		return ctrl.Result{Requeue: true}, err
	}

	logger.Info("Reconciling MicroFrontend", "MicroFrontend", mf.Name, "Namespace", mf.Namespace)

	// Handle finalizer
	if result, err := r.handleFinalizer(ctx, mf, finalizerName); result != nil {
		return *result, err
	}

	// Handle deletion
	if result, err := r.handleDeletion(ctx, req, mf, finalizerName); result != nil {
		return *result, err
	}

	// Update status
	statusUpdated := false
	originalStatus := mf.Status.DeepCopy()

	// Resolve service URL
	statusUpdated = r.resolveServiceURL(mf) || statusUpdated

	// Process MicroFrontendClass and namespace policy
	statusUpdated = r.processFrontendClass(ctx, mf) || statusUpdated

	// Check for import map conflicts if accepted and has import map
	if mf.Status.FrontendClassRef != nil && mf.Status.FrontendClassRef.Accepted {
		statusUpdated = r.checkImportMapConflicts(ctx, mf) || statusUpdated
	}

	// Update ObservedGeneration
	if mf.Status.ObservedGeneration != mf.Generation {
		mf.Status.ObservedGeneration = mf.Generation
		statusUpdated = true
	}

	// Update status if needed
	if statusUpdated {
		if err := r.Status().Update(ctx, mf); err != nil {
			logger.Error(err, "Failed to update MicroFrontend status")
			mf.Status = *originalStatus
			return ctrl.Result{Requeue: true}, err
		}
		logger.Info("Updated MicroFrontend status", "phase", mf.Status.Phase)
	}

	// Periodic reconciliation to re-check import map conflicts and other conditions
	// This ensures conflicts are automatically resolved when a blocking MicroFrontend is deleted
	return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
}

// handleFinalizer adds finalizer if not present.
func (r *MicroFrontendReconciler) handleFinalizer(ctx context.Context, mf *polyfeav1alpha1.MicroFrontend, finalizerName string) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)

	if !controllerutil.ContainsFinalizer(mf, finalizerName) {
		logger.Info("Adding finalizer")
		controllerutil.AddFinalizer(mf, finalizerName)
		if err := r.Update(ctx, mf); err != nil {
			logger.Error(err, "Failed to update MicroFrontend with finalizer")
			return &ctrl.Result{Requeue: true}, err
		}
		return &ctrl.Result{}, nil
	}
	return nil, nil
}

// handleDeletion handles resource deletion and finalizer cleanup.
func (r *MicroFrontendReconciler) handleDeletion(ctx context.Context, req ctrl.Request, mf *polyfeav1alpha1.MicroFrontend, finalizerName string) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)

	if mf.GetDeletionTimestamp() == nil {
		return nil, nil
	}

	if controllerutil.ContainsFinalizer(mf, finalizerName) {
		logger.Info("Running finalizer operations before deletion")
		if err := r.finalizeMicroFrontend(mf); err != nil {
			logger.Error(err, "Finalizer operations failed")
			return &ctrl.Result{Requeue: true}, nil
		}
		// Re-fetch in case of update
		if err := r.Get(ctx, req.NamespacedName, mf); err != nil {
			logger.Error(err, "Failed to re-fetch MicroFrontend after finalizer")
			return &ctrl.Result{Requeue: true}, err
		}
		controllerutil.RemoveFinalizer(mf, finalizerName)
		if err := r.Update(ctx, mf); err != nil {
			logger.Error(err, "Failed to remove finalizer")
			return &ctrl.Result{Requeue: true}, err
		}
	}
	return &ctrl.Result{}, nil
}

// resolveServiceURL resolves the service URL and updates status.
func (r *MicroFrontendReconciler) resolveServiceURL(mf *polyfeav1alpha1.MicroFrontend) bool {
	statusUpdated := false
	resolvedURL := mf.Spec.Service.ResolveServiceURL(mf.Namespace)

	if resolvedURL != mf.Status.ResolvedServiceURL {
		mf.Status.ResolvedServiceURL = resolvedURL
		statusUpdated = true
	}

	if resolvedURL != "" {
		polyfeav1alpha1.SetCondition(&mf.Status.Conditions, polyfeav1alpha1.ConditionTypeServiceResolved,
			metav1.ConditionTrue, polyfeav1alpha1.ReasonSuccessful, "Service URL resolved successfully")
	} else {
		polyfeav1alpha1.SetCondition(&mf.Status.Conditions, polyfeav1alpha1.ConditionTypeServiceResolved,
			metav1.ConditionFalse, polyfeav1alpha1.ReasonInvalidConfiguration, "Unable to resolve service URL")
	}

	return statusUpdated
}

// processFrontendClass fetches MicroFrontendClass and validates namespace policy.
func (r *MicroFrontendReconciler) processFrontendClass(ctx context.Context, mf *polyfeav1alpha1.MicroFrontend) bool {
	logger := log.FromContext(ctx)
	statusUpdated := false

	frontendClassName := DefaultFrontendClassName
	if mf.Spec.FrontendClass != nil && *mf.Spec.FrontendClass != "" {
		frontendClassName = *mf.Spec.FrontendClass
	}

	mfc := &polyfeav1alpha1.MicroFrontendClass{}
	err := r.Get(ctx, client.ObjectKey{Name: frontendClassName, Namespace: mf.Namespace}, mfc)

	if err != nil {
		// If not found in the MicroFrontend's namespace, search across all namespaces
		if apierrors.IsNotFound(err) {
			mfcList := &polyfeav1alpha1.MicroFrontendClassList{}
			if listErr := r.List(ctx, mfcList, client.InNamespace("")); listErr != nil {
				logger.Error(listErr, "Failed to list MicroFrontendClasses across namespaces")
				return r.handleFrontendClassNotFound(mf, frontendClassName, listErr, logger)
			}
			found := false
			for i := range mfcList.Items {
				if mfcList.Items[i].Name == frontendClassName {
					*mfc = mfcList.Items[i]
					found = true
					break
				}
			}
			if !found {
				return r.handleFrontendClassNotFound(mf, frontendClassName, err, logger)
			}
		} else {
			return r.handleFrontendClassNotFound(mf, frontendClassName, err, logger)
		}
	}

	return r.validateNamespacePolicy(mf, mfc, frontendClassName, logger) || statusUpdated
}

// handleFrontendClassNotFound handles the case when MicroFrontendClass is not found.
func (r *MicroFrontendReconciler) handleFrontendClassNotFound(mf *polyfeav1alpha1.MicroFrontend, frontendClassName string, err error, logger logr.Logger) bool {
	statusUpdated := false

	if !apierrors.IsNotFound(err) {
		logger.Error(err, "Failed to get MicroFrontendClass", "name", frontendClassName)
		polyfeav1alpha1.SetCondition(&mf.Status.Conditions, polyfeav1alpha1.ConditionTypeFrontendClassBound,
			metav1.ConditionFalse, polyfeav1alpha1.ReasonError, "Error retrieving MicroFrontendClass")
		mf.Status.Phase = polyfeav1alpha1.MicroFrontendPhaseFailed
		statusUpdated = true
	} else {
		polyfeav1alpha1.SetCondition(&mf.Status.Conditions, polyfeav1alpha1.ConditionTypeFrontendClassBound,
			metav1.ConditionFalse, polyfeav1alpha1.ReasonFrontendClassNotFound,
			"MicroFrontendClass not found in namespace")
		mf.Status.Phase = polyfeav1alpha1.MicroFrontendPhaseFailed
		statusUpdated = true
	}

	return statusUpdated
}

// validateNamespacePolicy validates namespace policy and updates status accordingly.
func (r *MicroFrontendReconciler) validateNamespacePolicy(mf *polyfeav1alpha1.MicroFrontend, mfc *polyfeav1alpha1.MicroFrontendClass, frontendClassName string, logger logr.Logger) bool {
	statusUpdated := false
	accepted := mfc.IsNamespaceAllowed(mf.Namespace)

	// Update FrontendClassRef
	if r.shouldUpdateFrontendClassRef(mf, frontendClassName, mfc.Namespace, accepted) {
		mf.Status.FrontendClassRef = &polyfeav1alpha1.MicroFrontendClassReference{
			Name:      frontendClassName,
			Namespace: mfc.Namespace,
			Accepted:  accepted,
		}
		statusUpdated = true
	}

	polyfeav1alpha1.SetCondition(&mf.Status.Conditions, polyfeav1alpha1.ConditionTypeFrontendClassBound,
		metav1.ConditionTrue, polyfeav1alpha1.ReasonSuccessful, "Bound to MicroFrontendClass")

	if accepted {
		statusUpdated = r.handleAcceptedMicroFrontend(mf, logger) || statusUpdated
	} else {
		statusUpdated = r.handleRejectedMicroFrontend(mf, frontendClassName, logger) || statusUpdated
	}

	return statusUpdated
}

// shouldUpdateFrontendClassRef checks if FrontendClassRef needs to be updated.
func (r *MicroFrontendReconciler) shouldUpdateFrontendClassRef(mf *polyfeav1alpha1.MicroFrontend, name, namespace string, accepted bool) bool {
	return mf.Status.FrontendClassRef == nil ||
		mf.Status.FrontendClassRef.Name != name ||
		mf.Status.FrontendClassRef.Namespace != namespace ||
		mf.Status.FrontendClassRef.Accepted != accepted
}

// handleAcceptedMicroFrontend updates status for accepted MicroFrontend.
func (r *MicroFrontendReconciler) handleAcceptedMicroFrontend(mf *polyfeav1alpha1.MicroFrontend, logger logr.Logger) bool {
	statusUpdated := false

	polyfeav1alpha1.SetCondition(&mf.Status.Conditions, polyfeav1alpha1.ConditionTypeNamespacePolicyValid,
		metav1.ConditionTrue, polyfeav1alpha1.ReasonSuccessful, "Namespace is allowed by policy")
	polyfeav1alpha1.SetCondition(&mf.Status.Conditions, polyfeav1alpha1.ConditionTypeAccepted,
		metav1.ConditionTrue, polyfeav1alpha1.ReasonSuccessful, "MicroFrontend accepted")

	// Clear rejection reason if previously set
	if mf.Status.RejectionReason != "" {
		mf.Status.RejectionReason = ""
		statusUpdated = true
	}

	// Check if there are import map conflicts
	hasImportMapConflicts := len(mf.Status.ImportMapConflicts) > 0

	// Determine overall phase and readiness
	serviceResolved := polyfeav1alpha1.IsConditionTrue(mf.Status.Conditions, polyfeav1alpha1.ConditionTypeServiceResolved)

	if hasImportMapConflicts {
		// MicroFrontend has import map conflicts - not ready
		if mf.Status.Phase != polyfeav1alpha1.MicroFrontendPhaseFailed {
			mf.Status.Phase = polyfeav1alpha1.MicroFrontendPhaseFailed
			statusUpdated = true
		}
		polyfeav1alpha1.SetCondition(&mf.Status.Conditions, polyfeav1alpha1.ConditionTypeReady,
			metav1.ConditionFalse, polyfeav1alpha1.ReasonError, "Import map conflicts detected")

		// Remove from repository so it won't be served
		if err := r.Repository.Delete(mf); err != nil {
			logger.Error(err, "Failed to remove MicroFrontend with conflicts from repository")
		} else {
			logger.Info("Removed MicroFrontend with import map conflicts from repository",
				"microfrontend", mf.Name,
				"namespace", mf.Namespace,
				"conflictCount", len(mf.Status.ImportMapConflicts))
		}
	} else if serviceResolved {
		// No conflicts and service resolved - ready to serve
		if mf.Status.Phase != polyfeav1alpha1.MicroFrontendPhaseReady {
			mf.Status.Phase = polyfeav1alpha1.MicroFrontendPhaseReady
			statusUpdated = true
		}

		// Store the MicroFrontend in the repository only if no conflicts
		if err := r.Repository.Store(mf); err != nil {
			logger.Error(err, "Failed to store MicroFrontend in repository")
			polyfeav1alpha1.SetCondition(&mf.Status.Conditions, polyfeav1alpha1.ConditionTypeReady,
				metav1.ConditionFalse, polyfeav1alpha1.ReasonError, "Failed to store in repository")
			mf.Status.Phase = polyfeav1alpha1.MicroFrontendPhaseFailed
			statusUpdated = true
		} else {
			polyfeav1alpha1.SetCondition(&mf.Status.Conditions, polyfeav1alpha1.ConditionTypeReady,
				metav1.ConditionTrue, polyfeav1alpha1.ReasonSuccessful, "MicroFrontend is ready")
		}
	} else {
		// Service not resolved yet - pending
		if mf.Status.Phase != polyfeav1alpha1.MicroFrontendPhasePending {
			mf.Status.Phase = polyfeav1alpha1.MicroFrontendPhasePending
			statusUpdated = true
		}
		polyfeav1alpha1.SetCondition(&mf.Status.Conditions, polyfeav1alpha1.ConditionTypeReady,
			metav1.ConditionFalse, polyfeav1alpha1.ReasonReconciling, "Service not yet resolved")
	}

	return statusUpdated
}

// handleRejectedMicroFrontend updates status for rejected MicroFrontend.
func (r *MicroFrontendReconciler) handleRejectedMicroFrontend(mf *polyfeav1alpha1.MicroFrontend, frontendClassName string, logger logr.Logger) bool {
	statusUpdated := false
	rejectionMsg := "Namespace not allowed by MicroFrontendClass namespace policy"

	polyfeav1alpha1.SetCondition(&mf.Status.Conditions, polyfeav1alpha1.ConditionTypeNamespacePolicyValid,
		metav1.ConditionFalse, polyfeav1alpha1.ReasonNamespaceNotAllowed, rejectionMsg)
	polyfeav1alpha1.SetCondition(&mf.Status.Conditions, polyfeav1alpha1.ConditionTypeAccepted,
		metav1.ConditionFalse, polyfeav1alpha1.ReasonNamespaceNotAllowed, rejectionMsg)
	polyfeav1alpha1.SetCondition(&mf.Status.Conditions, polyfeav1alpha1.ConditionTypeReady,
		metav1.ConditionFalse, polyfeav1alpha1.ReasonNamespaceNotAllowed, rejectionMsg)

	if mf.Status.RejectionReason != rejectionMsg {
		mf.Status.RejectionReason = rejectionMsg
		statusUpdated = true
	}
	if mf.Status.Phase != polyfeav1alpha1.MicroFrontendPhaseRejected {
		mf.Status.Phase = polyfeav1alpha1.MicroFrontendPhaseRejected
		statusUpdated = true
	}

	// Remove rejected MicroFrontend from repository so it won't be served
	if err := r.Repository.Delete(mf); err != nil {
		logger.Error(err, "Failed to delete rejected MicroFrontend from repository")
	} else {
		logger.Info("Removed rejected MicroFrontend from repository")
	}

	logger.Info("MicroFrontend rejected by namespace policy",
		"microfrontend", mf.Name,
		"namespace", mf.Namespace,
		"frontendClass", frontendClassName)

	return statusUpdated
}

// checkImportMapConflicts detects conflicts in import map entries with other MicroFrontends
// in the same FrontendClass. First-registered wins based on creation timestamp.
func (r *MicroFrontendReconciler) checkImportMapConflicts(ctx context.Context, mf *polyfeav1alpha1.MicroFrontend) bool {
	logger := log.FromContext(ctx)

	// If no import map, clear any existing conflicts and return
	if mf.Spec.ImportMap == nil || len(mf.Spec.ImportMap.Imports) == 0 {
		return r.clearImportMapConflicts(mf)
	}

	// Get frontend class name
	frontendClassName := r.getFrontendClassName(mf)

	// List all MicroFrontends
	mfList := &polyfeav1alpha1.MicroFrontendList{}
	if err := r.List(ctx, mfList, client.InNamespace("")); err != nil {
		logger.Error(err, "Failed to list MicroFrontends for import map conflict check")
		return false
	}

	// Build map of existing import map entries from other accepted MicroFrontends
	existingImports, existingScopes := r.buildExistingImportMaps(mf, frontendClassName, mfList)

	// Check for conflicts in this MicroFrontend's imports
	conflicts := r.findImportMapConflicts(mf, existingImports, existingScopes)

	// Update status if conflicts changed
	return r.updateImportMapConflictStatus(mf, conflicts, logger)
}

// clearImportMapConflicts clears import map conflicts from status if present
func (r *MicroFrontendReconciler) clearImportMapConflicts(mf *polyfeav1alpha1.MicroFrontend) bool {
	if len(mf.Status.ImportMapConflicts) > 0 {
		mf.Status.ImportMapConflicts = nil
		return true
	}
	return false
}

// getFrontendClassName returns the frontend class name for the MicroFrontend
func (r *MicroFrontendReconciler) getFrontendClassName(mf *polyfeav1alpha1.MicroFrontend) string {
	if mf.Spec.FrontendClass != nil && *mf.Spec.FrontendClass != "" {
		return *mf.Spec.FrontendClass
	}
	return DefaultFrontendClassName
}

// buildExistingImportMaps builds maps of existing import entries from other MicroFrontends
func (r *MicroFrontendReconciler) buildExistingImportMaps(
	mf *polyfeav1alpha1.MicroFrontend,
	frontendClassName string,
	mfList *polyfeav1alpha1.MicroFrontendList,
) (map[string]*importMapEntry, map[string]map[string]*importMapEntry) {
	existingImports := make(map[string]*importMapEntry)
	existingScopes := make(map[string]map[string]*importMapEntry)

	for i := range mfList.Items {
		other := &mfList.Items[i]

		if !r.shouldProcessMicroFrontend(mf, other, frontendClassName) {
			continue
		}

		r.processImportMapEntries(other, existingImports, existingScopes)
	}

	return existingImports, existingScopes
}

// shouldProcessMicroFrontend checks if a MicroFrontend should be processed for import map conflicts
func (r *MicroFrontendReconciler) shouldProcessMicroFrontend(
	mf *polyfeav1alpha1.MicroFrontend,
	other *polyfeav1alpha1.MicroFrontend,
	frontendClassName string,
) bool {
	// Skip self
	if other.Namespace == mf.Namespace && other.Name == mf.Name {
		return false
	}

	// Skip if not same frontend class
	otherClassName := r.getFrontendClassName(other)
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

// processImportMapEntries processes import map entries from a MicroFrontend
func (r *MicroFrontendReconciler) processImportMapEntries(
	mf *polyfeav1alpha1.MicroFrontend,
	existingImports map[string]*importMapEntry,
	existingScopes map[string]map[string]*importMapEntry,
) {
	// Process top-level imports
	for specifier, path := range mf.Spec.ImportMap.Imports {
		r.registerImportMapEntry(specifier, path, mf, existingImports)
	}

	// Process scoped imports
	for scope, scopeImports := range mf.Spec.ImportMap.Scopes {
		if existingScopes[scope] == nil {
			existingScopes[scope] = make(map[string]*importMapEntry)
		}
		for specifier, path := range scopeImports {
			r.registerImportMapEntry(specifier, path, mf, existingScopes[scope])
		}
	}
}

// registerImportMapEntry registers an import map entry, keeping the oldest one
func (r *MicroFrontendReconciler) registerImportMapEntry(
	specifier string,
	path string,
	mf *polyfeav1alpha1.MicroFrontend,
	entries map[string]*importMapEntry,
) {
	entry := &importMapEntry{
		path:              path,
		namespace:         mf.Namespace,
		name:              mf.Name,
		creationTimestamp: mf.CreationTimestamp,
	}

	if existing, ok := entries[specifier]; ok {
		// Keep the older one (first-registered)
		if mf.CreationTimestamp.Before(&existing.creationTimestamp) {
			entries[specifier] = entry
		}
	} else {
		entries[specifier] = entry
	}
}

// findImportMapConflicts finds conflicts between requested and existing import map entries
func (r *MicroFrontendReconciler) findImportMapConflicts(
	mf *polyfeav1alpha1.MicroFrontend,
	existingImports map[string]*importMapEntry,
	existingScopes map[string]map[string]*importMapEntry,
) []polyfeav1alpha1.ImportMapConflict {
	var conflicts []polyfeav1alpha1.ImportMapConflict

	// Check top-level imports
	conflicts = r.checkImportConflicts(mf, mf.Spec.ImportMap.Imports, existingImports, "", conflicts)

	// Check scoped imports
	for scope, scopeImports := range mf.Spec.ImportMap.Scopes {
		if existingScopeMap, ok := existingScopes[scope]; ok {
			conflicts = r.checkImportConflicts(mf, scopeImports, existingScopeMap, scope, conflicts)
		}
	}

	return conflicts
}

// checkImportConflicts checks for conflicts in a set of imports
func (r *MicroFrontendReconciler) checkImportConflicts(
	mf *polyfeav1alpha1.MicroFrontend,
	imports map[string]string,
	existingEntries map[string]*importMapEntry,
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

// updateImportMapConflictStatus updates the status if conflicts changed
func (r *MicroFrontendReconciler) updateImportMapConflictStatus(
	mf *polyfeav1alpha1.MicroFrontend,
	conflicts []polyfeav1alpha1.ImportMapConflict,
	logger logr.Logger,
) bool {
	if !importMapConflictsEqual(mf.Status.ImportMapConflicts, conflicts) {
		mf.Status.ImportMapConflicts = conflicts

		if len(conflicts) > 0 {
			logger.Info("Import map conflicts detected",
				"microfrontend", mf.Name,
				"namespace", mf.Namespace,
				"conflictCount", len(conflicts))
		}
		return true
	}
	return false
}

// importMapEntry tracks an import map entry's registration details
type importMapEntry struct {
	path              string
	namespace         string
	name              string
	creationTimestamp metav1.Time
}

// importMapConflictsEqual compares two slices of ImportMapConflict for equality
func importMapConflictsEqual(a, b []polyfeav1alpha1.ImportMapConflict) bool {
	if len(a) != len(b) {
		return false
	}

	// Create maps for comparison
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

// finalizeMicroFrontend performs cleanup before deletion.
func (r *MicroFrontendReconciler) finalizeMicroFrontend(mf *polyfeav1alpha1.MicroFrontend) error {
	logger := log.FromContext(context.Background())
	if err := r.Repository.Delete(mf); err != nil {
		logger.Error(err, "Failed to delete MicroFrontend from repository")
		return err
	}
	logger.Info("Finalizer cleanup complete", "MicroFrontend", mf)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MicroFrontendReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&polyfeav1alpha1.MicroFrontend{}).
		Named("microfrontend").
		Watches(
			&polyfeav1alpha1.MicroFrontendClass{},
			handler.EnqueueRequestsFromMapFunc(r.findMicroFrontendsForClass),
		).
		Complete(r)
}

// findMicroFrontendsForClass returns reconcile requests for all MicroFrontends
// that reference the given MicroFrontendClass.
func (r *MicroFrontendReconciler) findMicroFrontendsForClass(ctx context.Context, obj client.Object) []reconcile.Request {
	mfc, ok := obj.(*polyfeav1alpha1.MicroFrontendClass)
	if !ok {
		return nil
	}

	logger := log.FromContext(ctx)
	logger.Info("MicroFrontendClass changed, finding dependent MicroFrontends", "class", mfc.Name)

	// List all MicroFrontends across all namespaces
	mfList := &polyfeav1alpha1.MicroFrontendList{}
	if err := r.List(ctx, mfList, client.InNamespace("")); err != nil {
		logger.Error(err, "Failed to list MicroFrontends for class change")
		return nil
	}

	var requests []reconcile.Request
	for _, mf := range mfList.Items {
		// Get the frontend class name for this MicroFrontend
		frontendClassName := DefaultFrontendClassName
		if mf.Spec.FrontendClass != nil && *mf.Spec.FrontendClass != "" {
			frontendClassName = *mf.Spec.FrontendClass
		}

		// If this MicroFrontend references the changed class, enqueue it for reconciliation
		if frontendClassName == mfc.Name {
			requests = append(requests, reconcile.Request{
				NamespacedName: client.ObjectKey{
					Name:      mf.Name,
					Namespace: mf.Namespace,
				},
			})
			logger.Info("Enqueueing MicroFrontend for reconciliation due to class change",
				"microfrontend", mf.Name,
				"namespace", mf.Namespace,
				"class", mfc.Name)
		}
	}

	logger.Info("Enqueued MicroFrontends for reconciliation", "count", len(requests), "class", mfc.Name)
	return requests
}
