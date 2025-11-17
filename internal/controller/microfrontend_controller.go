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
	log := log.FromContext(ctx)

	mf := &polyfeav1alpha1.MicroFrontend{}
	if err := r.Get(ctx, req.NamespacedName, mf); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("MicroFrontend resource not found; assuming it was deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get MicroFrontend")
		return ctrl.Result{Requeue: true}, err
	}

	log.Info("Reconciling MicroFrontend", "MicroFrontend", mf)

	// Add finalizer if not present
	if !controllerutil.ContainsFinalizer(mf, finalizerName) {
		log.Info("Adding finalizer")
		controllerutil.AddFinalizer(mf, finalizerName)
		if err := r.Update(ctx, mf); err != nil {
			log.Error(err, "Failed to update MicroFrontend with finalizer")
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{}, nil
	}

	// Handle deletion
	if mf.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(mf, finalizerName) {
			log.Info("Running finalizer operations before deletion")
			if err := r.finalizeMicroFrontend(mf); err != nil {
				log.Error(err, "Finalizer operations failed")
				return ctrl.Result{Requeue: true}, nil
			}
			// Re-fetch in case of update
			if err := r.Get(ctx, req.NamespacedName, mf); err != nil {
				log.Error(err, "Failed to re-fetch MicroFrontend after finalizer")
				return ctrl.Result{Requeue: true}, err
			}
			controllerutil.RemoveFinalizer(mf, finalizerName)
			if err := r.Update(ctx, mf); err != nil {
				log.Error(err, "Failed to remove finalizer")
				return ctrl.Result{Requeue: true}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Store the MicroFrontend in the repository
	if err := r.Repository.Store(mf); err != nil {
		log.Error(err, "Failed to store MicroFrontend in repository")
		return ctrl.Result{Requeue: true}, err
	}

	return ctrl.Result{}, nil
}

// finalizeMicroFrontend performs cleanup before deletion.
func (r *MicroFrontendReconciler) finalizeMicroFrontend(mf *polyfeav1alpha1.MicroFrontend) error {
	log := log.FromContext(context.Background())
	if err := r.Repository.Delete(mf); err != nil {
		log.Error(err, "Failed to delete MicroFrontend from repository")
		return err
	}
	log.Info("Finalizer cleanup complete", "MicroFrontend", mf)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MicroFrontendReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&polyfeav1alpha1.MicroFrontend{}).
		Named("microfrontend").
		Complete(r)
}
