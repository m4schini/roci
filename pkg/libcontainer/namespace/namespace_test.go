package namespace

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"testing"
)

func TestFrom(t *testing.T) {
	spec := specs.Spec{Linux: &specs.Linux{Namespaces: []specs.LinuxNamespace{
		{Type: specs.PIDNamespace},
		{Type: specs.UserNamespace},
		{Type: specs.UTSNamespace},
		{Type: specs.NetworkNamespace},
		{Type: specs.MountNamespace},
	}}}

	ns, err := From(nil, spec)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for i, n := range ns {
		t.Log(i, n.Priority(), n.Type())
	}
}
