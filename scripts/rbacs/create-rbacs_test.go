package rbac_script

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Sample YAML input for testing
const sampleYAML = `
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
  namespace: default
spec:
  containers:
    - name: nginx
      image: nginx
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
        - name: nginx
          image: nginx
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: example-role
  namespace: default
rules:
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get", "watch", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: example-clusterrole
rules:
  - apiGroups: ["apps"]
    resources: ["deployments"]
    verbs: ["create", "delete"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  labels:
    app.kubernetes.io/instance: ibm-spectrum-scale
    app.kubernetes.io/name: operator
  name: ibm-spectrum-scale-leader-election-role
  namespace: ibm-spectrum-scale-operator
rules:
- apiGroups:
  - ""
  - coordination.k8s.io
  resources:
  - configmaps
  - leases
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
  - patch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  labels:
    app.kubernetes.io/instance: ibm-spectrum-scale
    app.kubernetes.io/name: cluster
  name: ibm-spectrum-scale-sysmon
  namespace: ibm-spectrum-scale
rules:
- apiGroups:
  - ""
  resources:
  - pods
  - services
  verbs:
  - get
  - list
- apiGroups:
  - apps
  resources:
  - deployments
  - statefulsets
  verbs:
  - get
  - list
---
`

// Test parsing normal Kubernetes objects
func TestExtractRBACRules(t *testing.T) {
	rules, _ := ExtractRBACRules([]byte(sampleYAML))

	expectedRules := map[schema.GroupVersionResource][]string{
		{Group: "core", Version: "v1", Resource: "pods"}: {
			"get", "list", "watch", "create", "update", "patch", "delete",
		},
		{Group: "apps", Version: "v1", Resource: "deployments"}: {
			"get", "list", "watch", "create", "update", "patch", "delete",
		},
		{Group: "core", Version: "v1", Resource: "pods"}: {
			"get", "watch", "list",
		},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles"}: {
			"get", "list", "watch", "create", "update", "patch", "delete",
		},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterroles"}: {
			"get", "list", "watch", "create", "update", "patch", "delete",
		},
		{Group: "apps", Version: "v1", Resource: "deployments"}: {
			"create", "delete",
		},
	}

	// Verify if the extracted rules match expected values
	for expectedGVR, expectedVerbs := range expectedRules {
		actualVerbs, exists := rules[expectedGVR]
		if !exists {
			t.Errorf("Missing RBAC rule for: %+v", expectedGVR)
			continue
		}

		if !equalStringSlices(actualVerbs, expectedVerbs) {
			t.Errorf("Mismatch for %+v\nExpected: %v\nGot: %v", expectedGVR, expectedVerbs, actualVerbs)
		}
	}
}

// Helper function to compare two string slices
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	exists := make(map[string]bool)
	for _, v := range a {
		exists[v] = true
	}
	for _, v := range b {
		if !exists[v] {
			return false
		}
	}
	return true
}

// Test RBAC marker generation
func TestGenerateRBACMarkers(t *testing.T) {
	rules, _ := ExtractRBACRules([]byte(sampleYAML))

	var output bytes.Buffer
	oldStdout := captureStdout(&output)
	defer restoreStdout(oldStdout)

	GenerateRBACMarkers(rules)

	expectedMarkers := []string{
		"// +kubebuilder:rbac:groups=core,resources=pods,verbs=get,list,watch,create,update,patch,delete",
		"// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get,list,watch,create,update,patch,delete",
		"// +kubebuilder:rbac:groups=core,resources=pods,verbs=get,watch,list",
		"// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get,list,watch,create,update,patch,delete",
		"// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=get,list,watch,create,update,patch,delete",
		"// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=create,delete",
	}

	outputLines := strings.Split(strings.TrimSpace(output.String()), "\n")

	for _, expected := range expectedMarkers {
		found := false
		for _, line := range outputLines {
			if line == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected RBAC marker not found:\n%s", expected)
		}
	}
}

// Helper function to capture stdout
func captureStdout(output *bytes.Buffer) *bytes.Buffer {
	old := bytes.NewBuffer([]byte{})
	old.Write(output.Bytes())
	output.Reset()
	return old
}

// Helper function to restore stdout
func restoreStdout(old *bytes.Buffer) {
	old.WriteTo(os.Stdout)
}
