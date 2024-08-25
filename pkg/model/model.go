package model

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"
)

// SyscallSignal converts a signal name (either numeric or symbolic) into a syscall.Signal type.
// If the signalName is numeric, it converts directly.
// If the signalName is symbolic (e.g., "INT" or "SIGINT"), it returns the corresponding signal.
// It returns an error if the signal name is unknown.
func SyscallSignal(signalName string) (syscall.Signal, error) {
	s, err := strconv.Atoi(signalName)
	if err == nil {
		return syscall.Signal(s), nil
	}

	if !strings.HasPrefix(signalName, "SIG") {
		signalName = "SIG" + signalName
	}

	signal, exists := signalMap[signalName]
	if !exists {
		return 0, fmt.Errorf("unknown signal name")
	}

	return signal, nil
}

var signalMap = map[string]syscall.Signal{
	"SIGHUP":    1,
	"SIGINT":    2,
	"SIGQUIT":   3,
	"SIGILL":    4,
	"SIGTRAP":   5,
	"SIGABRT":   6,
	"SIGBUS":    7,
	"SIGFPE":    8,
	"SIGKILL":   9,
	"SIGUSR1":   10,
	"SIGSEGV":   11,
	"SIGUSR2":   12,
	"SIGPIPE":   13,
	"SIGALRM":   14,
	"SIGTERM":   15,
	"SIGSTKFLT": 16,
	"SIGCHLD":   17,
	"SIGCONT":   18,
	"SIGSTOP":   19,
	"SIGTSTP":   20,
	"SIGTTIN":   21,
	"SIGTTOU":   22,
	"SIGURG":    23,
	"SIGXCPU":   24,
	"SIGXFSZ":   25,
	"SIGVTALRM": 26,
	"SIGPROF":   27,
	"SIGWINCH":  28,
	"SIGIO":     29,
	"SIGPWR":    30,
	"SIGSYS":    31,
}
