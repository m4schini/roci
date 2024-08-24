package namespace

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"go.uber.org/zap"
	"roci/pkg/libcontainer/ipc"
	"roci/pkg/logger"
	"slices"
	"syscall"
)

var log = logger.Log().Named("ns")

type Namespace interface {
	Priority() int
	Type() specs.LinuxNamespaceType
	CloneFlag() uintptr
	Finalize(spec specs.Spec) error
	IsSupported() bool
}

func From(runtime ipc.RuntimePipeWriter, spec specs.Spec) (namespaces []Namespace, err error) {
	log.Debug("Preparing namespaces from spec")

	namespacesSpec := spec.Linux.Namespaces
	namespaces = make([]Namespace, len(namespacesSpec))
	for i, namespace := range namespacesSpec {
		namespaces[i] = fromNamespaceType(runtime, spec, namespace.Type)
		log.Debug("prepared namespace", zap.Any("type", namespace.Type))
	}

	slices.SortFunc(namespaces, func(a, b Namespace) int {
		return b.Priority() - b.Priority()
	})

	return namespaces, nil
}

func fromNamespaceType(runtime ipc.RuntimePipeWriter, spec specs.Spec, namespaceType specs.LinuxNamespaceType) Namespace {
	switch namespaceType {
	case specs.PIDNamespace:
		return newPidNamespace(spec)
	case specs.IPCNamespace:
		return newIpcNamespace()
	case specs.TimeNamespace:
		return newTimeNamespace()
	case specs.UTSNamespace:
		return newUtsNamespace(spec)
	case specs.NetworkNamespace:
		return newNetworkNamespace()
	case specs.MountNamespace:
		return newMountNamespace()
	case specs.CgroupNamespace:
		return newCgroupNamespace()
	case specs.UserNamespace:
		return newUserNamespace(runtime)
	default:
		return nil
	}
}

func IsSupported(namespaceType specs.LinuxNamespaceType) bool {
	ns := fromNamespaceType(nil, specs.Spec{}, namespaceType)
	if ns == nil {
		return false
	}
	return ns.IsSupported()
}

func cloneFlag(namespaceType specs.LinuxNamespaceType) uintptr {
	switch namespaceType {
	case specs.PIDNamespace:
		return syscall.CLONE_NEWPID
	case specs.CgroupNamespace:
		return syscall.CLONE_NEWCGROUP
	case specs.MountNamespace:
		return syscall.CLONE_NEWNS
	case specs.IPCNamespace:
		return syscall.CLONE_NEWIPC
	case specs.NetworkNamespace:
		return syscall.CLONE_NEWNET
	case specs.TimeNamespace:
		return syscall.CLONE_NEWTIME
	case specs.UTSNamespace:
		return syscall.CLONE_NEWUTS
	case specs.UserNamespace:
		return syscall.CLONE_NEWUSER
	default:
		panic("unknown namespace type")
	}
}

func Unshare(namespace Namespace) error {
	return syscall.Unshare(int(namespace.CloneFlag()))
}

func logNsNotImplemented(log *zap.Logger) {
	log.Warn("namespace not implemented")
}
