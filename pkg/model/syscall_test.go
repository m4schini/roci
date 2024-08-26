package model

import (
	"syscall"
	"testing"
)

func TestSyscallSignal(t *testing.T) {
	tests := []struct {
		name        string
		signalName  string
		wantSignal  syscall.Signal
		expectError bool
	}{
		{
			name:       "Numeric signal",
			signalName: "9",
			wantSignal: syscall.Signal(9),
		},
		{
			name:       "Symbolic signal without SIG prefix",
			signalName: "INT",
			wantSignal: syscall.Signal(2),
		},
		{
			name:       "Symbolic signal with SIG prefix",
			signalName: "SIGTERM",
			wantSignal: syscall.Signal(15),
		},
		{
			name:        "Unknown signal name",
			signalName:  "SIGUNKNOWN",
			expectError: true,
		},
		{
			name:        "Empty signal name",
			signalName:  "",
			expectError: true,
		},
		{
			name:        "Non-numeric, non-symbolic signal",
			signalName:  "abc",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSignal, err := SyscallSignal(tt.signalName)
			if (err != nil) != tt.expectError {
				t.Errorf("SyscallSignal(%v) error = %v, expectError %v", tt.signalName, err, tt.expectError)
				return
			}
			if !tt.expectError && gotSignal != tt.wantSignal {
				t.Errorf("SyscallSignal(%v) = %v, want %v", tt.signalName, gotSignal, tt.wantSignal)
			}
		})
	}
}
