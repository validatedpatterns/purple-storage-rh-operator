package rbac_script_test

import (
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
			wantPanic: true,
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

	if p.Kind != "Deployment" {
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
	wantSubstr := "//+kubebuilder:rbac:groups=\"\",resources=Pods,namespace=default,verbs=create;delete;get;list;patch;update;watch"
	if rules[0] != wantSubstr {
		t.Errorf("Expected: %q\nGot:      %q", wantSubstr, rules[0])
	}
}
