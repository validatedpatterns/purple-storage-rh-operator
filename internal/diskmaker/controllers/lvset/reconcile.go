package lvset

import (
	"context"
	"fmt"
	"time"

	localv1alpha1 "github.com/darkdoc/purple-storage-rh-operator/api/v1alpha1"
	"github.com/darkdoc/purple-storage-rh-operator/internal/common"
	"github.com/darkdoc/purple-storage-rh-operator/internal/diskmaker"
	internal "github.com/darkdoc/purple-storage-rh-operator/internal/diskutils"
	"github.com/darkdoc/purple-storage-rh-operator/internal/localmetrics"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	provCommon "sigs.k8s.io/sig-storage-local-static-provisioner/pkg/common"
	provDeleter "sigs.k8s.io/sig-storage-local-static-provisioner/pkg/deleter"
)

const (
	// ComponentName for lvset symlinker
	ComponentName      = "localvolumeset-symlink-controller"
	pvOwnerKey         = "pvOwner"
	defaultRequeueTime = time.Minute
	fastRequeueTime    = 5 * time.Second
)

var nodeName string
var watchNamespace string

func init() {
	nodeName = common.GetNodeNameEnvVar()
	watchNamespace, _ = common.GetWatchNamespace()
}

// Reconcile reads that state of the cluster for a LocalVolumeSet object and makes changes based on the state read
// and what is in the LocalVolumeSet.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *LocalVolumeSetReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	requeueTime := defaultRequeueTime

	// Fetch the LocalVolumeSet instance
	lvset := &localv1alpha1.PurpleStorage{}
	err := r.Client.Get(ctx, request.NamespacedName, lvset)
	if err != nil {
		if kerrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	klog.InfoS("Reconciling LocalVolumeSet", "namespace", request.Namespace, "name", request.Name)

	err = common.ReloadRuntimeConfig(ctx, r.Client, request, r.nodeName, r.runtimeConfig)
	if err != nil {
		return ctrl.Result{}, err
	}

	// ignore LocalVolmeSets whose LabelSelector doesn't match this node
	// NodeSelectorTerms.MatchExpressions are ORed
	matches, err := common.NodeSelectorMatchesNodeLabels(r.runtimeConfig.Node, lvset.Spec.NodeSpec.Selector)
	if err != nil {
		klog.ErrorS(err, "failed to match nodeSelector to node labels")
		return ctrl.Result{}, err
	}

	if !matches {
		return ctrl.Result{}, nil
	}

	klog.InfoS("Looking for valid block devices", "namespace", request.Namespace, "name", request.Name)
	// list block devices
	blockDevices, badRows, err := internal.ListBlockDevices([]string{})
	if err != nil {
		msg := fmt.Sprintf("failed to list block devices: %v", err)
		r.eventReporter.Report(lvset, newDiskEvent(diskmaker.ErrorRunningBlockList, msg, "", corev1.EventTypeWarning))
		klog.Error(msg)
		return ctrl.Result{}, err
	} else if len(badRows) > 0 {
		msg := fmt.Sprintf("error parsing rows: %+v", badRows)
		r.eventReporter.Report(lvset, newDiskEvent(diskmaker.ErrorRunningBlockList, msg, "", corev1.EventTypeWarning))
		klog.Error(msg)
	}

	// find disks that match lvset filters and matchers
	validDevices, delayedDevices := r.getValidDevices(lvset, blockDevices)

	// update metrics for unmatched disks
	localmetrics.SetLVSUnmatchedDiskMetric(nodeName, len(blockDevices)-len(validDevices))

	// shorten the requeueTime if there are delayed devices
	if len(delayedDevices) > 1 && requeueTime == defaultRequeueTime {
		requeueTime = deviceMinAge / 2
	}

	return ctrl.Result{Requeue: true, RequeueAfter: requeueTime}, nil
}

// runs filters and matchers on the blockDeviceList and returns valid devices
// and devices that are not considered old enough to be valid yet
// i.e. if the device is younger than deviceMinAge
// if the waitingDevices list is nonempty, the operator should requeueue
func (r *LocalVolumeSetReconciler) getValidDevices(lvset *localv1alpha1.PurpleStorage, blockDevices []internal.BlockDevice) ([]internal.BlockDevice, []internal.BlockDevice) {
	validDevices := make([]internal.BlockDevice, 0)
	delayedDevices := make([]internal.BlockDevice, 0)
	// get valid devices
DeviceLoop:
	for _, blockDevice := range blockDevices {

		// store device in deviceAgeMap
		r.deviceAgeMap.storeDeviceAge(blockDevice.KName)

		for name, filter := range FilterMap {
			var valid bool
			var err error
			valid, err = filter(blockDevice, nil)
			if err != nil {
				klog.ErrorS(err, "filter error", "device",
					blockDevice.Name, "filter", name)
				valid = false
				continue DeviceLoop
			} else if !valid {
				klog.InfoS("filter negative", "device",
					blockDevice.Name, "filter", name)
				continue DeviceLoop
			}
		}

		// check if the device is older than deviceMinAge
		isOldEnough := r.deviceAgeMap.isOlderThan(blockDevice.KName)

		// skip devices younger than deviceMinAge
		if !isOldEnough {
			delayedDevices = append(delayedDevices, blockDevice)
			// record DiscoveredDevice event
			if lvset != nil {
				r.eventReporter.Report(
					lvset,
					newDiskEvent(
						DiscoveredNewDevice,
						fmt.Sprintf("found possible matching disk, waiting %v to claim", deviceMinAge),
						blockDevice.KName, corev1.EventTypeNormal,
					),
				)
			}
			continue DeviceLoop
		}

		for name, matcher := range matcherMap {
			valid, err := matcher(blockDevice, lvset.Spec.NodeSpec.DeviceInclusionSpec)
			if err != nil {
				klog.ErrorS(err, "match error", "device",
					blockDevice.Name, "filter", name)
				valid = false
				continue DeviceLoop
			} else if !valid {
				klog.InfoS("match negative", "device",
					blockDevice.Name, "filter", name)
				continue DeviceLoop
			}
		}
		klog.InfoS("matched disk", "device", blockDevice.Name)
		// handle valid disk
		validDevices = append(validDevices, blockDevice)

	}
	return validDevices, delayedDevices
}

type LocalVolumeSetReconciler struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	Client        client.Client
	Scheme        *runtime.Scheme
	nodeName      string
	eventReporter *eventReporter
	cacheSynced   bool
	// map from KNAME of device to time when the device was first observed since the process started
	deviceAgeMap *ageMap

	// static-provisioner stuff
	cleanupTracker *provDeleter.CleanupStatusTracker
	runtimeConfig  *provCommon.RuntimeConfig
	deleter        *provDeleter.Deleter
}

func NewLocalVolumeSetReconciler(client client.Client, scheme *runtime.Scheme, time timeInterface, cleanupTracker *provDeleter.CleanupStatusTracker, rc *provCommon.RuntimeConfig) *LocalVolumeSetReconciler {
	deleter := provDeleter.NewDeleter(rc, cleanupTracker)
	eventReporter := newEventReporter(rc.Recorder)
	lvsReconciler := &LocalVolumeSetReconciler{
		Client:         client,
		Scheme:         scheme,
		nodeName:       nodeName,
		eventReporter:  eventReporter,
		deviceAgeMap:   newAgeMap(time),
		cleanupTracker: cleanupTracker,
		runtimeConfig:  rc,
		deleter:        deleter,
	}

	return lvsReconciler
}

func (r *LocalVolumeSetReconciler) WithManager(mgr ctrl.Manager) error {
	err := ctrl.NewControllerManagedBy(mgr).
		// set to 1 explicitly, despite it being the default, as the reconciler is not thread-safe.
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		For(&localv1alpha1.PurpleStorage{}).
		Complete(r)

	return err
}
