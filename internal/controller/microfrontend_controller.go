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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
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

	return ctrl.Result{}, nil
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

	frontendClassName := "polyfea-controller-default"
	if mf.Spec.FrontendClass != nil && *mf.Spec.FrontendClass != "" {
		frontendClassName = *mf.Spec.FrontendClass
	}

	mfc := &polyfeav1alpha1.MicroFrontendClass{}
	err := r.Get(ctx, client.ObjectKey{Name: frontendClassName, Namespace: mf.Namespace}, mfc)

	if err != nil {
		return r.handleFrontendClassNotFound(mf, frontendClassName, err, logger)
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

	// Determine overall phase
	if polyfeav1alpha1.IsConditionTrue(mf.Status.Conditions, polyfeav1alpha1.ConditionTypeServiceResolved) {
		if mf.Status.Phase != polyfeav1alpha1.MicroFrontendPhaseReady {
			mf.Status.Phase = polyfeav1alpha1.MicroFrontendPhaseReady
			statusUpdated = true
		}
	} else {
		if mf.Status.Phase != polyfeav1alpha1.MicroFrontendPhasePending {
			mf.Status.Phase = polyfeav1alpha1.MicroFrontendPhasePending
			statusUpdated = true
		}
	}

	// Store the MicroFrontend in the repository only if accepted
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

	logger.Info("MicroFrontend rejected by namespace policy",
		"microfrontend", mf.Name,
		"namespace", mf.Namespace,
		"frontendClass", frontendClassName)

	return statusUpdated
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
		Complete(r)
}
