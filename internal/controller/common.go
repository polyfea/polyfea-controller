package controller

import (
	"context"

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

// FindMicroFrontendClassByName looks up a cluster-scoped MicroFrontendClass by name.
func FindMicroFrontendClassByName(ctx context.Context, c client.Client, name string) (*polyfeav1alpha1.MicroFrontendClass, error) {
	mfc := &polyfeav1alpha1.MicroFrontendClass{}
	if err := c.Get(ctx, client.ObjectKey{Name: name}, mfc); err != nil {
		return nil, err
	}
	return mfc, nil
}

// FindMicroFrontendByName looks up a MicroFrontend by name in the given namespace.
func FindMicroFrontendByName(ctx context.Context, c client.Client, name, namespace string) (*polyfeav1alpha1.MicroFrontend, error) {
	mf := &polyfeav1alpha1.MicroFrontend{}
	if err := c.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, mf); err != nil {
		return nil, err
	}
	return mf, nil
}
