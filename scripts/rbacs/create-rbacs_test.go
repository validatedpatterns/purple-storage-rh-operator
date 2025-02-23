package rbac_script

import (
	"fmt"
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
`

// Test parsing normal Kubernetes objects
func TestExtractRBACRules(t *testing.T) {
	rules, _ := ExtractRBACRules([]byte(sampleYAML))
	fmt.Printf("Rules ZOZZO1: %d - %v\n", len(rules), rules)
	defaultVerbs := NewStringSet()
	for _, s := range []string{"get", "list", "watch", "create", "update", "patch", "delete"} {
		defaultVerbs.Add(s)
	}
	expectedRules := map[schema.GroupVersionResource]StringSet{
		{Group: "", Version: "v1", Resource: "pods"}: defaultVerbs}
	// {Group: "apps", Version: "v1", Resource: "deployments"}:
	// 	NewStringSetFromList([]string{"get", "list", "watch", "create", "update", "patch", "delete"}),
	// ,
	// {Group: "core", Version: "v1", Resource: "pods"}:
	// 	"get", "watch", "list",
	// },
	// {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "roles"}: {
	// 	"get", "list", "watch", "create", "update", "patch", "delete",
	// },
	// {Group: "rbac.authorization.k8s.io", Version: "v1", Resource: "clusterroles"}: {
	// 	"get", "list", "watch", "create", "update", "patch", "delete",
	// },
	// {Group: "apps", Version: "v1", Resource: "deployments"}: {
	// 	"create", "delete",
	// },

	fmt.Printf("Rules ZOZZO2: %d - %v\n", len(expectedRules), expectedRules)

	// Verify if the extracted rules match expected values
	for expectedGVR, expectedVerbs := range expectedRules {
		fmt.Printf("Expected %v: %v", expectedGVR, expectedVerbs)
		actualVerbs, exists := rules[expectedGVR]
		fmt.Printf("Rule: %v", actualVerbs)

		if !exists {
			t.Errorf("Missing RBAC rule for: %+v", expectedGVR)
			continue
		}

		if !actualVerbs.Equals(expectedVerbs) {
			t.Errorf("Mismatch for %+v\nExpected: %v\nGot: %v", expectedGVR, expectedVerbs, actualVerbs)
		}
	}
}

// Test RBAC marker generation
func TestGenerateRBACMarkers(t *testing.T) {
	rules, _ := ExtractRBACRules([]byte(sampleYAML))
	markers := GenerateRBACMarkers(rules)

	expectedMarkers := []string{
		"//+kubebuilder:rbac:groups=core,resources=pods,verbs=create,delete,get,list,patch,update,watch",
	}

	for _, expected := range expectedMarkers {
		found := false
		t.Logf("Expect: %s", expected)
		for _, line := range markers {
			t.Logf("Marker: %s", line)
			if strings.TrimSpace(line) == strings.TrimSpace(expected) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected RBAC marker not found:\n%s\nFound:\n%s\n", expected, markers)
		}
	}
}
