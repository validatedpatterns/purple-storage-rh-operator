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

package controller

import (
	"context"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	purplev1alpha1 "github.com/darkdoc/purple-storage-rh-operator/api/v1alpha1"
	mfc "github.com/manifestival/controller-runtime-client"
	"github.com/manifestival/manifestival"
)

// PurpleStorageReconciler reconciles a PurpleStorage object
type PurpleStorageReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	config        *rest.Config
	dynamicClient dynamic.Interface
}

// Basic Operator RBACs
//+kubebuilder:rbac:groups=purple.purplestorage.com,resources=purplestorages,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=purple.purplestorage.com,resources=purplestorages/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=purple.purplestorage.com,resources=purplestorages/finalizers,verbs=update

// Operator needs to create some machine configs
//+kubebuilder:rbac:groups=machineconfiguration.openshift.io,resources=machineconfigs,verbs=get;list;watch;create;update;patch;delete

// Below rules are generated via ./scripts/create-rbac.sh
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings,verbs=list;watch;delete;update;get;create;patch
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=list;watch;delete;update;get;create;patch
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=list;watch;delete;update;get;create;patch
//+kubebuilder:rbac:groups=apiextensions.k8s.io,resources=customresourcedefinitions,verbs=list;watch;delete;update;get;create;patch
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=list;watch;delete;update;get;create;patch
//+kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=mutatingwebhookconfigurations,verbs=list;watch;delete;update;get;create;patch
//+kubebuilder:rbac:groups="",resources=namespaces,verbs=list;watch;delete;update;get;create;patch
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=list;watch;delete;update;get;create;patch
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=list;watch;delete;update;get;create;patch
//+kubebuilder:rbac:groups=security.openshift.io,resources=securitycontextconstraints,verbs=list;watch;delete;update;get;create;patch
//+kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=list;watch;delete;update;get;create;patch
//+kubebuilder:rbac:groups="",resources=services,verbs=list;watch;delete;update;get;create;patch
//+kubebuilder:rbac:groups=admissionregistration.k8s.io,resources=validatingwebhookconfigurations,verbs=list;watch;delete;update;get;create;patch

// There is a bug in IBMs role granting silly permissions, so here we work around that:
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=list;watch;delete;update;get;create;patch
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=configmaps,verbs=list;watch;delete;update;get;create;patch
//+kubebuilder:rbac:groups="",resources=leases,verbs=list;watch;delete;update;get;create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PurpleStorage object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *PurpleStorageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	purplestorage := &purplev1alpha1.PurpleStorage{}
	err := r.Get(ctx, req.NamespacedName, purplestorage)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	install_path := fmt.Sprintf("files/%s/install.yaml", purplestorage.Spec.Ibm_spectrum_scale_container_native_version)
	_, err = os.Stat(install_path)
	if os.IsNotExist(err) {
		install_path = fmt.Sprintf("/%s", install_path)
		_, err = os.Stat(install_path)
		if os.IsNotExist(err) {
			return ctrl.Result{}, err
		}
	}
	installManifest, err := manifestival.NewManifest(install_path, manifestival.UseClient(mfc.NewClient(r.Client)))
	if err != nil {
		return ctrl.Result{}, err
	}
	log.Log.Info(fmt.Sprintf("Applying manifest from %s", install_path))

	if err := installManifest.Apply(); err != nil {
		log.Log.Error(err, "Error applying manifest")
		return reconcile.Result{}, err
	}
	log.Log.Info(fmt.Sprintf("Applied manifest from %s", install_path))

	new_mc := NewMachineConfig(purplestorage.Spec.Machineconfig.Labels)
	gvr := schema.GroupVersionResource{
		Group:    "machineconfiguration.openshift.io",
		Version:  "v1",
		Resource: "machineconfigs",
	}

	old_mc, err := r.dynamicClient.Resource(gvr).Get(ctx, new_mc.GetName(), metav1.GetOptions{})
	if err != nil {
		log.Log.Info(fmt.Sprintf("Creating machineconfig"))
		err = r.Client.Create(ctx, new_mc)
		if err != nil {
			return ctrl.Result{}, err
		}
		log.Log.Info(fmt.Sprintf("Created machineconfig"))

	} else {
		log.Log.Info(fmt.Sprintf("Updating machineconfig"))
		new_mc.SetResourceVersion(old_mc.GetResourceVersion())
		err = r.Client.Update(ctx, new_mc)
		if err != nil {
			return ctrl.Result{}, err
		}
		log.Log.Info(fmt.Sprintf("Updated machineconfig"))
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PurpleStorageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	var err error
	r.config = mgr.GetConfig()
	if r.dynamicClient, err = dynamic.NewForConfig(r.config); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&purplev1alpha1.PurpleStorage{}).
		Complete(r)
}
