package namespace

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"syscall"
)

// Time namespace
type timeNS struct {
}

func newTimeNamespace() *timeNS {
	return new(timeNS)
}

func (t *timeNS) IsSupported() bool {
	return true
}

func (t *timeNS) Priority() int {
	return 0
}

func (t *timeNS) Type() specs.LinuxNamespaceType {
	return specs.TimeNamespace
}

func (t *timeNS) CloneFlag() uintptr {
	return syscall.CLONE_NEWTIME
}

func (t *timeNS) Finalize(spec specs.Spec) error {
	return nil
}
