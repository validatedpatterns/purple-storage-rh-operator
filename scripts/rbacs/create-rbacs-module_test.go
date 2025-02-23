package rbac_script_test

import (
	"fmt"
	"reflect"
	"testing"

	rbac_script "github.com/darkdoc/purple-storage-rh-operator/scripts/rbacs"
)

func TestAddStringUnique(t *testing.T) {
	tests := []struct {
		name      string
		initial   []string
		value     string
		want      []string
		wantPanic bool
	}{
		{
			name:      "add to empty slice",
			initial:   []string{},
			value:     "apple",
			want:      []string{"apple"},
			wantPanic: false,
		},
		{
			name:      "add new element",
			initial:   []string{"apple", "banana"},
			value:     "orange",
			want:      []string{"apple", "banana", "orange"},
			wantPanic: false,
		},
		{
			name:      "duplicate element -> panic",
			initial:   []string{"apple", "banana"},
			value:     "banana",
			want:      []string{"apple", "banana"},
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.wantPanic && r == nil {
					t.Errorf("Expected panic, but did not panic")
				} else if !tt.wantPanic && r != nil {
					t.Errorf("Did not expect panic, but got one: %v", r)
				}
			}()

			got := rbac_script.AddStringUnique(tt.initial, tt.value)
			if !tt.wantPanic && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddStringUnique() = %v; want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToPlural(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Empty string",
			input: "",
			want:  "",
		},
		{
			name:  "Already plural",
			input: "pods",
			want:  "pods",
		},
		{
			name:  "Singular: Pod",
			input: "Pod",
			want:  "Pods",
		},
		{
			name:  "Singular: Deployment",
			input: "Deployment",
			want:  "Deployments",
		},
		{
			name:  "Ends with 's': Days",
			input: "Days",
			want:  "Days", // remains unchanged
		},
		{
			name:  "Ends with 's' (lowercase): roles",
			input: "roles",
			want:  "roles", // remains unchanged
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rbac_script.ConvertToPlural(tt.input)
			if got != tt.want {
				t.Errorf("ConvertToPlural(%q) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewPermission(t *testing.T) {
	raw := map[string]interface{}{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata": map[string]interface{}{
			"name":      "my-deploy",
			"namespace": "default",
		},
	}
	p := rbac_script.NewPermission(raw)
	if p == nil {
		t.Fatal("Expected non-nil Permission")
	}

	if p.Kind != "deployment" {
		t.Errorf("Expected kind=Deployment, got %s", p.Kind)
	}
	if p.Group != "apps" {
		t.Errorf("Expected group=apps, got %s", p.Group)
	}
	if p.Version != "v1" {
		t.Errorf("Expected version=v1, got %s", p.Version)
	}
	if p.Name != "my-deploy" {
		t.Errorf("Expected name=my-deploy, got %s", p.Name)
	}
	if p.Namespace != "default" {
		t.Errorf("Expected namespace=default, got %s", p.Namespace)
	}
}

func TestPermission_RBACRule_NonRole(t *testing.T) {
	raw := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Pod",
		"metadata": map[string]interface{}{
			"name":      "my-pod",
			"namespace": "default",
		},
	}
	p := rbac_script.NewPermission(raw)
	if p == nil {
		t.Fatal("NewPermission returned nil")
	}
	rules := p.RBACRule()

	if len(rules) != 1 {
		t.Fatalf("Expected 1 RBAC rule, got %d", len(rules))
	}

	// The rule should have default verbs: get, list, watch, create, update, patch, delete
	wantSubstr := "//+kubebuilder:rbac:groups=\"\",namespace=default,resources=pods,verbs=create;delete;get;list;patch;update;watch"
	if rules[0] != wantSubstr {
		t.Errorf("Expected: %q\nGot:      %q", wantSubstr, rules[0])
	}
}

func TestPermission_RBACRule_Role(t *testing.T) {
	roleYAML := map[string]interface{}{
		"apiVersion": "rbac.authorization.k8s.io/v1",
		"kind":       "Role",
		"metadata": map[string]interface{}{
			"name":      "example-role",
			"namespace": "default",
		},
		"rules": []interface{}{
			map[string]interface{}{
				"apiGroups": []interface{}{""},
				"resources": []interface{}{"pods"},
				"verbs":     []interface{}{"get", "list", "watch"},
			},
			map[string]interface{}{
				"apiGroups": []interface{}{"apps"},
				"resources": []interface{}{"deployments"},
				"verbs":     []interface{}{"create", "delete"},
			},
		},
	}

	p := rbac_script.NewPermission(roleYAML)
	if p == nil {
		t.Fatal("Expected a valid Permission for Role")
	}
	rules := p.RBACRule()
	// We expect 2 lines of output (one for pods, one for deployments)
	if len(rules) != 2 {
		t.Fatalf("Expected 2 lines of RBAC markers, got %d", len(rules))
	}

	// Check for pods rule
	wantPods := "//+kubebuilder:rbac:groups=\"\",namespace=default,resources=pods,verbs=get;list;watch"
	if rules[0] != wantPods && rules[1] != wantPods {
		t.Errorf("Could not find pods RBAC rule in output:\n%v", rules)
	}

	// Check for deployments rule
	wantDeploy := "//+kubebuilder:rbac:groups=apps,namespace=default,resources=deployments,verbs=create;delete"
	if rules[0] != wantDeploy && rules[1] != wantDeploy {
		t.Errorf("Could not find deployments RBAC rule in output:\n%v", rules)
	}
}

func TestIntegration_ExtractAndGenerateMarkers(t *testing.T) {
	yamlContent := `
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
  namespace: default
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deploy
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: example-role
  namespace: default
rules:
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get", "list", "watch"]
  - apiGroups: ["apps"]
    resources: ["deployments"]
    verbs: ["create", "delete"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: example-clusterrole
rules:
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["list", "watch", "create"]
`

	perms, err := rbac_script.ExtractRBACRules([]byte(yamlContent))
	if err != nil {
		t.Fatalf("ExtractRBACRules error: %v", err)
	}
	if len(perms) != 4 {
		t.Fatalf("Expected 4 objects from YAML, got %d", len(perms))
	}

	rbacs := rbac_script.GenerateRBACMarkers(perms)

	// We expect lines for:
	// Pod (default verbs)
	// Deployment (default verbs)
	// Role (2 lines)
	// ClusterRole (1+ lines depending on rules)
	// Let's do a minimal check that some lines appear
	checks := []string{
		"//+kubebuilder:rbac:groups=\"\",namespace=default,resources=pods,verbs=create;delete;get;list;patch;update;watch",
		"//+kubebuilder:rbac:groups=apps,namespace=default,resources=deployments,verbs=create;delete;get;list;patch;update;watch",
		"//+kubebuilder:rbac:groups=\"\",namespace=default,resources=pods,verbs=get;list;watch",
		"//+kubebuilder:rbac:groups=apps,namespace=default,resources=deployments,verbs=create;delete",
		"//+kubebuilder:rbac:groups=\"\",resources=secrets,verbs=create;list;watch", // clusterrole
	}
	if len(checks) != len(rbacs) {
		t.Fatalf("Expected %d RBAC markers, got %d", len(checks), len(rbacs))
	}
	for i := range rbacs {
		fmt.Printf("RBAC : %s\n", rbacs[i])
		fmt.Printf("CHECK: %v\n", checks[i])
	}

	for i := 0; i < len(checks); i++ {
		if rbacs[i] != checks[i] {
			t.Errorf("Expected %s got %s\n", checks[i], rbacs[i])
		} else {
			t.Logf("SAME: %s\n", checks[i])
		}
	}
}
