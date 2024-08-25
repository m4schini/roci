package namespace

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"go.uber.org/zap"
	"syscall"
)

// Network Namespace
type netNS struct {
	logger *zap.Logger
}

func newNetworkNamespace() *netNS {
	ns := new(netNS)
	ns.logger = log.Named("net")
	return ns
}

func (n *netNS) IsSupported() bool {
	// Network isolation is not implemented in roci
	return false
}

func (n *netNS) Priority() int {
	return 0
}

func (n *netNS) Type() specs.LinuxNamespaceType {
	return specs.NetworkNamespace
}

func (n *netNS) CloneFlag() uintptr {
	logNsNotImplemented(n.logger)
	return syscall.CLONE_NEWNET
}

func (n *netNS) Finalize(spec specs.Spec) error {
	logNsNotImplemented(n.logger)
	return nil
}
