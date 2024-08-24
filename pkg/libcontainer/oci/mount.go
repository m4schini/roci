package oci

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"strings"
	"syscall"
)

func MountOptions(mount *specs.Mount) (flags int, opts string) {
	return ParseMountOptions(mount.Options)
}

func ParseMountOptions(options []string) (flags int, opts string) {
	var data []string

	for _, opt := range options {
		switch opt {
		case "async":
			flags &= ^syscall.MS_SYNCHRONOUS
		case "atime":
			flags &= ^syscall.MS_NOATIME
		case "bind":
			flags |= syscall.MS_BIND
		case "defaults":
			// No need to set any flags for defaults
		case "dev":
			flags &= ^syscall.MS_NODEV
		case "diratime":
			flags &= ^syscall.MS_NODIRATIME
		case "dirsync":
			flags |= syscall.MS_DIRSYNC
		case "exec":
			flags &= ^syscall.MS_NOEXEC
		case "iversion":
			flags |= syscall.MS_I_VERSION
		case "loud":
			// No direct mapping; potentially logging behavior
		case "mand":
			flags |= syscall.MS_MANDLOCK
		case "noatime":
			flags |= syscall.MS_NOATIME
		case "nodev":
			flags |= syscall.MS_NODEV
		case "nodiratime":
			flags |= syscall.MS_NODIRATIME
		case "noexec":
			flags |= syscall.MS_NOEXEC
		case "noiversion":
			flags &= ^syscall.MS_I_VERSION
		case "nomand":
			flags &= ^syscall.MS_MANDLOCK
		case "nosuid":
			flags |= syscall.MS_NOSUID
		case "private":
			flags |= syscall.MS_PRIVATE
		case "rbind":
			flags |= syscall.MS_BIND | syscall.MS_REC
		case "relatime":
			flags |= syscall.MS_RELATIME
		case "remount":
			flags |= syscall.MS_REMOUNT
		case "ro":
			flags |= syscall.MS_RDONLY
		case "rprivate":
			flags |= syscall.MS_PRIVATE | syscall.MS_REC
		case "rshared":
			flags |= syscall.MS_SHARED | syscall.MS_REC
		case "rslave":
			flags |= syscall.MS_SLAVE | syscall.MS_REC
		case "runbindable":
			flags |= syscall.MS_UNBINDABLE | syscall.MS_REC
		case "rw":
			flags &= ^syscall.MS_RDONLY
		case "shared":
			flags |= syscall.MS_SHARED
		case "silent":
			flags |= syscall.MS_SILENT
		case "slave":
			flags |= syscall.MS_SLAVE
		case "strictatime":
			flags |= syscall.MS_STRICTATIME
		case "suid":
			flags &= ^syscall.MS_NOSUID
		case "sync":
			flags |= syscall.MS_SYNCHRONOUS
		case "tmpcopyup":
			data = append(data, "tmpcopyup")
		case "unbindable":
			flags |= syscall.MS_UNBINDABLE
		default:
			data = append(data, opt)
		}
	}

	return flags, strings.Join(data, ",")
}
