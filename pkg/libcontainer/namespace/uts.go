package namespace

import (
	"errors"
	"github.com/opencontainers/runtime-spec/specs-go"
	"syscall"
)

type utsNS struct {
	specs specs.Spec
}

func (u *utsNS) IsSupported() bool {
	return true
}

func (u *utsNS) Priority() int {
	return 0
}

func newUtsNamespace(spec specs.Spec) *utsNS {
	return &utsNS{specs: spec}
}

func (u *utsNS) Type() specs.LinuxNamespaceType {
	return specs.UTSNamespace
}

func (u *utsNS) CloneFlag() uintptr {
	return syscall.CLONE_NEWUTS
}

func (u *utsNS) Finalize(spec specs.Spec) error {
	return errors.Join(
		setHostname(spec.Hostname),
		setDomainname(spec.Domainname),
	)
}

func setHostname(hostname string) (err error) {
	if hostname == "" {
		return nil
	}
	return syscall.Sethostname([]byte(hostname))
}

func setDomainname(domainName string) (err error) {
	if domainName == "" {
		return nil
	}
	return syscall.Setdomainname([]byte(domainName))
}
