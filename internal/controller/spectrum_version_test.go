package controller

import "testing"

func TestIsOpenShiftSupported(t *testing.T) {
	tests := []struct {
		ibmVersion string
		ocpVersion string
		expected   bool
	}{
		{"5.1.5.0", "4.9", true},          // Expected to be supported
		{"5.1.5.0", "4.12", false},        // Not in the supported list
		{"5.1.7.0", "4.11", true},         // Supported
		{"5.1.7.0", "4.13", false},        // Not supported
		{"5.1.9.1", "4.12", true},         // Supported
		{"5.1.9.1", "4.15", false},        // Not supported
		{"5.2.2.0", "4.17", true},         // Supported
		{"5.2.2.0", "4.18", false},        // Not supported
		{"5.2.2.0", "4.15", true},         // Supported
		{"invalid_version", "4.9", false}, // Invalid IBM Storage Scale version
	}

	for _, tt := range tests {
		result := IsOpenShiftSupported(tt.ibmVersion, tt.ocpVersion)
		if result != tt.expected {
			t.Errorf("IsOpenShiftSupported(%s, %s) = %v; expected %v", tt.ibmVersion, tt.ocpVersion, result, tt.expected)
		}
	}
}
