package oci

import (
	"fmt"
	"github.com/opencontainers/runtime-spec/specs-go"
)

const (
	reducedOciVersionTag = "+modified"
)

var Version = fmt.Sprintf("%v%v", specs.Version, reducedOciVersionTag)

func Rootfs(spec *specs.Root) string {
	if spec == nil {
		return ""
	}
	return spec.Path
}
