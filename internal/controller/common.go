package controller

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	polyfeav1alpha1 "github.com/polyfea/polyfea-controller/api/v1alpha1"
)

const (
	PortName                 = "webserver"
	DefaultFrontendClassName = "polyfea-controller-default"
	FinalizerName            = "polyfea.github.io/finalizer"
)

// GetFrontendClassName returns the frontend class name for a MicroFrontend spec,
// falling back to DefaultFrontendClassName if not specified.
func GetFrontendClassName(frontendClass *string) string {
	if frontendClass != nil && *frontendClass != "" {
		return *frontendClass
	}
	return DefaultFrontendClassName
}

func ptr[T any](v T) *T { return &v }

// FindMicroFrontendClassByName looks up a MicroFrontendClass by name, first in the given
// namespace then across all namespaces. Returns the object and nil error if found, or
// a not-found/other error if not.
func FindMicroFrontendClassByName(ctx context.Context, c client.Client, name, namespace string) (*polyfeav1alpha1.MicroFrontendClass, error) {
	mfc := &polyfeav1alpha1.MicroFrontendClass{}
	err := c.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, mfc)
	if err == nil {
		return mfc, nil
	}
	if !apierrors.IsNotFound(err) {
		return nil, err
	}

	// Search across all namespaces
	mfcList := &polyfeav1alpha1.MicroFrontendClassList{}
	if listErr := c.List(ctx, mfcList, client.InNamespace("")); listErr != nil {
		return nil, listErr
	}
	for i := range mfcList.Items {
		if mfcList.Items[i].Name == name {
			return &mfcList.Items[i], nil
		}
	}
	return nil, err // return original not-found error
}

// FindMicroFrontendByName looks up a MicroFrontend by name, first in the given
// namespace then across all namespaces. Returns the object and nil error if found, or
// a not-found/other error if not.
func FindMicroFrontendByName(ctx context.Context, c client.Client, name, namespace string) (*polyfeav1alpha1.MicroFrontend, error) {
	mf := &polyfeav1alpha1.MicroFrontend{}
	err := c.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, mf)
	if err == nil {
		return mf, nil
	}
	if !apierrors.IsNotFound(err) {
		return nil, err
	}

	// Search across all namespaces
	mfList := &polyfeav1alpha1.MicroFrontendList{}
	if listErr := c.List(ctx, mfList); listErr != nil {
		return nil, listErr
	}
	for i := range mfList.Items {
		if mfList.Items[i].Name == name {
			return &mfList.Items[i], nil
		}
	}
	return nil, err // return original not-found error
}
