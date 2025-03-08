package api

import (
	purplev1alpha1 "github.com/validatedpatterns/purple-storage-rh-operator/api/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, purplev1alpha1.SchemeBuilder.AddToScheme)
}
