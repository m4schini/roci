package namespace

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"go.uber.org/zap"
	"syscall"
)

// Cgroup Namespace
type cgroupNS struct {
	log *zap.Logger
}

func newCgroupNamespace() *cgroupNS {
	ns := new(cgroupNS)
	ns.log = log.Named("cgroup")
	return ns
}

func (c *cgroupNS) IsSupported() bool {

	return false
}

func (c *cgroupNS) Priority() int {
	return 0
}

func (c *cgroupNS) Type() specs.LinuxNamespaceType {
	return specs.CgroupNamespace
}

func (c *cgroupNS) CloneFlag() uintptr {
	logNsNotImplemented(c.log)
	return syscall.CLONE_NEWCGROUP
}

func (c *cgroupNS) Finalize(spec specs.Spec) error {
	logNsNotImplemented(c.log)
	return nil
}
