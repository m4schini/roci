package model

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestExitCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{
			name:     "ErrExist",
			err:      ErrExist,
			wantCode: ErrExistExit,
		},
		{
			name:     "ErrInvalidID",
			err:      ErrInvalidID,
			wantCode: ErrInvalidIDExit,
		},
		{
			name:     "ErrNotExist",
			err:      ErrNotExist,
			wantCode: ErrNotExistExit,
		},
		{
			name:     "ErrRunning",
			err:      ErrRunning,
			wantCode: ErrRunningExit,
		},
		{
			name:     "ErrNotRunning",
			err:      ErrNotRunning,
			wantCode: ErrNotRunningExit,
		},
		{
			name:     "ErrNoSudo",
			err:      ErrNoSudo,
			wantCode: ErrNoSudoExit,
		},
		{
			name:     "os.IsNotExist",
			err:      os.ErrNotExist,
			wantCode: ErrFileNotExistExit,
		},
		{
			name:     "os.IsExist",
			err:      os.ErrExist,
			wantCode: ErrFileExistExit,
		},
		{
			name:     "context.Canceled",
			err:      context.Canceled,
			wantCode: ErrContextCanceledExit,
		},
		{
			name:     "Unknown error",
			err:      errors.New("unknown error"),
			wantCode: UnknownErrorExit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExitCode(tt.err); got != tt.wantCode {
				t.Errorf("ExitCode(%v) = %d, want %d", tt.err, got, tt.wantCode)
			}
		})
	}
}
