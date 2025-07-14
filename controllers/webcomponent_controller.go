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
	Repository repository.Repository[*polyfeav1alpha1.WebComponent]
}

//+kubebuilder:rbac:groups=polyfea.github.io,resources=webcomponents,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=polyfea.github.io,resources=webcomponents/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=polyfea.github.io,resources=webcomponents/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

func (r *WebComponentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	const webComponentFinalizer = "polyfea.github.io/finalizer"
	log := log.FromContext(ctx)

	webComponent := &polyfeav1alpha1.WebComponent{}
	if err := r.Get(ctx, req.NamespacedName, webComponent); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("WebComponent resource not found. Ignoring since object must be deleted!")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get WebComponent!")
		return ctrl.Result{Requeue: true}, err
	}

	log.Info("Reconciling WebComponent.", "WebComponent", webComponent)

	if !controllerutil.ContainsFinalizer(webComponent, webComponentFinalizer) {
		log.Info("Adding Finalizer for WebComponent.")
		controllerutil.AddFinalizer(webComponent, webComponentFinalizer)
		if err := r.Update(ctx, webComponent); err != nil {
			log.Error(err, "Failed to update custom resource to add finalizer!")
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{}, nil
	}

	if webComponent.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(webComponent, webComponentFinalizer) {
			log.Info("Performing finalizer operations for the WebComponent before deleting the custom resource.")
			if err := r.finalizeOperationsForWebComponent(webComponent); err != nil {
				log.Error(err, "Failed to perform finalizer operations for the WebComponent!")
				return ctrl.Result{Requeue: true}, nil
			}
			controllerutil.RemoveFinalizer(webComponent, webComponentFinalizer)
			if err := r.Update(ctx, webComponent); err != nil {
				log.Error(err, "Failed to remove finalizer for WebComponent!")
				return ctrl.Result{Requeue: true}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if err := r.Repository.Store(webComponent); err != nil {
		log.Error(err, "Failed to store WebComponent in repository!")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *WebComponentReconciler) finalizeOperationsForWebComponent(webComponent *polyfeav1alpha1.WebComponent) error {
	log := log.FromContext(context.Background())
	if err := r.Repository.Delete(webComponent); err != nil {
		log.Error(err, "Failed to delete WebComponent from repository!")
		return err
	}
	log.Info("Finalizer cleanup complete for WebComponent.", "WebComponent", webComponent)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebComponentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&polyfeav1alpha1.WebComponent{}).
		Complete(r)
}
