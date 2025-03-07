package nodedaemon

import (
	"context"
	"strings"
	"time"

	localv1alpha1 "github.com/darkdoc/purple-storage-rh-operator/api/v1alpha1"
	"github.com/darkdoc/purple-storage-rh-operator/internal/common"
	"github.com/darkdoc/purple-storage-rh-operator/internal/localmetrics"
	prometheusv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	errorutils "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	oldProvisionerName     = "localvolumeset-local-provisioner"
	oldLVDiskMakerPrefix   = "local-volume-diskmaker-"
	oldLVProvisionerPrefix = "local-volume-provisioner-"
	appLabelKey            = "app"
	// ProvisionerName is the name of the local-static-provisioner daemonset
	ProvisionerName = "local-provisioner"
	// DiskMakerName is the name of the diskmaker-manager daemonset
	DiskMakerName = "diskmaker-manager"

	// dataHashAnnotationKey = "purple.purplestorage.com/configMapDataHash"

	orphanLSOServiceMonitorName = "local-storage-operator-metrics"
)

// DaemonReconciler reconciles all LocalVolumeSet objects at once
type DaemonReconciler struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	Client                        client.Client
	Scheme                        *runtime.Scheme
	deletedStaticProvisioner      bool
	deletedOrphanedServiceMonitor bool
}

// Reconcile reads that state of the cluster for a PurpleStorage object and makes changes based on the state read
// and what is in the PurpleStorage.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *DaemonReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	err := r.cleanupOldObjects(ctx, request.Namespace)
	if err != nil {
		return ctrl.Result{}, err
	}

	tolerations, ownerRefs, nodeSelector, err := r.aggregateDeamonInfo(ctx, request)
	if err != nil {
		return ctrl.Result{}, err
	}

	// enable service and servicemonitor for diskmaker daemonset
	metricsExportor := localmetrics.NewExporter(ctx, r.Client, common.DiskMakerServiceName, request.Namespace, common.DiskMakerMetricsServingCert,
		ownerRefs, DiskMakerName)
	if err := metricsExportor.EnableMetricsExporter(); err != nil {
		klog.ErrorS(err, "failed to create service and servicemonitors for diskmaker daemonset")
		return ctrl.Result{}, err
	}

	if err := localmetrics.CreateOrUpdateAlertRules(ctx, r.Client, request.Namespace, DiskMakerName, ownerRefs); err != nil {
		klog.ErrorS(err, "failed to create alerting rules")
		return ctrl.Result{}, err
	}

	diskMakerDSMutateFn := getDiskMakerDSMutateFn(request, tolerations, ownerRefs, nodeSelector)
	ds, opResult, err := CreateOrUpdateDaemonset(ctx, r.Client, diskMakerDSMutateFn)
	if err != nil {
		return ctrl.Result{}, err
	} else if opResult == controllerutil.OperationResultUpdated || opResult == controllerutil.OperationResultCreated {
		klog.InfoS("daemonset changed", "dsName", ds.GetName(), "opResult", opResult)
	}

	return ctrl.Result{}, err
}

// do a one-time delete of the old object that are not needed in this release
func (r *DaemonReconciler) cleanupOldObjects(ctx context.Context, namespace string) error {
	errs := make([]error, 0)
	err := r.cleanupOldDaemonsets(ctx, namespace)
	if err != nil {
		errs = append(errs, err)
	}
	err = r.cleanupOrphanServiceMonitor(ctx, namespace)
	if err != nil {
		errs = append(errs, err)
	}
	return errorutils.NewAggregate(errs)
}

// do a one-time delete of the old static-provisioner daemonset
func (r *DaemonReconciler) cleanupOldDaemonsets(ctx context.Context, namespace string) error {
	if r.deletedStaticProvisioner {
		return nil
	}

	// search for old localvolume daemons
	dsList := &appsv1.DaemonSetList{}
	err := r.Client.List(ctx, dsList, client.InNamespace(namespace))
	if err != nil {
		klog.ErrorS(err, "could not list daemonsets")
		return err
	}
	appNameList := make([]string, 0)
	for _, ds := range dsList.Items {
		appLabel, found := ds.ObjectMeta.Labels[appLabelKey]
		if !found {
			continue
		} else if strings.HasPrefix(appLabel, oldLVDiskMakerPrefix) || strings.HasPrefix(appLabel, oldLVProvisionerPrefix) {
			// remember name to watch for pods to delete
			appNameList = append(appNameList, appLabel)
			// delete daemonset
			err = r.Client.Delete(ctx, &ds)
			if err != nil && !(errors.IsNotFound(err) || errors.IsGone(err)) {
				klog.ErrorS(err, "could not delete daemonset", "dsName", ds.Name)
				return err
			}
		}
	}

	// wait for pods to die
	err = wait.ExponentialBackoff(wait.Backoff{
		Cap:      time.Minute * 2,
		Duration: time.Second,
		Factor:   1.7,
		Jitter:   1,
		Steps:    20,
	}, func() (done bool, err error) {
		podList := &corev1.PodList{}
		allGone := false
		// search for any pods with label 'app' in appNameList
		appNameList = append(appNameList, oldProvisionerName)
		requirement, err := labels.NewRequirement(appLabelKey, selection.In, appNameList)
		if err != nil {
			klog.ErrorS(err, "failed to compose labelselector requirement",
				"appLabelKey", appLabelKey, "appNameList", appNameList)
			return false, err
		}
		selector := labels.NewSelector().Add(*requirement)
		err = r.Client.List(context.TODO(), podList, client.MatchingLabelsSelector{Selector: selector})
		if err != nil && !errors.IsNotFound(err) {
			return false, err
		} else if len(podList.Items) == 0 {
			allGone = true
		}
		klog.Infof("waiting for 0 pods with app label %q, found = %v",
			oldProvisionerName, len(podList.Items))
		return allGone, nil
	})
	if err != nil {
		klog.ErrorS(err, "could not determine that old provisioner pods were deleted")
		return err
	}
	r.deletedStaticProvisioner = true
	return nil
}

// do a one-time delete of the old LSO ServiceMonitor from 4.8
func (r *DaemonReconciler) cleanupOrphanServiceMonitor(ctx context.Context, namespace string) error {
	if r.deletedOrphanedServiceMonitor {
		return nil
	}

	obj := prometheusv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      orphanLSOServiceMonitorName,
		},
	}

	err := r.Client.Delete(ctx, &obj)
	if err != nil && !errors.IsNotFound(err) {
		klog.V(2).ErrorS(err, "Could not delete ServiceMonitor")
		return err
	}

	klog.V(2).InfoS("Orphaned ServiceMonitor deleted", "namespace", obj.Namespace, "name", obj.Name)
	r.deletedOrphanedServiceMonitor = true
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DaemonReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// The controller will ignore the name part of the enqueued request as
	// every reconcile gathers multiple resources an acts on a few one-per-namespace objects.
	enqueueOnlyNamespace := handler.EnqueueRequestsFromMapFunc(
		func(ctx context.Context, obj client.Object) []reconcile.Request {
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{Namespace: obj.GetNamespace()},
			}
			return []reconcile.Request{req}
		})

	return ctrl.NewControllerManagedBy(mgr).
		For(&localv1alpha1.PurpleStorage{}).
		// watch provisioner, diskmaker-manager daemonsets
		Watches(&appsv1.DaemonSet{}, enqueueOnlyNamespace, builder.WithPredicates(common.EnqueueOnlyLabeledSubcomponents(DiskMakerName, ProvisionerName))).
		Complete(r)
}
