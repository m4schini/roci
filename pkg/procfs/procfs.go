package procfs

import (
	"syscall"
)

// defaultPath is the default base path to the procfs directory, typically located at "/proc".
const (
	defaultPath = "/proc"
)

// Root is an instance of FS initialized with the default procfs path.
var Root = &FS{procfsPath: defaultPath}

// FS represents a filesystem structure for interacting with procfs.
type FS struct {
	// procfsPath is the path to the procfs directory, used to access process-specific resources.
	procfsPath string
}

// IsProcessRunning checks if a process with the given PID is running.
func IsProcessRunning(pid int) bool {
	if pid == 0 {
		return false
	}

	// Send a signal 0 to the process to check if it's running
	// syscall.Kill returns nil if the process is running, or an error otherwise
	err := syscall.Kill(pid, 0)
	if err == nil {
		return true
	}

	// Check if the error is due to the process not existing
	if err == syscall.ESRCH {
		return false
	}

	// For other errors, it's safer to assume the process is not running
	return false
}

func WaitForProcessStop(pid int) error {
	var status syscall.WaitStatus

	for {
		// Wait for the process with the given PID to change state.
		_, err := syscall.Wait4(pid, &status, syscall.WSTOPPED, nil)
		if err != nil {
			return err
		}

		if !IsProcessRunning(pid) {
			return nil
		}
	}
}
