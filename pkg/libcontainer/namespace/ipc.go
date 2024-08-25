package namespace

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"syscall"
)

// IPC Namespace
type ipcNS struct {
}

func newIpcNamespace() *ipcNS {
	return new(ipcNS)
}

func (i *ipcNS) IsSupported() bool {
	return true
}

func (i *ipcNS) Priority() int {
	return 0
}

func (i *ipcNS) Type() specs.LinuxNamespaceType {
	return specs.IPCNamespace
}

func (i *ipcNS) CloneFlag() uintptr {
	return syscall.CLONE_NEWIPC
}

func (i *ipcNS) Finalize(spec specs.Spec) error {
	return nil
}
