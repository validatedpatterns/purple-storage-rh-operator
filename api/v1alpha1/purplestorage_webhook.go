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
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var purplestoragelog = logf.Log.WithName("purplestorage-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *PurpleStorage) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&PurpleStorage{}).
		WithValidator(&PurpleStorageValidator{}).
		Complete()
}

// +kubebuilder:object:generate=false
// +k8s:deepcopy-gen=false
// +k8s:openapi-gen=false
// PurpleStorageValidator is responsible for setting default values on the PurpleStorage resources
// when created or updated.
//
// NOTE: The +kubebuilder:object:generate=false and +k8s:deepcopy-gen=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type PurpleStorageValidator struct{}

// FIXME(bandini): This needs to be reviewed more in detail. I added sideEffects=none to get it passing but not 100% sure about it
// +kubebuilder:webhook:verbs=create;update,path=/validate-purple-purplestorage-com-v1alpha1-purplestorage,mutating=false,failurePolicy=fail,groups=purple.purplestorage.com,resources=purplestorages,versions=v1alpha1,name=vpurplestorage.kb.io,admissionReviewVersions=v1,sideEffects=none

var _ webhook.CustomValidator = &PurpleStorageValidator{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *PurpleStorageValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	p, ok := obj.(*PurpleStorage)
	if !ok {
		return nil, fmt.Errorf("expected a PurpleStorage object but got %T", obj)
	}
	purplestoragelog.Info("validate create", "name", p.Name)

	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *PurpleStorageValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	p, ok := oldObj.(*PurpleStorage)
	if !ok {
		return nil, fmt.Errorf("expected a PurpleStorage object but got %T", oldObj)
	}
	purplestoragelog.Info("validate update", "name", p.Name)

	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *PurpleStorageValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	p, ok := obj.(*PurpleStorage)
	if !ok {
		return nil, fmt.Errorf("expected a PurpleStorage object but got %T", obj)
	}
	purplestoragelog.Info("validate create", "name", p.Name)

	return nil, nil
}
