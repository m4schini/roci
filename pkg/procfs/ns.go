package procfs

import (
	"fmt"
	"github.com/opencontainers/runtime-spec/specs-go"
	"path/filepath"
	"syscall"
)

// Setns changes the namespace of the current process to the namespace of the specified process ID (pid).
func (F *FS) Setns(pid Pid, namespaceType specs.LinuxNamespaceType) error {
	// Open the namespace fd in the procfs directory
	nsFdPath := filepath.Join(F.procfsPath, pid.String(), "ns", string(namespaceType))
	fd, err := syscall.Open(nsFdPath, syscall.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	// Perform the setns syscall to change the current process's namespace to the target namespace.
	// 308 is the syscall number for setns on Linux.
	setnserr, _, _ := syscall.RawSyscall(308, uintptr(fd), 0, 0)
	if setnserr != 0 {
		return fmt.Errorf("%v", setnserr)
	}

	return nil
}
