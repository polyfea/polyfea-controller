package repository

import (
	"github.com/polyfea/polyfea-controller/api/v1alpha1"
)

type MicrofrontendFilterFunc func(mf v1alpha1.MicroFrontend) bool
type MicrofrontendClassFilterFunc func(mf v1alpha1.MicroFrontendClass) bool
type WebComponentFilterFunc func(mf v1alpha1.WebComponent) bool

type PolyfeaRepository interface {
	StoreMicrofrontend(microfrontend v1alpha1.MicroFrontend) error
	StoreMicrofrontendClass(microfrontendClass v1alpha1.MicroFrontendClass) error
	StoreWebComponent(webComponent v1alpha1.WebComponent) error

	GetMicrofrontends(filter MicrofrontendFilterFunc) ([]v1alpha1.MicroFrontend, error)
	GetMicrofrontendClasses(filter MicrofrontendClassFilterFunc) ([]v1alpha1.MicroFrontendClass, error)
	GetWebComponents(filter WebComponentFilterFunc) ([]v1alpha1.WebComponent, error)
}

type InMemoryPolyfeaRepository struct {
	microfrontedns       map[string]v1alpha1.MicroFrontend
	microfrontendClasses map[string]v1alpha1.MicroFrontendClass
	webComponents        map[string]v1alpha1.WebComponent
}

func NewInMemoryPolyfeaRepository() *InMemoryPolyfeaRepository {
	return &InMemoryPolyfeaRepository{
		microfrontedns:       map[string]v1alpha1.MicroFrontend{},
		microfrontendClasses: map[string]v1alpha1.MicroFrontendClass{},
		webComponents:        map[string]v1alpha1.WebComponent{},
	}
}

func (r *InMemoryPolyfeaRepository) StoreMicrofrontend(microfrontend v1alpha1.MicroFrontend) error {
	r.microfrontedns[microfrontend.Name] = microfrontend
	return nil
}

func (r *InMemoryPolyfeaRepository) StoreMicrofrontendClass(microfrontendClass v1alpha1.MicroFrontendClass) error {
	r.microfrontendClasses[microfrontendClass.Name] = microfrontendClass
	return nil
}

func (r *InMemoryPolyfeaRepository) StoreWebComponent(webComponent v1alpha1.WebComponent) error {
	r.webComponents[webComponent.Name] = webComponent
	return nil
}

func (r *InMemoryPolyfeaRepository) GetMicrofrontends(filter MicrofrontendFilterFunc) ([]v1alpha1.MicroFrontend, error) {
	var result []v1alpha1.MicroFrontend
	for _, mf := range r.microfrontedns {
		if filter(mf) {
			result = append(result, mf)
		}
	}
	return result, nil
}

func (r *InMemoryPolyfeaRepository) GetMicrofrontendClasses(filter MicrofrontendClassFilterFunc) ([]v1alpha1.MicroFrontendClass, error) {
	var result []v1alpha1.MicroFrontendClass
	for _, mf := range r.microfrontendClasses {
		if filter(mf) {
			result = append(result, mf)
		}
	}
	return result, nil
}

func (r *InMemoryPolyfeaRepository) GetWebComponents(filter WebComponentFilterFunc) ([]v1alpha1.WebComponent, error) {
	var result []v1alpha1.WebComponent
	for _, mf := range r.webComponents {
		if filter(mf) {
			result = append(result, mf)
		}
	}
	return result, nil
}
