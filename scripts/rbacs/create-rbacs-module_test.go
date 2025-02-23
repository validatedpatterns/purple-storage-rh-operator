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
