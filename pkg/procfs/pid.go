package procfs

import (
	"fmt"
	"math"
)

type Pid int

const (
	PidSelf Pid = math.MinInt
)

func (p Pid) String() string {
	return pid(p)
}

func pid(pid Pid) string {
	switch pid {
	case PidSelf:
		return "self"
	default:
		return fmt.Sprintf("%d", pid)
	}
}
