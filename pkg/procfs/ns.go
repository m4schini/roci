package procfs

import (
	"fmt"
	"github.com/opencontainers/runtime-spec/specs-go"
	"path/filepath"
	"syscall"
)

func (F *FS) Setns(pid Pid, namespaceType specs.LinuxNamespaceType) error {
	nsFdPath := filepath.Join(F.BasePath, pid.String(), "ns", string(namespaceType))
	fd, err := syscall.Open(nsFdPath, syscall.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	// 308 == setns
	setnserr, _, _ := syscall.RawSyscall(308, uintptr(fd), 0, 0)
	if setnserr != 0 {
		return fmt.Errorf("%v", setnserr)
	}

	return nil
}
