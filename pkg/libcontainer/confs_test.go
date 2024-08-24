package libcontainer

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"testing"
)

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestRoot(t *testing.T) {
	rootDir := "/tmp/roci"
	root, err := NewContainerFS(rootDir)
	checkErr(t, err)

	c, err := root.Create("test", specs.Spec{
		Version: "test",
	})
	checkErr(t, err)

	t.Log(c.id)
	t.Log(c.stateDir)
	t.Log(c.config)
}
