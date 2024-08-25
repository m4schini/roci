package namespace

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"go.uber.org/zap"
	"syscall"
)

// PID namespace
type pidNS struct {
	logger *zap.Logger
	spec   specs.Spec
}

func newPidNamespace(spec specs.Spec) *pidNS {
	return &pidNS{spec: spec, logger: zap.L().Named("namespace").Named("pid")}
}

func (p *pidNS) IsSupported() bool {
	return true
}

func (p *pidNS) Priority() int {
	return 2
}

func (p *pidNS) Type() specs.LinuxNamespaceType {
	return specs.PIDNamespace
}

func (p *pidNS) CloneFlag() uintptr {
	return syscall.CLONE_NEWPID
}

func (p *pidNS) Finalize(spec specs.Spec) error {
	// the pid namespace needs to be finalized by mounting the procfs in the container root filesystem.
	// But this is not implemented here, instead its part of the rootfs.FinalizeRootfs
	return nil
}
