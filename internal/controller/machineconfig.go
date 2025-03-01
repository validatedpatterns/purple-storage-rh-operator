package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	machineconfigv1 "github.com/openshift/api/machineconfiguration/v1"
	ctrlcommon "github.com/openshift/machine-config-operator/pkg/controller/common"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// apiVersion: machineconfiguration.openshift.io/v1
// kind: MachineConfig
// metadata:
//   labels:
//     machineconfiguration.openshift.io/role: worker
//   name: 00-worker-ibm-spectrum-scale-kernel-devel
// spec:
//   config:
//     ignition:
//       version: 3.2.0
//   extensions:
//   - kernel-devel

func NewMachineConfig(labels map[string]string) *machineconfigv1.MachineConfig {
	return &machineconfigv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: machineconfigv1.SchemeGroupVersion.String(),
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "00-worker-ibm-spectrum-scale-kernel-devel",
			Labels: labels,
		},
		Spec: *NewMachineConfigSpec(),
	}
}

func NewMachineConfigSpec() *machineconfigv1.MachineConfigSpec {
	tmpIgnCfg := ctrlcommon.NewIgnConfig()
	rawTmpIgnCfg, _ := json.Marshal(tmpIgnCfg)

	return &machineconfigv1.MachineConfigSpec{
		// config is a Ignition Config object.
		// +optional
		Config: runtime.RawExtension{
			Raw: rawTmpIgnCfg,
		},
		// extensions contains a list of additional features that can be enabled on host
		// +listType=atomic
		// +optional
		Extensions: []string{"kernel-devel"},
	}
}

// WaitForMachineConfigPoolUpdated polls the MachineConfigPool until it shows Updated=True
func WaitForMachineConfigPoolUpdated(ctx context.Context, client dynamic.Interface, mcpName string) error {
	mcpGVR := schema.GroupVersionResource{
		Group:    "machineconfiguration.openshift.io",
		Version:  "v1",
		Resource: "machineconfigpools",
	}

	// 1. Get the latest MachineConfigPool object
	mcp, err := client.Resource(mcpGVR).Get(ctx, mcpName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get MachineConfigPool %q: %w", mcpName, err)
	}

	// 2. Check if the MCP is "Updated"
	updated, reason, err := isMachineConfigPoolUpdating(mcp)
	if err != nil {
		return fmt.Errorf("failed to parse MCP conditions: %w", err)
	}

	if updated {
		fmt.Printf("MachineConfigPool %q has completed (Updated=True)\n", mcpName)
		return nil
	}

	return fmt.Errorf("MachineConfigPool %q not updated yet. Reason: %s\n", mcpName, reason)
}

// FIXME(bandini): For now we check for the conditions and if one is of type Updated and status True we return true.
// Maybe we should actually check that machineCount == readyMachineCount == updatedMachineCount ?
// conditions:
//   - lastTransitionTime: "2025-03-01T11:31:08Z"
//     message: ""
//     reason: ""
//     status: "False"
//     type: RenderDegraded
//   - lastTransitionTime: "2025-03-01T11:31:13Z"
//     message: ""
//     reason: ""
//     status: "False"
//     type: NodeDegraded
//   - lastTransitionTime: "2025-03-01T11:31:13Z"
//     message: ""
//     reason: ""
//     status: "False"
//     type: Degraded
//   - lastTransitionTime: "2025-03-01T12:02:56Z"
//     message: All nodes are updated with MachineConfig rendered-worker-7056a3d79e377146ee42af553e10ee68
//     reason: ""
//     status: "True"
//     type: Updated
//   - lastTransitionTime: "2025-03-01T12:02:56Z"
//     message: ""
//     reason: ""
//     status: "False"
//     type: Updating
//
// isMachineConfigPoolUpdated checks the MCP's status conditions to see if it has Updated=True
func isMachineConfigPoolUpdating(mcp *unstructured.Unstructured) (bool, string, error) {
	conditions, found, err := unstructured.NestedSlice(mcp.Object, "status", "conditions")
	if err != nil {
		return false, "", err
	}
	if !found {
		return false, "no conditions found", nil
	}

	// Parse each condition and see if type=Updated with status=True
	for _, c := range conditions {
		cond, ok := c.(map[string]interface{})
		if !ok {
			return false, "", errors.New("condition is not in expected map format")
		}

		condType, _, _ := unstructured.NestedString(cond, "type")
		condStatus, _, _ := unstructured.NestedString(cond, "status")
		condReason, _, _ := unstructured.NestedString(cond, "reason")

		if condType == "Updated" && condStatus == "True" {
			return true, condReason, nil
		}
	}

	// Not updated yet
	return false, "Updated=False", nil
}
