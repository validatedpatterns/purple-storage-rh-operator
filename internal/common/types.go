package common

const (
	// LocalVolumeOwnerNameForPV stores name of LocalVolume that created this PV
	LocalVolumeOwnerNameForPV = "storage.openshift.com/local-volume-owner-name"
	// LocalVolumeOwnerNamespaceForPV stores namespace of LocalVolume that created this PV
	LocalVolumeOwnerNamespaceForPV = "storage.openshift.com/local-volume-owner-namespace"

	// PVOwnerKindLabel stores the namespace of the CR that created this PV
	PVOwnerKindLabel = "storage.openshift.com/owner-kind"
	// PVOwnerNameLabel stores the name of the CR that created this PV
	PVOwnerNameLabel = "storage.openshift.com/owner-name"
	// PVOwnerNamespaceLabel stores the namespace of the CR that created this PV
	PVOwnerNamespaceLabel = "storage.openshift.com/owner-namespace"
	// PVDeviceNameLabel is the KNAME of the device
	PVDeviceNameLabel = "storage.openshift.com/device-name"
	// PVDeviceIDLabel is the id of the device
	PVDeviceIDLabel = "storage.openshift.com/device-id"
)

// DeprecatedLabels: these labels were deprecated because the potential values weren't all compatible label values
// they have been move to annotations
var DeprecatedLabels = []string{PVDeviceNameLabel, PVDeviceIDLabel}
