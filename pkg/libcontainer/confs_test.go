package libcontainer

import (
	"testing"
)

func TestValidateId(t *testing.T) {
	tests := []struct {
		id    string
		valid bool
	}{
		{"", false},        // Empty string
		{"abc", true},      // Valid ID
		{"ABC123", true},   // Valid ID
		{"abc_123", true},  // Valid ID with underscore
		{"abc-123", true},  // Valid ID with minus
		{"abc+123", true},  // Valid ID with plus
		{"abc.123", true},  // Valid ID with period
		{".", false},       // Invalid ID: single dot
		{"..", false},      // Invalid ID: double dots
		{"abc@123", false}, // Invalid ID: invalid character '@'
		{"abc$123", false}, // Invalid ID: invalid character '$'
		{"abc 123", false}, // Invalid ID: space
		{"abc/123", false}, // Invalid ID: invalid character '/'
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			valid := validateIdFormat(tt.id)
			if valid != tt.valid {
				t.Errorf("validateID(%q), expected: %v, actual: %v", tt.id, tt.valid, valid)
			}
		})
	}

}
