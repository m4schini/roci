package rootfs

import (
	"errors"
	"github.com/opencontainers/runtime-spec/specs-go"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"roci/pkg/libcontainer/oci"
	"roci/pkg/logger"
	"syscall"
)

var SpecDevSymlinks = []DevSymlink{
	{Source: "/proc/self/fd", Target: "/dev/fd"},
	{Source: "/proc/self/fd/0", Target: "/dev/stdin"},
	{Source: "/proc/self/fd/1", Target: "/dev/stdout"},
	{Source: "/proc/self/fd/2", Target: "/dev/stderr"},
}

type DevSymlink struct {
	Source string
	Target string
}

func FinalizeRootfs(rootfs string, spec *specs.Spec) (err error) {
	var (
		log              = logger.Log().With(zap.String("rootfs", rootfs))
		mounts           = spec.Mounts
		setupDevRequired = checkSetupDevRequired(mounts)
	)

	err = syscall.Chdir(rootfs)
	if err != nil {
		return err
	}

	for _, mount := range mounts {
		err = mountInRootfs(rootfs, mount)
		if err != nil {
			log.Warn("mount failed", zap.String("type", mount.Type), zap.String("dest", mount.Destination))
			continue
		}
	}

	err = syscall.Chroot(rootfs)
	if err != nil {
		return err
	}

	err = syscall.Chdir("/")
	if err != nil {
		return err
	}

	if setupDevRequired {
		log.Debug("setup dev is required")
		err = setupDev(rootfs, spec)
		if err != nil {
			log.Warn("setup dev failed")
			return err
		}
	}

	return nil
}

func CleanRootfs(rootfs string, spec *specs.Spec) (err error) {
	var (
		mounts = spec.Mounts
	)

	for _, mount := range mounts {
		err = unmountInRootfs(rootfs, mount.Destination, true)
		if err != nil {
			continue
		}
	}

	_ = unmountInRootfs(rootfs, "dev", false)

	return nil
}

func mountInRootfs(rootfs string, mount specs.Mount) (err error) {
	destination := filepath.Join(rootfs, mount.Destination)
	if err := os.MkdirAll(destination, 0o755); err != nil {
		return err
	}

	flags, data := oci.MountOptions(&mount)
	wd, _ := os.Getwd()
	logger.Log().Debug("mounting", zap.String("cwd", wd), zap.String("source", mount.Source), zap.String("dest", destination), zap.String("type", mount.Type))
	return syscall.Mount(mount.Source, destination, mount.Type, uintptr(flags), data)
}

func unmountInRootfs(rootfs, mountDestination string, force bool) (err error) {
	destination := filepath.Join(rootfs, mountDestination)
	_, err = os.Stat(destination)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
	}

	var flags int = syscall.MNT_DETACH | syscall.MNT_FORCE

	logger.Log().Debug("unmounting", zap.String("dest", destination), zap.Bool("force", force))
	err = syscall.Unmount(destination, flags)
	if err != nil {
		logger.Log().Warn("failed to unmount", zap.Error(err), zap.String("dest", destination), zap.Bool("force", force))
		return err
	}

	return nil
}

// checkSetupDevRequired returns true if /dev needs to be set up.
// TODO cite runc
func checkSetupDevRequired(mounts []specs.Mount) bool {
	for _, m := range mounts {
		//TODO utils.clean?

		if m.Type == "bind" && filepath.Clean(m.Destination) == "/dev" {
			return false
		}
	}
	return true
}

func setupDev(rootfs string, spec *specs.Spec) (err error) {
	devdir := filepath.Join(rootfs, "dev")
	err = os.MkdirAll(devdir, 0755)
	if err != nil {
		return err
	}

	err = createDevs(spec)
	if err != nil {
		return err
	}

	return createDevSymlinks(rootfs)
}

func createDevs(spec *specs.Spec) error {
	//TODO create devices
	logger.Log().Warn("create devs is not implemented")
	return nil
}

func createDevSymlinks(rootfs string) error {
	for _, link := range SpecDevSymlinks {
		var (
			source = link.Source
			target = filepath.Join(rootfs, link.Target)
		)
		if err := os.Symlink(source, target); err != nil && !os.IsExist(err) {
			return err
		}
	}

	return nil
}
