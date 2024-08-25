package procfs

import (
	"fmt"
	"math"
)

// Pid represents a process ID.
type Pid int

const (
	// PidSelf is a helper constant to specify that "self" should be returned
	PidSelf Pid = math.MinInt
)

func (p Pid) String() string {
	return pid(p)
}

// pid converts a Pid to its string representation.
// It returns "self" for PidSelf and the numeric string representation for other PIDs.
func pid(pid Pid) string {
	switch pid {
	case PidSelf:
		return "self"
	default:
		return fmt.Sprintf("%d", pid)
	}
}
