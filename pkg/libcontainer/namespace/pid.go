package namespace

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"go.uber.org/zap"
	"os"
	"syscall"
)

type pidNS struct {
	logger *zap.Logger
	spec   specs.Spec
}

func newPidNamespace(spec specs.Spec) *pidNS {
	return &pidNS{spec: spec, logger: zap.L().Named("namespace").Named("pid")}
}

func (p *pidNS) IsSupported() bool {
	return true
}

func (p *pidNS) Priority() int {
	return 2
}

func (p *pidNS) Type() specs.LinuxNamespaceType {
	return specs.PIDNamespace
}

func (p *pidNS) CloneFlag() uintptr {
	p.logger.Debug("current pid", zap.Int("pid", os.Getpid()))
	return syscall.CLONE_NEWPID
}

func (p *pidNS) Finalize(spec specs.Spec) error {
	log := p.logger

	//log.Debug("mounting procfs")
	//if err := syscall.Mount("proc", "proc", "proc", 0, ""); err != nil {
	//	return err
	//}
	p.logger.Debug("current pid", zap.Int("pid", os.Getpid()))

	log.Debug("applied namespace")
	return nil
}
