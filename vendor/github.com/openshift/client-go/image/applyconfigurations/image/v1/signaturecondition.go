// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	imagev1 "github.com/openshift/api/image/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SignatureConditionApplyConfiguration represents a declarative configuration of the SignatureCondition type for use
// with apply.
type SignatureConditionApplyConfiguration struct {
	Type               *imagev1.SignatureConditionType `json:"type,omitempty"`
	Status             *corev1.ConditionStatus         `json:"status,omitempty"`
	LastProbeTime      *metav1.Time                    `json:"lastProbeTime,omitempty"`
	LastTransitionTime *metav1.Time                    `json:"lastTransitionTime,omitempty"`
	Reason             *string                         `json:"reason,omitempty"`
	Message            *string                         `json:"message,omitempty"`
}

// SignatureConditionApplyConfiguration constructs a declarative configuration of the SignatureCondition type for use with
// apply.
func SignatureCondition() *SignatureConditionApplyConfiguration {
	return &SignatureConditionApplyConfiguration{}
}

// WithType sets the Type field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Type field is set to the value of the last call.
func (b *SignatureConditionApplyConfiguration) WithType(value imagev1.SignatureConditionType) *SignatureConditionApplyConfiguration {
	b.Type = &value
	return b
}

// WithStatus sets the Status field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Status field is set to the value of the last call.
func (b *SignatureConditionApplyConfiguration) WithStatus(value corev1.ConditionStatus) *SignatureConditionApplyConfiguration {
	b.Status = &value
	return b
}

// WithLastProbeTime sets the LastProbeTime field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LastProbeTime field is set to the value of the last call.
func (b *SignatureConditionApplyConfiguration) WithLastProbeTime(value metav1.Time) *SignatureConditionApplyConfiguration {
	b.LastProbeTime = &value
	return b
}

// WithLastTransitionTime sets the LastTransitionTime field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LastTransitionTime field is set to the value of the last call.
func (b *SignatureConditionApplyConfiguration) WithLastTransitionTime(value metav1.Time) *SignatureConditionApplyConfiguration {
	b.LastTransitionTime = &value
	return b
}

// WithReason sets the Reason field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Reason field is set to the value of the last call.
func (b *SignatureConditionApplyConfiguration) WithReason(value string) *SignatureConditionApplyConfiguration {
	b.Reason = &value
	return b
}

// WithMessage sets the Message field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Message field is set to the value of the last call.
func (b *SignatureConditionApplyConfiguration) WithMessage(value string) *SignatureConditionApplyConfiguration {
	b.Message = &value
	return b
}
