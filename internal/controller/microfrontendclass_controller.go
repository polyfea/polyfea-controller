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
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

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
// +kubebuilder:rbac:groups=polyfea.github.io,resources=microfrontends,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

const (
	OperatorServiceSelectorName  = "app"
	OperatorServiceSelectorValue = "polyfea-webserver"
)

// Reconcile moves the current state of the cluster closer to the desired state.
func (r *MicroFrontendClassReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
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

	if !controllerutil.ContainsFinalizer(mfc, FinalizerName) {
		logger.Info("Adding Finalizer for MicroFrontendClass")
		controllerutil.AddFinalizer(mfc, FinalizerName)
		if err := r.Update(ctx, mfc); err != nil {
			logger.Error(err, "Failed to update custom resource to add finalizer")
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{}, nil
	}

	if mfc.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(mfc, FinalizerName) {
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
			controllerutil.RemoveFinalizer(mfc, FinalizerName)
			if err := r.Update(ctx, mfc); err != nil {
				logger.Error(err, "Failed to remove finalizer")
				return ctrl.Result{Requeue: true}, err
			}
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, nil
	}

	logger.Info("Storing MicroFrontendClass in repository", "MicroFrontendClass", mfc.Name)
	if err := r.Repository.Store(mfc); err != nil {
		logger.Error(err, "Failed to store MicroFrontendClass in repository")
		polyfeav1alpha1.SetCondition(&mfc.Status.Conditions, polyfeav1alpha1.ConditionTypeReady,
			metav1.ConditionFalse, polyfeav1alpha1.ReasonError, "Failed to store in repository")
		mfc.Status.Phase = polyfeav1alpha1.MicroFrontendClassPhaseInvalid
		if err := r.Status().Update(ctx, mfc); err != nil {
			logger.Error(err, "Failed to update MicroFrontendClass status")
		}
		return ctrl.Result{Requeue: true}, err
	}

	statusUpdated := false
	originalStatus := mfc.Status.DeepCopy()

	mfList := &polyfeav1alpha1.MicroFrontendList{}
	if err := r.List(ctx, mfList, client.InNamespace("")); err != nil {
		logger.Error(err, "Failed to list MicroFrontends")
	} else {
		acceptedCount := int32(0)
		rejectedCount := int32(0)

		for _, mf := range mfList.Items {
			ref := mf.Status.FrontendClassRef
			if ref == nil || ref.Name != mfc.Name || ref.Namespace != mfc.Namespace {
				continue
			}
			if ref.Accepted {
				acceptedCount++
			} else {
				rejectedCount++
			}
		}

		if mfc.Status.AcceptedMicroFrontends != acceptedCount {
			mfc.Status.AcceptedMicroFrontends = acceptedCount
			statusUpdated = true
		}
		if mfc.Status.RejectedMicroFrontends != rejectedCount {
			mfc.Status.RejectedMicroFrontends = rejectedCount
			statusUpdated = true
		}
	}

	polyfeav1alpha1.SetCondition(&mfc.Status.Conditions, polyfeav1alpha1.ConditionTypeReady,
		metav1.ConditionTrue, polyfeav1alpha1.ReasonSuccessful, "MicroFrontendClass is ready")

	if mfc.Status.Phase != polyfeav1alpha1.MicroFrontendClassPhaseReady {
		mfc.Status.Phase = polyfeav1alpha1.MicroFrontendClassPhaseReady
		statusUpdated = true
	}

	if mfc.Status.ObservedGeneration != mfc.Generation {
		mfc.Status.ObservedGeneration = mfc.Generation
		statusUpdated = true
	}

	if statusUpdated {
		if err := r.Status().Update(ctx, mfc); err != nil {
			logger.Error(err, "Failed to update MicroFrontendClass status")
			mfc.Status = *originalStatus
			return ctrl.Result{Requeue: true}, err
		}
		logger.Info("Updated MicroFrontendClass status",
			"phase", mfc.Status.Phase,
			"accepted", mfc.Status.AcceptedMicroFrontends,
			"rejected", mfc.Status.RejectedMicroFrontends)
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
		Watches(
			&polyfeav1alpha1.MicroFrontend{},
			handler.EnqueueRequestsFromMapFunc(r.findMicroFrontendClassForMicroFrontend),
		).
		Complete(r)
}

// findMicroFrontendClassForMicroFrontend maps a MicroFrontend change to a reconcile
// request for the MicroFrontendClass it is bound to (via Status.FrontendClassRef).
// This ensures the MFC's accepted/rejected counts stay up-to-date whenever an MF
// is reconciled and its status changes.
func (r *MicroFrontendClassReconciler) findMicroFrontendClassForMicroFrontend(ctx context.Context, obj client.Object) []reconcile.Request {
	mf, ok := obj.(*polyfeav1alpha1.MicroFrontend)
	if !ok {
		return nil
	}
	ref := mf.Status.FrontendClassRef
	if ref == nil {
		return nil
	}
	return []reconcile.Request{
		{NamespacedName: types.NamespacedName{Name: ref.Name, Namespace: ref.Namespace}},
	}
}
