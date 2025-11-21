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

package v1alpha1

import (
	"testing"
)

func TestServiceReference_ResolveServiceURL(t *testing.T) {
	tests := []struct {
		name             string
		serviceRef       *ServiceReference
		defaultNamespace string
		expectedURL      string
	}{
		{
			name:             "URI is set - should return URI directly",
			serviceRef:       &ServiceReference{URI: ptr("https://cdn.example.com")},
			defaultNamespace: "default",
			expectedURL:      "https://cdn.example.com",
		},
		{
			name:             "URI with path - should return full URI",
			serviceRef:       &ServiceReference{URI: ptr("http://external-service.com/api")},
			defaultNamespace: "default",
			expectedURL:      "http://external-service.com/api",
		},
		{
			name: "In-cluster service with all fields",
			serviceRef: &ServiceReference{
				Name:      ptr("my-service"),
				Namespace: ptr("custom-namespace"),
				Port:      ptr(int32(8080)),
				Scheme:    ptr("https"),
			},
			defaultNamespace: "default",
			expectedURL:      "https://my-service.custom-namespace.svc.cluster.local:8080",
		},
		{
			name: "In-cluster service with defaults",
			serviceRef: &ServiceReference{
				Name: ptr("my-service"),
			},
			defaultNamespace: "default",
			expectedURL:      "http://my-service.default.svc.cluster.local:80",
		},
		{
			name: "In-cluster service using default namespace",
			serviceRef: &ServiceReference{
				Name:   ptr("my-service"),
				Port:   ptr(int32(3000)),
				Scheme: ptr("http"),
			},
			defaultNamespace: "production",
			expectedURL:      "http://my-service.production.svc.cluster.local:3000",
		},
		{
			name: "In-cluster service with custom port",
			serviceRef: &ServiceReference{
				Name:      ptr("api-service"),
				Namespace: ptr("backend"),
				Port:      ptr(int32(9000)),
			},
			defaultNamespace: "default",
			expectedURL:      "http://api-service.backend.svc.cluster.local:9000",
		},
		{
			name: "In-cluster service with https scheme",
			serviceRef: &ServiceReference{
				Name:      ptr("secure-service"),
				Namespace: ptr("default"),
				Scheme:    ptr("https"),
			},
			defaultNamespace: "default",
			expectedURL:      "https://secure-service.default.svc.cluster.local:80",
		},
		{
			name:             "Nil service reference",
			serviceRef:       nil,
			defaultNamespace: "default",
			expectedURL:      "",
		},
		{
			name:             "Empty service reference",
			serviceRef:       &ServiceReference{},
			defaultNamespace: "default",
			expectedURL:      "",
		},
		{
			name: "In-cluster service with custom domain",
			serviceRef: &ServiceReference{
				Name:      ptr("my-service"),
				Namespace: ptr("default"),
				Port:      ptr(int32(80)),
				Scheme:    ptr("http"),
				Domain:    ptr("svc.cluster.example.com"),
			},
			defaultNamespace: "default",
			expectedURL:      "http://my-service.default.svc.cluster.example.com:80",
		},
		{
			name: "In-cluster service with custom domain and custom port",
			serviceRef: &ServiceReference{
				Name:      ptr("api-service"),
				Namespace: ptr("backend"),
				Port:      ptr(int32(9000)),
				Scheme:    ptr("https"),
				Domain:    ptr("custom.local"),
			},
			defaultNamespace: "default",
			expectedURL:      "https://api-service.backend.custom.local:9000",
		},
		{
			name: "In-cluster service with empty domain string uses default",
			serviceRef: &ServiceReference{
				Name:      ptr("my-service"),
				Namespace: ptr("default"),
				Domain:    ptr(""),
			},
			defaultNamespace: "default",
			expectedURL:      "http://my-service.default.svc.cluster.local:80",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.serviceRef.ResolveServiceURL(tt.defaultNamespace)
			if result != tt.expectedURL {
				t.Errorf("ResolveServiceURL() = %v, want %v", result, tt.expectedURL)
			}
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
