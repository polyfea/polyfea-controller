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

	polyfeav1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/internal/repository"
)

// WebComponentReconciler reconciles a WebComponent object
type WebComponentReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Repository repository.Repository[*polyfeav1alpha1.WebComponent]
}

// +kubebuilder:rbac:groups=polyfea.github.io,resources=webcomponents,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=polyfea.github.io,resources=webcomponents/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=polyfea.github.io,resources=webcomponents/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

func (r *WebComponentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	const webComponentFinalizer = "polyfea.github.io/finalizer"
	logger := log.FromContext(ctx)

	webComponent := &polyfeav1alpha1.WebComponent{}
	if err := r.Get(ctx, req.NamespacedName, webComponent); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("WebComponent resource not found. Ignoring since object must be deleted!")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get WebComponent!")
		return ctrl.Result{Requeue: true}, err
	}

	logger.Info("Reconciling WebComponent.", "WebComponent", webComponent)

	if !controllerutil.ContainsFinalizer(webComponent, webComponentFinalizer) {
		logger.Info("Adding Finalizer for WebComponent.")
		controllerutil.AddFinalizer(webComponent, webComponentFinalizer)
		if err := r.Update(ctx, webComponent); err != nil {
			logger.Error(err, "Failed to update custom resource to add finalizer!")
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{}, nil
	}

	if webComponent.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(webComponent, webComponentFinalizer) {
			logger.Info("Performing finalizer operations for the WebComponent before deleting the custom resource.")
			if err := r.finalizeOperationsForWebComponent(webComponent); err != nil {
				logger.Error(err, "Failed to perform finalizer operations for the WebComponent!")
				return ctrl.Result{Requeue: true}, nil
			}
			controllerutil.RemoveFinalizer(webComponent, webComponentFinalizer)
			if err := r.Update(ctx, webComponent); err != nil {
				logger.Error(err, "Failed to remove finalizer for WebComponent!")
				return ctrl.Result{Requeue: true}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Update status
	statusUpdated := false
	originalStatus := webComponent.Status.DeepCopy()

	// Check if MicroFrontend reference exists
	if webComponent.Spec.MicroFrontend != nil && *webComponent.Spec.MicroFrontend != "" {
		mfName := *webComponent.Spec.MicroFrontend
		mf := &polyfeav1alpha1.MicroFrontend{}

		// Look for MicroFrontend in the same namespace
		mfFound := false
		if err := r.Get(ctx, client.ObjectKey{Name: mfName, Namespace: webComponent.Namespace}, mf); err == nil {
			mfFound = true
		} else if !apierrors.IsNotFound(err) {
			logger.Error(err, "Failed to get MicroFrontend", "name", mfName)
			polyfeav1alpha1.SetCondition(&webComponent.Status.Conditions, polyfeav1alpha1.ConditionTypeMicroFrontendResolved,
				metav1.ConditionFalse, polyfeav1alpha1.ReasonError, "Error retrieving MicroFrontend")
			webComponent.Status.Phase = polyfeav1alpha1.WebComponentPhaseFailed
			statusUpdated = true
		}

		// Update MicroFrontendRef
		if webComponent.Status.MicroFrontendRef == nil ||
			webComponent.Status.MicroFrontendRef.Name != mfName ||
			webComponent.Status.MicroFrontendRef.Namespace != webComponent.Namespace ||
			webComponent.Status.MicroFrontendRef.Found != mfFound {
			webComponent.Status.MicroFrontendRef = &polyfeav1alpha1.ObjectReference{
				Name:      mfName,
				Namespace: webComponent.Namespace,
				Found:     mfFound,
			}
			statusUpdated = true
		}

		if mfFound {
			polyfeav1alpha1.SetCondition(&webComponent.Status.Conditions, polyfeav1alpha1.ConditionTypeMicroFrontendResolved,
				metav1.ConditionTrue, polyfeav1alpha1.ReasonSuccessful, "MicroFrontend found and resolved")

			// Check if the MicroFrontend is ready
			if polyfeav1alpha1.IsReady(mf.Status.Conditions) {
				polyfeav1alpha1.SetCondition(&webComponent.Status.Conditions, polyfeav1alpha1.ConditionTypeReady,
					metav1.ConditionTrue, polyfeav1alpha1.ReasonSuccessful, "WebComponent is ready")
				if webComponent.Status.Phase != polyfeav1alpha1.WebComponentPhaseReady {
					webComponent.Status.Phase = polyfeav1alpha1.WebComponentPhaseReady
					statusUpdated = true
				}
			} else {
				polyfeav1alpha1.SetCondition(&webComponent.Status.Conditions, polyfeav1alpha1.ConditionTypeReady,
					metav1.ConditionFalse, polyfeav1alpha1.ReasonReconciling, "Waiting for MicroFrontend to be ready")
				if webComponent.Status.Phase != polyfeav1alpha1.WebComponentPhasePending {
					webComponent.Status.Phase = polyfeav1alpha1.WebComponentPhasePending
					statusUpdated = true
				}
			}

			// Store in repository only if MicroFrontend is ready
			if err := r.Repository.Store(webComponent); err != nil {
				logger.Error(err, "Failed to store WebComponent in repository!")
				polyfeav1alpha1.SetCondition(&webComponent.Status.Conditions, polyfeav1alpha1.ConditionTypeReady,
					metav1.ConditionFalse, polyfeav1alpha1.ReasonError, "Failed to store in repository")
				webComponent.Status.Phase = polyfeav1alpha1.WebComponentPhaseFailed
				statusUpdated = true
			}
		} else {
			polyfeav1alpha1.SetCondition(&webComponent.Status.Conditions, polyfeav1alpha1.ConditionTypeMicroFrontendResolved,
				metav1.ConditionFalse, polyfeav1alpha1.ReasonMicroFrontendNotFound, "MicroFrontend not found in namespace")
			polyfeav1alpha1.SetCondition(&webComponent.Status.Conditions, polyfeav1alpha1.ConditionTypeReady,
				metav1.ConditionFalse, polyfeav1alpha1.ReasonMicroFrontendNotFound, "MicroFrontend not found")
			if webComponent.Status.Phase != polyfeav1alpha1.WebComponentPhaseMicroFrontendNotFound {
				webComponent.Status.Phase = polyfeav1alpha1.WebComponentPhaseMicroFrontendNotFound
				statusUpdated = true
			}
			logger.Info("MicroFrontend not found", "webComponent", webComponent.Name, "microFrontend", mfName)

			// Still store in repository even if MicroFrontend is not found
			// This allows the WebComponent to be registered and potentially work with other mechanisms
			if err := r.Repository.Store(webComponent); err != nil {
				logger.Error(err, "Failed to store WebComponent in repository!")
			}
		}
	} else {
		// No MicroFrontend reference
		polyfeav1alpha1.SetCondition(&webComponent.Status.Conditions, polyfeav1alpha1.ConditionTypeMicroFrontendResolved,
			metav1.ConditionFalse, polyfeav1alpha1.ReasonInvalidConfiguration, "No MicroFrontend reference specified")
		polyfeav1alpha1.SetCondition(&webComponent.Status.Conditions, polyfeav1alpha1.ConditionTypeReady,
			metav1.ConditionFalse, polyfeav1alpha1.ReasonInvalidConfiguration, "No MicroFrontend reference")
		if webComponent.Status.Phase != polyfeav1alpha1.WebComponentPhaseFailed {
			webComponent.Status.Phase = polyfeav1alpha1.WebComponentPhaseFailed
			statusUpdated = true
		}

		// Still store it in repository as it might be used differently
		if err := r.Repository.Store(webComponent); err != nil {
			logger.Error(err, "Failed to store WebComponent in repository!")
		}
	}

	// Update ObservedGeneration
	if webComponent.Status.ObservedGeneration != webComponent.Generation {
		webComponent.Status.ObservedGeneration = webComponent.Generation
		statusUpdated = true
	}

	// Update status if needed
	if statusUpdated {
		if err := r.Status().Update(ctx, webComponent); err != nil {
			logger.Error(err, "Failed to update WebComponent status")
			webComponent.Status = *originalStatus
			return ctrl.Result{Requeue: true}, err
		}
		logger.Info("Updated WebComponent status", "phase", webComponent.Status.Phase)
	}

	return ctrl.Result{}, nil
}

func (r *WebComponentReconciler) finalizeOperationsForWebComponent(webComponent *polyfeav1alpha1.WebComponent) error {
	logger := log.FromContext(context.Background())
	if err := r.Repository.Delete(webComponent); err != nil {
		logger.Error(err, "Failed to delete WebComponent from repository!")
		return err
	}
	logger.Info("Finalizer cleanup complete for WebComponent.", "WebComponent", webComponent)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebComponentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&polyfeav1alpha1.WebComponent{}).
		Complete(r)
}
