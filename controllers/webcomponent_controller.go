/*
Copyright 2023.

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

package controllers

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	polyfeav1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
)

// WebComponentReconciler reconciles a WebComponent object
type WebComponentReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Recorder   record.EventRecorder
	Repository repository.PolyfeaRepository[*polyfeav1alpha1.WebComponent]
}

//+kubebuilder:rbac:groups=polyfea.github.io,resources=webcomponents,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=polyfea.github.io,resources=webcomponents/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=polyfea.github.io,resources=webcomponents/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *WebComponentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	const webComponentFinalizer = "polyfea.github.io/finalizer"

	log := log.FromContext(ctx)

	// Fetch the webComponent instance
	// The purpose is check if the Custom Resource for the Kind WebComponent
	// is applied on the cluster if not we return nil to stop the reconciliation
	webComponent := &polyfeav1alpha1.WebComponent{}
	err := r.Get(ctx, req.NamespacedName, webComponent)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// If the custom resource is not found then, it usually means that it was deleted or not created
			// In this way, we will stop the reconciliation
			log.Info("WebComponent resource not found. Ignoring since object must be deleted!")
			return ctrl.Result{}, nil
		}

		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get WebComponent!")
		return ctrl.Result{}, err
	}

	log.Info("Reconciling WebComponent.", "WebComponent", webComponent)

	// Let's add a finalizer. Then, we can define some operations which should
	// occurs before the custom resource to be deleted.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/finalizers
	if !controllerutil.ContainsFinalizer(webComponent, webComponentFinalizer) {
		log.Info("Adding Finalizer for WebComponent.")
		if ok := controllerutil.AddFinalizer(webComponent, webComponentFinalizer); !ok {
			log.Error(err, "Failed to add finalizer into the custom resource!")
			return ctrl.Result{Requeue: true}, nil
		}

		if err = r.Update(ctx, webComponent); err != nil {
			log.Error(err, "Failed to update custom resource to add finalizer!")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	// Check if the WebComponent instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isWebComponentMarkedToBeDeleted := webComponent.GetDeletionTimestamp() != nil
	if isWebComponentMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(webComponent, webComponentFinalizer) {
			log.Info("Performing finalizer operations for the WebComponent before deleting the custom resource.")

			if err := r.finalizeOperationsForWebComponent(webComponent); err != nil {
				log.Error(err, "Failed to perform finalizer operations for the WebComponent!")
				return ctrl.Result{Requeue: true}, nil
			}

			if err := r.Get(ctx, req.NamespacedName, webComponent); err != nil {
				log.Error(err, "Failed to re-fetch WebComponent!")
				return ctrl.Result{}, err
			}

			log.Info("Removing Finalizer for WebComponent after successfully performing the operations.")
			if ok := controllerutil.RemoveFinalizer(webComponent, webComponentFinalizer); !ok {
				log.Error(err, "Failed to remove finalizer for WebComponent!")
				return ctrl.Result{Requeue: true}, nil
			}

			if err := r.Update(ctx, webComponent); err != nil {
				log.Error(err, "Failed to remove finalizer for WebComponent!")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, nil
	}

	r.Repository.StoreItem(webComponent)

	return ctrl.Result{}, nil
}

func (r *WebComponentReconciler) finalizeOperationsForWebComponent(webComponent *polyfeav1alpha1.WebComponent) error {
	log := log.FromContext(context.Background())
	r.Repository.DeleteItem(webComponent)
	log.Info("Removing finalizer from WebComponent.", "WebComponent", webComponent)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebComponentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&polyfeav1alpha1.WebComponent{}).
		Complete(r)
}
