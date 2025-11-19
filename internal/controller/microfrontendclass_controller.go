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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	polyfeav1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/internal/repository"
)

// MicroFrontendClassReconciler reconciles a MicroFrontendClass object.
type MicroFrontendClassReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Repository repository.Repository[*polyfeav1alpha1.MicroFrontendClass]
}

// +kubebuilder:rbac:groups=polyfea.github.io,resources=microfrontendclasses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=polyfea.github.io,resources=microfrontendclasses/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=polyfea.github.io,resources=microfrontendclasses/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch

const (
	OperatorServiceSelectorName  = "app"
	OperatorServiceSelectorValue = "polyfea-webserver"
)

// Reconcile moves the current state of the cluster closer to the desired state.
func (r *MicroFrontendClassReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	const finalizerName = "polyfea.github.io/finalizer"
	logger := log.FromContext(ctx)

	mfc := &polyfeav1alpha1.MicroFrontendClass{}
	if err := r.Get(ctx, req.NamespacedName, mfc); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("MicroFrontendClass resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get MicroFrontendClass")
		return ctrl.Result{Requeue: true}, err
	}

	logger.Info("Reconciling MicroFrontendClass", "MicroFrontendClass", mfc)

	// Add finalizer if not present
	if !controllerutil.ContainsFinalizer(mfc, finalizerName) {
		logger.Info("Adding Finalizer for MicroFrontendClass")
		controllerutil.AddFinalizer(mfc, finalizerName)
		if err := r.Update(ctx, mfc); err != nil {
			logger.Error(err, "Failed to update custom resource to add finalizer")
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{}, nil
	}

	// Handle deletion
	if mfc.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(mfc, finalizerName) {
			logger.Info("Performing finalizer operations before deletion")
			if err := r.finalizeOperationsForMicroFrontendClass(mfc); err != nil {
				logger.Error(err, "Failed to perform finalizer operations")
				return ctrl.Result{Requeue: true}, nil
			}
			if err := r.Get(ctx, req.NamespacedName, mfc); err != nil {
				logger.Error(err, "Failed to re-fetch MicroFrontendClass")
				return ctrl.Result{Requeue: true}, err
			}
			logger.Info("Removing Finalizer for MicroFrontendClass after successful operations")
			controllerutil.RemoveFinalizer(mfc, finalizerName)
			if err := r.Update(ctx, mfc); err != nil {
				logger.Error(err, "Failed to remove finalizer")
				return ctrl.Result{Requeue: true}, err
			}
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, nil
	}

	// Store the MicroFrontendClass in the repository
	logger.Info("Storing MicroFrontendClass in repository", "MicroFrontendClass", mfc)
	if err := r.Repository.Store(mfc); err != nil {
		logger.Error(err, "Failed to store MicroFrontendClass in repository")
		return ctrl.Result{Requeue: true}, err
	}

	return ctrl.Result{}, nil
}

func (r *MicroFrontendClassReconciler) finalizeOperationsForMicroFrontendClass(mfc *polyfeav1alpha1.MicroFrontendClass) error {
	logger := log.FromContext(context.Background())
	err := r.Repository.Delete(mfc)
	if err != nil {
		logger.Error(err, "Failed to remove MicroFrontendClass from repository", "MicroFrontendClass", mfc)
		return err
	}
	logger.Info("Removing finalizer from MicroFrontendClass", "MicroFrontendClass", mfc)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MicroFrontendClassReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&polyfeav1alpha1.MicroFrontendClass{}).
		Complete(r)
}
