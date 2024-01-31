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
	"reflect"
	"slices"

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/go-logr/logr"
	polyfeav1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
)

// MicroFrontendClassReconciler reconciles a MicroFrontendClass object
type MicroFrontendClassReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	Recorder   record.EventRecorder
	Repository repository.PolyfeaRepository[*polyfeav1alpha1.MicroFrontendClass]
}

const (
	PortName                     = "webserver"
	OperatorServiceSelectorName  = "app"
	OperatorServiceSelectorValue = "polyfea-webserver"
)

//+kubebuilder:rbac:groups=polyfea.github.io,resources=microfrontendclasses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=polyfea.github.io,resources=microfrontendclasses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=polyfea.github.io,resources=microfrontendclasses/finalizers,verbs=update
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=httproutes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *MicroFrontendClassReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	const microFrontendClassFinalizer = "polyfea.github.io/finalizer"

	log := log.FromContext(ctx)

	// Fetch the microFrontendClass instance
	// The purpose is check if the Custom Resource for the Kind MicroFrontendClass
	// is applied on the cluster if not we return nil to stop the reconciliation
	microFrontendClass := &polyfeav1alpha1.MicroFrontendClass{}
	err := r.Get(ctx, req.NamespacedName, microFrontendClass)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// If the custom resource is not found then, it usually means that it was deleted or not created
			// In this way, we will stop the reconciliation
			log.Info("MicroFrontendClass resource not found. Ignoring since object must be deleted!")
			return ctrl.Result{}, nil
		}

		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get MicroFrontendClass!")
		return ctrl.Result{}, err
	}

	log.Info("Reconciling MicroFrontendClass.", "MicroFrontendClass", microFrontendClass)

	// Let's add a finalizer. Then, we can define some operations which should
	// occurs before the custom resource to be deleted.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/finalizers
	if !controllerutil.ContainsFinalizer(microFrontendClass, microFrontendClassFinalizer) {
		log.Info("Adding Finalizer for MicroFrontendClass.")
		if ok := controllerutil.AddFinalizer(microFrontendClass, microFrontendClassFinalizer); !ok {
			log.Error(err, "Failed to add finalizer into the custom resource!")
			return ctrl.Result{Requeue: true}, nil
		}

		if err = r.Update(ctx, microFrontendClass); err != nil {
			log.Error(err, "Failed to update custom resource to add finalizer!")
			return ctrl.Result{}, err
		}

		return ctrl.Result{Requeue: true}, nil
	}

	// Check if the MicroFrontendClass instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isMicroFrontendClassMarkedToBeDeleted := microFrontendClass.GetDeletionTimestamp() != nil
	if isMicroFrontendClassMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(microFrontendClass, microFrontendClassFinalizer) {
			log.Info("Performing finalizer operations for the MicroFrontendClass before deleting the custom resource.")

			if err := r.finalizeOperationsForMicroFrontendClass(microFrontendClass); err != nil {
				log.Error(err, "Failed to perform finalizer operations for the MicroFrontendClass!")
				return ctrl.Result{Requeue: true}, nil
			}

			if err := r.Get(ctx, req.NamespacedName, microFrontendClass); err != nil {
				log.Error(err, "Failed to re-fetch MicroFrontendClass!")
				return ctrl.Result{}, err
			}

			log.Info("Removing Finalizer for MicroFrontendClass after successfully performing the operations.")
			if ok := controllerutil.RemoveFinalizer(microFrontendClass, microFrontendClassFinalizer); !ok {
				log.Error(err, "Failed to remove finalizer for MicroFrontendClass!")
				return ctrl.Result{Requeue: true}, nil
			}

			if err := r.Update(ctx, microFrontendClass); err != nil {
				log.Error(err, "Failed to remove finalizer for MicroFrontendClass!")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, nil
	}

	httpRoute := &gatewayv1.HTTPRoute{}
	ingress := &networkingv1.Ingress{}

	log.Info("Checking if HttpRoute and Ingress resources exist for MicroFrontendClass.")
	err = r.Get(ctx, client.ObjectKey{Namespace: microFrontendClass.Namespace, Name: microFrontendClass.Name}, httpRoute)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("HttpRoute resource not found. Setting to nil")
			httpRoute = nil
		} else {
			log.Error(err, "Failed to get HttpRoute!")
			return ctrl.Result{}, err
		}
	}

	err = r.Get(ctx, client.ObjectKey{Namespace: microFrontendClass.Namespace, Name: microFrontendClass.Name}, ingress)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Ingress resource not found. Setting to nil")
			ingress = nil
		} else {
			log.Error(err, "Failed to get Ingress!")
			return ctrl.Result{}, err
		}
	}

	log.Info("Checking if HttpRoute and Ingress resources should be created for MicroFrontendClass.")
	if microFrontendClass.Spec.Routing != nil && httpRoute == nil && ingress == nil {
		operatorService := r.getOperatorService(ctx, log)

		if operatorService == nil {
			log.Error(err, "Failed to get OperatorService!")
			return ctrl.Result{}, err
		}

		if microFrontendClass.Spec.Routing.ParentRefs != nil {
			err = r.createHttpRoute(ctx, log, microFrontendClass, operatorService)

			if err != nil {
				log.Error(err, "Failed to create HttpRoute!")
				return ctrl.Result{}, err
			}
		}

		if microFrontendClass.Spec.Routing.IngressClassName != nil {
			err = r.createIngress(ctx, log, microFrontendClass, operatorService)

			if err != nil {
				log.Error(err, "Failed to create HttpRoute!")
				return ctrl.Result{}, err
			}
		}

	}

	log.Info("Checking if HttpRoute and Ingress resources should be updated for MicroFrontendClass.")
	if microFrontendClass.Spec.Routing != nil && httpRoute != nil {
		if microFrontendClass.Spec.Routing.ParentRefs == nil && microFrontendClass.Spec.Routing.IngressClassName != nil {
			err = r.Delete(ctx, httpRoute)

			if err != nil {
				log.Error(err, "Failed to delete HttpRoute!")
				return ctrl.Result{}, err
			}

			err = r.createIngress(ctx, log, microFrontendClass, r.getOperatorService(ctx, log))

			if err != nil {
				log.Error(err, "Failed to create Ingress!")
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		}

		operatorService := r.getOperatorService(ctx, log)

		if operatorService == nil {
			log.Error(err, "Failed to get OperatorService!")
			return ctrl.Result{}, err
		}
		update := createRouteForMicroFrontendClass(microFrontendClass, operatorService)
		controllerutil.SetControllerReference(microFrontendClass, update, r.Scheme)

		if httpRoute.Name != update.Name || httpRoute.Namespace != update.Namespace || !reflect.DeepEqual(httpRoute.Spec, update.Spec) {
			log.Info("Updating HttpRoute for MicroFrontendClass.")
			update.ResourceVersion = httpRoute.ResourceVersion
			err = r.Update(ctx, update)

			if err != nil {
				log.Error(err, "Failed to update HttpRoute!")
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		}
	}

	log.Info("Checking if HttpRoute and Ingress resources should be deleted for MicroFrontendClass.")
	if microFrontendClass.Spec.Routing != nil && ingress != nil {
		if microFrontendClass.Spec.Routing.IngressClassName == nil && microFrontendClass.Spec.Routing.ParentRefs != nil {
			err = r.Delete(ctx, ingress)

			if err != nil {
				log.Error(err, "Failed to delete Ingress!")
				return ctrl.Result{}, err
			}

			err = r.createHttpRoute(ctx, log, microFrontendClass, r.getOperatorService(ctx, log))

			if err != nil {
				log.Error(err, "Failed to create HttpRoute!")
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		}

		operatorService := r.getOperatorService(ctx, log)

		if operatorService == nil {
			log.Error(err, "Failed to get OperatorService!")
			return ctrl.Result{}, err
		}

		update := createIngressForMicroFrontendClass(microFrontendClass, operatorService)
		controllerutil.SetControllerReference(microFrontendClass, update, r.Scheme)
		if ingress.Name != update.Name || ingress.Namespace != update.Namespace || !reflect.DeepEqual(ingress.Spec, update.Spec) {
			log.Info("Updating Ingress for MicroFrontendClass.")
			update.ResourceVersion = ingress.ResourceVersion
			err = r.Update(ctx, update)

			if err != nil {
				log.Error(err, "Failed to update HttpRoute!")
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		}
	}

	if microFrontendClass.Spec.Routing == nil && httpRoute != nil {
		err = r.Delete(ctx, httpRoute)

		if err != nil {
			log.Error(err, "Failed to delete HttpRoute!")
			return ctrl.Result{}, err
		}
	}

	if microFrontendClass.Spec.Routing == nil && ingress != nil {
		err = r.Delete(ctx, ingress)

		if err != nil {
			log.Error(err, "Failed to delete Ingress!")
			return ctrl.Result{}, err
		}
	}

	err = r.Repository.StoreItem(microFrontendClass)

	if err != nil {
		log.Error(err, "Failed to store MicroFrontendClass into repository!")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *MicroFrontendClassReconciler) finalizeOperationsForMicroFrontendClass(microFrontendClass *polyfeav1alpha1.MicroFrontendClass) error {
	log := log.FromContext(context.Background())
	r.Repository.DeleteItem(microFrontendClass)
	log.Info("Removing finalizer from MicroFrontendClass.", "MicroFrontendClass", microFrontendClass)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MicroFrontendClassReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&polyfeav1alpha1.MicroFrontendClass{}).
		Owns(&gatewayv1.HTTPRoute{}).
		Owns(&networkingv1.Ingress{}).
		Complete(r)
}

func (r *MicroFrontendClassReconciler) createHttpRoute(ctx context.Context, log logr.Logger, microFrontendClass *polyfeav1alpha1.MicroFrontendClass, operatorService *corev1.Service) error {
	httpRoute := createRouteForMicroFrontendClass(microFrontendClass, operatorService)
	err := controllerutil.SetControllerReference(microFrontendClass, httpRoute, r.Scheme)

	if err != nil {
		log.Error(err, "Failed to set controller reference for HttpRoute!")
		return err
	}

	err = r.Create(ctx, httpRoute)

	if err != nil {
		log.Error(err, "Failed to create HttpRoute!")
		return err
	}

	return nil
}

func (r *MicroFrontendClassReconciler) createIngress(ctx context.Context, log logr.Logger, microFrontendClass *polyfeav1alpha1.MicroFrontendClass, operatorService *corev1.Service) error {
	ingress := createIngressForMicroFrontendClass(microFrontendClass, operatorService)
	controllerutil.SetControllerReference(microFrontendClass, ingress, r.Scheme)
	err := r.Create(ctx, ingress)

	if err != nil {
		log.Error(err, "Failed to create Ingress!")
		return err
	}

	return nil
}

func (r *MicroFrontendClassReconciler) getOperatorService(ctx context.Context, log logr.Logger) *corev1.Service {
	listOptions := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set{OperatorServiceSelectorName: OperatorServiceSelectorValue}),
	}

	var services corev1.ServiceList

	if err := r.List(ctx, &services, listOptions); err != nil || len(services.Items) == 0 {
		log.Error(err, "Failed to list services or no services found!")
		return nil
	}

	return &services.Items[0]
}

func createRouteForMicroFrontendClass(microFrontendClass *polyfeav1alpha1.MicroFrontendClass, operatorService *corev1.Service) *gatewayv1.HTTPRoute {
	kind := gatewayv1.Kind("Service")
	name := gatewayv1.ObjectName(operatorService.Name)
	portIndex := slices.IndexFunc(operatorService.Spec.Ports, func(port corev1.ServicePort) bool {
		return port.Name == PortName
	})
	port := gatewayv1.PortNumber(operatorService.Spec.Ports[portIndex].Port)

	return &gatewayv1.HTTPRoute{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      microFrontendClass.Name,
			Namespace: microFrontendClass.Namespace,
		},
		Spec: gatewayv1.HTTPRouteSpec{
			CommonRouteSpec: gatewayv1.CommonRouteSpec{
				ParentRefs: microFrontendClass.Spec.Routing.ParentRefs,
			},
			Rules: []gatewayv1.HTTPRouteRule{
				{
					Matches: []gatewayv1.HTTPRouteMatch{
						{
							Path: &gatewayv1.HTTPPathMatch{
								Value: microFrontendClass.Spec.BaseUri,
							},
						},
					},
					BackendRefs: []gatewayv1.HTTPBackendRef{
						{
							BackendRef: gatewayv1.BackendRef{
								BackendObjectReference: gatewayv1.BackendObjectReference{
									Kind: &kind,
									Name: name,
									Port: &port,
								},
							},
						},
					},
				},
			},
		},
	}
}

func createIngressForMicroFrontendClass(microFrontendClass *polyfeav1alpha1.MicroFrontendClass, operatorService *corev1.Service) *networkingv1.Ingress {
	pathType := networkingv1.PathTypePrefix

	return &networkingv1.Ingress{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      microFrontendClass.Name,
			Namespace: microFrontendClass.Namespace,
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: microFrontendClass.Spec.Routing.IngressClassName,
			Rules: []networkingv1.IngressRule{
				{
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     *microFrontendClass.Spec.BaseUri,
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: operatorService.Name,
											Port: networkingv1.ServiceBackendPort{
												Name: PortName,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
