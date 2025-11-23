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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Common condition types following Kubernetes conventions
const (
	// ConditionTypeReady indicates that the resource is ready for use
	ConditionTypeReady string = "Ready"

	// ConditionTypeAccepted indicates that the resource has been accepted
	ConditionTypeAccepted string = "Accepted"

	// ConditionTypeAvailable indicates that the resource is available
	ConditionTypeAvailable string = "Available"
)

// MicroFrontend-specific condition types
const (
	// ConditionTypeServiceResolved indicates that the service reference has been resolved
	ConditionTypeServiceResolved string = "ServiceResolved"

	// ConditionTypeFrontendClassBound indicates that the MicroFrontend is bound to a MicroFrontendClass
	ConditionTypeFrontendClassBound string = "FrontendClassBound"

	// ConditionTypeNamespacePolicyValid indicates that the MicroFrontend satisfies the namespace policy
	ConditionTypeNamespacePolicyValid string = "NamespacePolicyValid"
)

// WebComponent-specific condition types
const (
	// ConditionTypeMicroFrontendResolved indicates that the referenced MicroFrontend has been resolved
	ConditionTypeMicroFrontendResolved string = "MicroFrontendResolved"
)

// Common condition reasons
const (
	// ReasonSuccessful indicates a successful operation
	ReasonSuccessful string = "Successful"

	// ReasonInvalidConfiguration indicates an invalid configuration
	ReasonInvalidConfiguration string = "InvalidConfiguration"

	// ReasonNamespaceNotAllowed indicates that the namespace is not allowed by policy
	ReasonNamespaceNotAllowed string = "NamespaceNotAllowed"

	// ReasonFrontendClassNotFound indicates that the referenced MicroFrontendClass was not found
	ReasonFrontendClassNotFound string = "FrontendClassNotFound"

	// ReasonServiceNotFound indicates that the referenced Service was not found
	ReasonServiceNotFound string = "ServiceNotFound"

	// ReasonMicroFrontendNotFound indicates that the referenced MicroFrontend was not found
	ReasonMicroFrontendNotFound string = "MicroFrontendNotFound"

	// ReasonReconciling indicates that the resource is being reconciled
	ReasonReconciling string = "Reconciling"

	// ReasonError indicates that an error occurred during reconciliation
	ReasonError string = "Error"
)

// Phase constants for MicroFrontend
const (
	MicroFrontendPhasePending  string = "Pending"
	MicroFrontendPhaseReady    string = "Ready"
	MicroFrontendPhaseFailed   string = "Failed"
	MicroFrontendPhaseRejected string = "Rejected"
)

// Phase constants for WebComponent
const (
	WebComponentPhasePending               string = "Pending"
	WebComponentPhaseReady                 string = "Ready"
	WebComponentPhaseFailed                string = "Failed"
	WebComponentPhaseMicroFrontendNotFound string = "MicroFrontendNotFound"
)

// Phase constants for MicroFrontendClass
const (
	MicroFrontendClassPhaseReady   string = "Ready"
	MicroFrontendClassPhaseInvalid string = "Invalid"
)

// SetCondition adds or updates a condition in the conditions list
func SetCondition(conditions *[]metav1.Condition, conditionType string, status metav1.ConditionStatus, reason, message string) {
	now := metav1.Now()
	newCondition := metav1.Condition{
		Type:               conditionType,
		Status:             status,
		ObservedGeneration: 0, // This should be set by the caller if needed
		LastTransitionTime: now,
		Reason:             reason,
		Message:            message,
	}

	// Find if the condition already exists
	for i, condition := range *conditions {
		if condition.Type == conditionType {
			// Only update if the status has changed
			if condition.Status != status {
				newCondition.LastTransitionTime = now
			} else {
				newCondition.LastTransitionTime = condition.LastTransitionTime
			}
			(*conditions)[i] = newCondition
			return
		}
	}

	// Condition doesn't exist, append it
	*conditions = append(*conditions, newCondition)
}

// RemoveCondition removes a condition from the conditions list
func RemoveCondition(conditions *[]metav1.Condition, conditionType string) {
	for i, condition := range *conditions {
		if condition.Type == conditionType {
			*conditions = append((*conditions)[:i], (*conditions)[i+1:]...)
			return
		}
	}
}

// GetCondition returns the condition with the given type
func GetCondition(conditions []metav1.Condition, conditionType string) *metav1.Condition {
	for i := range conditions {
		if conditions[i].Type == conditionType {
			return &conditions[i]
		}
	}
	return nil
}

// IsConditionTrue returns true if the condition is present and set to true
func IsConditionTrue(conditions []metav1.Condition, conditionType string) bool {
	condition := GetCondition(conditions, conditionType)
	return condition != nil && condition.Status == metav1.ConditionTrue
}

// IsReady returns true if the Ready condition is present and set to true
func IsReady(conditions []metav1.Condition) bool {
	return IsConditionTrue(conditions, ConditionTypeReady)
}
