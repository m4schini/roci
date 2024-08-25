package procfs

import (
	"fmt"
	"testing"
)

func TestPid_String_PidSelf(t *testing.T) {
	pid := PidSelf

	str := fmt.Sprintf("%v", pid)
	if str != "self" {
		t.FailNow()
	}
}
