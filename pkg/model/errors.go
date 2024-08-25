package model

import (
	"context"
	"errors"
	"os"
)

// These errors are copied from
// github.com/opencontainers/runc/libcontainer/error.go
var (
	// ErrExist indicates that a container with the given ID already exists.
	ErrExist     = errors.New("container with given ID already exists")
	ErrExistExit = 101

	// ErrInvalidID indicates that the provided container ID format is invalid.
	ErrInvalidID     = errors.New("invalid container ID format")
	ErrInvalidIDExit = 102

	// ErrNotExist indicates that the specified container does not exist.
	ErrNotExist     = errors.New("container does not exist")
	ErrNotExistExit = 103

	// ErrRunning indicates that the container is still running and cannot perform the requested operation.
	ErrRunning     = errors.New("container still running")
	ErrRunningExit = 104

	// ErrNotRunning indicates that the container is not currently running.
	ErrNotRunning     = errors.New("container not running")
	ErrNotRunningExit = 105
)

var (
	// ErrNoSudo indicates that the runtime caller has no superuser privileges (sudo)
	ErrNoSudo     = errors.New("runtime needs to be run as sudo")
	ErrNoSudoExit = 10

	ErrFileNotExistExit    = 3
	ErrFileExistExit       = 31
	ErrContextCanceledExit = 2
	UnknownErrorExit       = 1
)

// ExitCode maps an error to a return code
func ExitCode(err error) int {
	switch {
	case errors.Is(err, ErrExist):
		return ErrExistExit
	case errors.Is(err, ErrInvalidID):
		return ErrInvalidIDExit
	case errors.Is(err, ErrNotExist):
		return ErrNotExistExit
	case errors.Is(err, ErrRunning):
		return ErrRunningExit
	case errors.Is(err, ErrNotRunning):
		return ErrNotRunningExit
	case errors.Is(err, ErrNoSudo):
		return ErrNoSudoExit
	case os.IsNotExist(err):
		return ErrFileNotExistExit
	case os.IsExist(err):
		return ErrFileExistExit
	case errors.Is(err, context.Canceled):
		return ErrContextCanceledExit
	default:
		return UnknownErrorExit
	}
}
