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

func ptr[T any](v T) *T { return &v }

// FindMicroFrontendClassByName looks up a MicroFrontendClass by name in the specified namespace.
func FindMicroFrontendClassByName(ctx context.Context, c client.Client, name, namespace string) (*polyfeav1alpha1.MicroFrontendClass, error) {
	mfc := &polyfeav1alpha1.MicroFrontendClass{}
	err := c.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, mfc)
	if err != nil {
		return nil, err
	}
	return mfc, nil
}

// FindMicroFrontendByName looks up a MicroFrontend by name in the specified namespace.
func FindMicroFrontendByName(ctx context.Context, c client.Client, name, namespace string) (*polyfeav1alpha1.MicroFrontend, error) {
	mf := &polyfeav1alpha1.MicroFrontend{}
	err := c.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, mf)
	if err != nil {
		return nil, err
	}
	return mf, nil
}
