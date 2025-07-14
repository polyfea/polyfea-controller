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
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/go-logr/logr"
	polyfeav1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
	"github.com/polyfea/polyfea-controller/repository"
)

// MicroFrontendClassReconciler reconciles a MicroFrontendClass object.
type MicroFrontendClassReconciler struct {
	client.Client
	Scheme            *runtime.Scheme
	Recorder          record.EventRecorder
	Repository        repository.Repository[*polyfeav1alpha1.MicroFrontendClass]
	selfRef           controller.Controller
	cacheRef          cache.Cache
	isAlreadyWatching bool
}

//+kubebuilder:rbac:groups=polyfea.github.io,resources=microfrontendclasses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=polyfea.github.io,resources=microfrontendclasses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=polyfea.github.io,resources=microfrontendclasses/finalizers,verbs=update
//+kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=httproutes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch

const (
	OperatorServiceSelectorName  = "app"
	OperatorServiceSelectorValue = "polyfea-webserver"
)

// Reconcile moves the current state of the cluster closer to the desired state.
func (r *MicroFrontendClassReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	const finalizerName = "polyfea.github.io/finalizer"
	log := log.FromContext(ctx)

	mfc := &polyfeav1alpha1.MicroFrontendClass{}
	if err := r.Get(ctx, req.NamespacedName, mfc); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("MicroFrontendClass resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get MicroFrontendClass")
		return ctrl.Result{Requeue: true}, err
	}

	log.Info("Reconciling MicroFrontendClass", "MicroFrontendClass", mfc)

	// Add finalizer if not present
	if !controllerutil.ContainsFinalizer(mfc, finalizerName) {
		log.Info("Adding Finalizer for MicroFrontendClass")
		controllerutil.AddFinalizer(mfc, finalizerName)
		if err := r.Update(ctx, mfc); err != nil {
			log.Error(err, "Failed to update custom resource to add finalizer")
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{}, nil
	}

	// Handle deletion
	if mfc.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(mfc, finalizerName) {
			log.Info("Performing finalizer operations before deletion")
			if err := r.finalizeOperationsForMicroFrontendClass(mfc); err != nil {
				log.Error(err, "Failed to perform finalizer operations")
				return ctrl.Result{Requeue: true}, nil
			}
			if err := r.Get(ctx, req.NamespacedName, mfc); err != nil {
				log.Error(err, "Failed to re-fetch MicroFrontendClass")
				return ctrl.Result{Requeue: true}, err
			}
			log.Info("Removing Finalizer for MicroFrontendClass after successful operations")
			controllerutil.RemoveFinalizer(mfc, finalizerName)
			if err := r.Update(ctx, mfc); err != nil {
				log.Error(err, "Failed to remove finalizer")
				return ctrl.Result{Requeue: true}, err
			}
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, nil
	}

	// Resource existence checks
	httpRoute := &gatewayv1.HTTPRoute{}
	httpRoutePresent := true
	ingress := &networkingv1.Ingress{}

	log.Info("Checking if HttpRoute and Ingress resources exist for MicroFrontendClass")
	err := r.Get(ctx, client.ObjectKey{Namespace: mfc.Namespace, Name: mfc.Name}, httpRoute)
	if err != nil {
		if apierrors.IsNotFound(err) {
			if !r.isAlreadyWatching {
				r.selfRef.Watch(source.Kind(r.cacheRef, &gatewayv1.HTTPRoute{}), &handler.EnqueueRequestForObject{})
				r.isAlreadyWatching = true
			}
			log.Info("HttpRoute resource not found. Setting to nil")
			httpRoute = nil
		} else if meta.IsNoMatchError(err) {
			log.Info("HttpRoute CRD not installed. It will not be created")
			httpRoutePresent = false
			httpRoute = nil
		} else {
			log.Error(err, "Failed to get HttpRoute")
			return ctrl.Result{Requeue: true}, err
		}
	} else if !r.isAlreadyWatching {
		r.selfRef.Watch(source.Kind(r.cacheRef, &gatewayv1.HTTPRoute{}), &handler.EnqueueRequestForObject{})
		r.isAlreadyWatching = true
	}

	err = r.Get(ctx, client.ObjectKey{Namespace: mfc.Namespace, Name: mfc.Name}, ingress)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Ingress resource not found. Setting to nil")
			ingress = nil
		} else {
			log.Error(err, "Failed to get Ingress")
			return ctrl.Result{Requeue: true}, err
		}
	}

	// Resource creation logic
	log.Info("Checking if HttpRoute and Ingress resources should be created")
	if mfc.Spec.Routing != nil && httpRoute == nil && ingress == nil {
		operatorService := r.getOperatorService(ctx, log, mfc)
		if operatorService == nil {
			log.Error(err, "Failed to get OperatorService")
			return ctrl.Result{Requeue: true}, err
		}
		if mfc.Spec.Routing.ParentRefs != nil && httpRoutePresent {
			if err := r.createHttpRoute(ctx, log, mfc, operatorService); err != nil {
				if apierrors.IsNotFound(err) {
					log.Info("HttpRoute CRD not installed. It will not be created")
					httpRoutePresent = false
				} else {
					log.Error(err, "Failed to create HttpRoute")
					return ctrl.Result{Requeue: true}, err
				}
			}
			if err := r.recreateRefGrant(ctx, log, mfc, operatorService); err != nil {
				log.Error(err, "Failed to recreate ReferenceGrant")
				return ctrl.Result{Requeue: true}, err
			}
		}
		if mfc.Spec.Routing.IngressClassName != nil {
			if err := r.createIngress(ctx, log, mfc, operatorService); err != nil {
				log.Error(err, "Failed to create Ingress")
				return ctrl.Result{Requeue: true}, err
			}
		}
	}

	// Resource update logic
	log.Info("Checking if HttpRoute and Ingress resources should be updated")
	if mfc.Spec.Routing != nil && httpRoute != nil {
		if mfc.Spec.Routing.ParentRefs == nil && mfc.Spec.Routing.IngressClassName != nil {
			if err := r.Delete(ctx, httpRoute); err != nil {
				log.Error(err, "Failed to delete HttpRoute")
				return ctrl.Result{Requeue: true}, err
			}
			if err := r.createIngress(ctx, log, mfc, r.getOperatorService(ctx, log, mfc)); err != nil {
				log.Error(err, "Failed to create Ingress")
				return ctrl.Result{Requeue: true}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
		operatorService := r.getOperatorService(ctx, log, mfc)
		if operatorService == nil {
			log.Error(err, "Failed to get OperatorService")
			return ctrl.Result{Requeue: true}, err
		}
		update := createRouteForMicroFrontendClass(mfc, operatorService)
		if err := controllerutil.SetControllerReference(mfc, update, r.Scheme); err != nil {
			log.Error(err, "Failed to set controller reference for HttpRoute")
			return ctrl.Result{Requeue: true}, err
		}
		if areNotSameHttpRoute(httpRoute, update) {
			log.Info("Updating HttpRoute for MicroFrontendClass")
			update.ResourceVersion = httpRoute.ResourceVersion
			if err := r.Update(ctx, update); err != nil {
				log.Error(err, "Failed to update HttpRoute")
				return ctrl.Result{Requeue: true}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
	}

	if mfc.Spec.Routing != nil && ingress != nil {
		if mfc.Spec.Routing.IngressClassName == nil && mfc.Spec.Routing.ParentRefs != nil {
			if err := r.Delete(ctx, ingress); err != nil {
				log.Error(err, "Failed to delete Ingress")
				return ctrl.Result{Requeue: true}, err
			}
			if err := r.createHttpRoute(ctx, log, mfc, r.getOperatorService(ctx, log, mfc)); err != nil {
				log.Error(err, "Failed to create HttpRoute")
				return ctrl.Result{Requeue: true}, err
			}
			if err := r.recreateRefGrant(ctx, log, mfc, r.getOperatorService(ctx, log, mfc)); err != nil {
				log.Error(err, "Failed to recreate ReferenceGrant")
				return ctrl.Result{Requeue: true}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
		operatorService := r.getOperatorService(ctx, log, mfc)
		if operatorService == nil {
			log.Error(err, "Failed to get OperatorService")
			return ctrl.Result{Requeue: true}, err
		}
		update := createIngressForMicroFrontendClass(mfc, operatorService)
		controllerutil.SetControllerReference(mfc, update, r.Scheme)
		if areNotSameIngress(ingress, update) {
			log.Info("Updating Ingress for MicroFrontendClass")
			update.ResourceVersion = ingress.ResourceVersion
			if err := r.Update(ctx, update); err != nil {
				log.Error(err, "Failed to update Ingress")
				return ctrl.Result{Requeue: true}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
	}

	// Resource deletion logic
	log.Info("Checking if HttpRoute and Ingress resources should be deleted")
	if mfc.Spec.Routing == nil && httpRoute != nil {
		if err := r.Delete(ctx, httpRoute); err != nil {
			log.Error(err, "Failed to delete HttpRoute")
			return ctrl.Result{Requeue: true}, err
		}
	}
	if mfc.Spec.Routing == nil && ingress != nil {
		if err := r.Delete(ctx, ingress); err != nil {
			log.Error(err, "Failed to delete Ingress")
			return ctrl.Result{Requeue: true}, err
		}
	}

	// Store the MicroFrontendClass in the repository
	if err := r.Repository.Store(mfc); err != nil {
		log.Error(err, "Failed to store MicroFrontendClass in repository")
		return ctrl.Result{Requeue: true}, err
	}

	return ctrl.Result{}, nil
}

func (r *MicroFrontendClassReconciler) finalizeOperationsForMicroFrontendClass(mfc *polyfeav1alpha1.MicroFrontendClass) error {
	log := log.FromContext(context.Background())
	r.Repository.Delete(mfc)
	log.Info("Removing finalizer from MicroFrontendClass", "MicroFrontendClass", mfc)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MicroFrontendClassReconciler) SetupWithManager(mgr ctrl.Manager) error {
	selfRef, err := ctrl.NewControllerManagedBy(mgr).
		For(&polyfeav1alpha1.MicroFrontendClass{}).
		Owns(&networkingv1.Ingress{}).
		Watches(&corev1.Service{}, handler.EnqueueRequestsFromMapFunc(r.findObjectsForService)).
		Build(r)
	if err != nil {
		return err
	}
	r.selfRef = selfRef
	r.cacheRef = mgr.GetCache()
	return nil
}

func (r *MicroFrontendClassReconciler) recreateRefGrant(ctx context.Context, log logr.Logger, mfc *polyfeav1alpha1.MicroFrontendClass, operatorService *corev1.Service) error {
	refGrant := &gatewayv1alpha2.ReferenceGrant{}
	err := r.Get(ctx, client.ObjectKey{Namespace: operatorService.Namespace, Name: mfc.Name}, refGrant)
	if err != nil {
		if apierrors.IsNotFound(err) {
			refGrant = &gatewayv1alpha2.ReferenceGrant{
				ObjectMeta: ctrl.ObjectMeta{
					Name:      mfc.Name,
					Namespace: operatorService.Namespace,
				},
				Spec: gatewayv1alpha2.ReferenceGrantSpec{
					From: []gatewayv1alpha2.ReferenceGrantFrom{
						{
							Group:     "gateway.networking.k8s.io",
							Kind:      "HTTPRoute",
							Namespace: gatewayv1.Namespace(mfc.Namespace),
						},
					},
					To: []gatewayv1alpha2.ReferenceGrantTo{
						{
							Kind: "Service",
						},
					},
				},
			}
		} else {
			log.Error(err, "Failed to get ReferenceGrant")
			return err
		}
	} else {
		if err := r.Delete(ctx, refGrant); err != nil {
			log.Error(err, "Failed to delete ReferenceGrant")
			return err
		}
	}
	if err := r.Create(ctx, refGrant); err != nil {
		log.Error(err, "Failed to create ReferenceGrant")
		return err
	}
	return nil
}

func (r *MicroFrontendClassReconciler) findObjectsForService(ctx context.Context, service client.Object) []reconcile.Request {
	var requests []reconcile.Request
	log := log.FromContext(ctx)

	microFrontendClasses := &polyfeav1alpha1.MicroFrontendClassList{}
	if err := r.List(ctx, microFrontendClasses); err != nil {
		log.Error(err, "Failed to list MicroFrontendClasses")
		return nil
	}

	for _, mfc := range microFrontendClasses.Items {
		if string(mfc.UID) == service.GetAnnotations()["Owner"] {
			requests = append(requests, reconcile.Request{
				NamespacedName: client.ObjectKey{
					Namespace: mfc.Namespace,
					Name:      mfc.Name,
				},
			})
		}
	}
	return requests
}

func areNotSameHttpRoute(httpRoute *gatewayv1.HTTPRoute, update *gatewayv1.HTTPRoute) bool {
	return httpRoute.Name != update.Name ||
		httpRoute.Namespace != update.Namespace ||
		!reflect.DeepEqual(httpRoute.Spec.ParentRefs, update.Spec.ParentRefs) ||
		*httpRoute.Spec.Rules[0].Matches[0].Path.Value != *update.Spec.Rules[0].Matches[0].Path.Value ||
		httpRoute.Spec.Rules[0].BackendRefs[0].BackendRef.Name != update.Spec.Rules[0].BackendRefs[0].BackendRef.Name ||
		*httpRoute.Spec.Rules[0].BackendRefs[0].BackendRef.Port != *update.Spec.Rules[0].BackendRefs[0].BackendRef.Port
}

func areNotSameIngress(ingress *networkingv1.Ingress, update *networkingv1.Ingress) bool {
	return ingress.Name != update.Name ||
		ingress.Namespace != update.Namespace ||
		*ingress.Spec.IngressClassName != *update.Spec.IngressClassName ||
		ingress.Spec.Rules[0].HTTP.Paths[0].Path != update.Spec.Rules[0].HTTP.Paths[0].Path ||
		ingress.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Name != update.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Name ||
		ingress.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Port.Name != update.Spec.Rules[0].HTTP.Paths[0].Backend.Service.Port.Name
}

func (r *MicroFrontendClassReconciler) createHttpRoute(ctx context.Context, log logr.Logger, mfc *polyfeav1alpha1.MicroFrontendClass, operatorService *corev1.Service) error {
	httpRoute := createRouteForMicroFrontendClass(mfc, operatorService)
	if err := controllerutil.SetControllerReference(mfc, httpRoute, r.Scheme); err != nil {
		return err
	}
	return r.Create(ctx, httpRoute)
}

func (r *MicroFrontendClassReconciler) createIngress(ctx context.Context, log logr.Logger, mfc *polyfeav1alpha1.MicroFrontendClass, operatorService *corev1.Service) error {
	ingress := createIngressForMicroFrontendClass(mfc, operatorService)
	if err := controllerutil.SetControllerReference(mfc, ingress, r.Scheme); err != nil {
		return err
	}
	return r.Create(ctx, ingress)
}

func (r *MicroFrontendClassReconciler) getOperatorService(ctx context.Context, log logr.Logger, mfc *polyfeav1alpha1.MicroFrontendClass) *corev1.Service {
	listOptions := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set{OperatorServiceSelectorName: OperatorServiceSelectorValue}),
	}
	var services corev1.ServiceList
	if err := r.List(ctx, &services, listOptions); err != nil || len(services.Items) == 0 {
		log.Error(err, "Failed to list services or no services found")
		return nil
	}
	if services.Items[0].Annotations == nil {
		services.Items[0].Annotations = make(map[string]string)
	}
	services.Items[0].Annotations["Owner"] = string(mfc.UID)
	r.Update(ctx, &services.Items[0])
	return &services.Items[0]
}

func createRouteForMicroFrontendClass(mfc *polyfeav1alpha1.MicroFrontendClass, operatorService *corev1.Service) *gatewayv1.HTTPRoute {
	kind := gatewayv1.Kind("Service")
	name := gatewayv1.ObjectName(operatorService.Name)
	namespace := gatewayv1.Namespace(operatorService.Namespace)
	portIndex := slices.IndexFunc(operatorService.Spec.Ports, func(port corev1.ServicePort) bool {
		return port.Name == PortName
	})
	port := gatewayv1.PortNumber(operatorService.Spec.Ports[portIndex].Port)

	return &gatewayv1.HTTPRoute{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      mfc.Name,
			Namespace: mfc.Namespace,
		},
		Spec: gatewayv1.HTTPRouteSpec{
			CommonRouteSpec: gatewayv1.CommonRouteSpec{
				ParentRefs: mfc.Spec.Routing.ParentRefs,
			},
			Rules: []gatewayv1.HTTPRouteRule{
				{
					Matches: []gatewayv1.HTTPRouteMatch{
						{
							Path: &gatewayv1.HTTPPathMatch{
								Value: mfc.Spec.BaseUri,
							},
						},
					},
					BackendRefs: []gatewayv1.HTTPBackendRef{
						{
							BackendRef: gatewayv1.BackendRef{
								BackendObjectReference: gatewayv1.BackendObjectReference{
									Kind:      &kind,
									Name:      name,
									Port:      &port,
									Namespace: &namespace,
								},
							},
						},
					},
				},
			},
		},
	}
}

func createIngressForMicroFrontendClass(mfc *polyfeav1alpha1.MicroFrontendClass, operatorService *corev1.Service) *networkingv1.Ingress {
	pathType := networkingv1.PathTypePrefix
	return &networkingv1.Ingress{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      mfc.Name,
			Namespace: mfc.Namespace,
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: mfc.Spec.Routing.IngressClassName,
			Rules: []networkingv1.IngressRule{
				{
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     *mfc.Spec.BaseUri,
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
