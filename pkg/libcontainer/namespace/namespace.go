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

// Namespace defines an interface for Linux namespaces and their containerization.
type Namespace interface {

	// Priority returns the priority of the namespace, which is used for sorting
	// namespaces in the From function. Namespaces with higher priorities are
	// processed earlier.
	Priority() int

	// Type returns the OCI LinuxNamespaceType key for the namespace
	Type() specs.LinuxNamespaceType

	// CloneFlag returns the flag that should be used with the `clone` or `unshare`
	// syscall. This flag determines the specific namespace that will be created
	// or unshared during these syscalls.
	CloneFlag() uintptr

	// Finalize performs additional steps required to complete the containerization
	// process after the `clone` or `unshare` syscall has been executed.
	Finalize(spec specs.Spec) error

	// IsSupported checks whether the namespace is supported by the runtime.
	// It returns `false` if the namespace cannot be used.
	IsSupported() bool
}

// From prepares a list of namespaces based on the provided spec.
// It takes a RuntimePipeWriter and a spec as inputs and returns a slice of Namespace and an error if any occurs.
// The namespaces are sorted based on their priority before being returned.
func From(runtime ipc.RuntimePipeWriter, spec specs.Spec) (namespaces []Namespace, err error) {
	log.Debug("Preparing namespaces from spec")

	namespacesSpec := spec.Linux.Namespaces
	namespaces = make([]Namespace, len(namespacesSpec))
	for i, namespace := range namespacesSpec {
		namespaces[i] = fromNamespaceType(runtime, spec, namespace.Type)
		log.Debug("prepared namespace", zap.Any("type", namespace.Type))
	}

	// Sort the namespaces based on their priority.
	slices.SortFunc(namespaces, func(a, b Namespace) int {
		return b.Priority() - b.Priority()
	})

	return namespaces, nil
}

// Unshare detaches the provided namespace from its parent by using the syscall.Unshare function.
func Unshare(namespace Namespace) error {
	return syscall.Unshare(int(namespace.CloneFlag()))
}

// fromNamespaceType maps a LinuxNamespaceType to the corresponding Namespace object.
// It returns the appropriate Namespace implementation based on the type, or nil if the type is unrecognized.
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
