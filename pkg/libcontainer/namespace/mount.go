package namespace

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"go.uber.org/zap"
	"syscall"
)

type mountNS struct {
	log *zap.Logger
}

func (m *mountNS) IsSupported() bool {
	return true
}

func newMountNamespace() *mountNS {
	ns := new(mountNS)
	ns.log = log.Named("mnt")
	return ns
}

func (m *mountNS) Priority() int {
	return 1
}

func (m *mountNS) Type() specs.LinuxNamespaceType {
	return specs.MountNamespace
}

func (m *mountNS) CloneFlag() uintptr {
	return syscall.CLONE_NEWNS
}

func (m *mountNS) Finalize(spec specs.Spec) error {
	return nil
}
