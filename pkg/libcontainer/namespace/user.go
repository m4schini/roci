package namespace

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"roci/pkg/procfs"
	"syscall"
)

type userNS struct {
	idMapper procfs.IdMapper
}

func newUserNamespace(idMapper procfs.IdMapper) *userNS {
	ns := new(userNS)
	ns.idMapper = idMapper
	return ns
}

func (u *userNS) IsSupported() bool {
	return true
}

func (u *userNS) Priority() int {
	return 3
}

func (u *userNS) Type() specs.LinuxNamespaceType {
	return specs.UserNamespace
}

func (u *userNS) CloneFlag() uintptr {
	return syscall.CLONE_NEWUSER
}

func (u *userNS) Finalize(spec specs.Spec) (err error) {
	insideUid := spec.Process.User.UID
	if err = u.idMapper.MapUid(procfs.PidSelf, insideUid, uint32(syscall.Getuid())); err != nil {
		return err
	}

	insideGid := spec.Process.User.GID
	if err = u.idMapper.MapGid(procfs.PidSelf, insideGid, uint32(syscall.Getgid())); err != nil {
		return err
	}

	if err = syscall.Setuid(int(insideUid)); err != nil {
		return err
	}

	if err = syscall.Setgid(int(insideGid)); err != nil {
		return err
	}

	return
}
