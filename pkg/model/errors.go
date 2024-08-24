package model

import (
	"context"
	"errors"
	"os"
)

var (
	ErrExist      = errors.New("container with given ID already exists")
	ErrInvalidID  = errors.New("invalid container ID format")
	ErrNotExist   = errors.New("container does not exist")
	ErrPaused     = errors.New("container paused")
	ErrRunning    = errors.New("container still running")
	ErrNotRunning = errors.New("container not running")
	ErrNotPaused  = errors.New("container not paused")
)

// These errors are copied from
// github.com/opencontainers/runc/libcontainer/error.go
// TODO License and copyright notice?

var (
	ErrNoSudo = errors.New("runtime needs to be run as sudo")
)

func ExitCode(err error) int {
	switch {
	case errors.Is(err, ErrExist):
		return 1001
	case errors.Is(err, ErrInvalidID):
		return 1002
	case errors.Is(err, ErrNotExist):
		return 1003
	case errors.Is(err, ErrNoSudo):
		return 10
	case os.IsNotExist(err):
		return 3
	case os.IsExist(err):
		return 31
	case errors.Is(err, context.Canceled):
		return 2
	default:
		return 1
	}
}
