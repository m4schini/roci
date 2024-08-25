package initp

import (
	"errors"
	"github.com/opencontainers/runtime-spec/specs-go"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"roci/pkg/libcontainer/ipc"
	"roci/pkg/libcontainer/namespace"
	"roci/pkg/libcontainer/rootfs"
	"roci/pkg/logger"
	"strings"
	"syscall"
)

func Init(stateDir string, spec specs.Spec) (err error) {
	var (
		log        = logger.Log()
		rootfsPath = spec.Root.Path
		wd, _      = os.Getwd()
	)
	log.Debug("change dir to rootfs", zap.String("rootfs", rootfsPath), zap.String("wd", wd))
	err = syscall.Chdir(rootfsPath)
	if err != nil {
		return err
	}
	log.Debug("init started", zap.Int("pid", os.Getpid()))

	log.Debug("create runtime.pipe client")
	runtime, err := ipc.NewRuntimePipeWriter(stateDir)
	if err != nil {
		return err
	}

	waitForStart := make(chan struct{})
	go func() {
		defer close(waitForStart)
		log.Debug("create init.pipe listener")
		pipe, err := ipc.NewInitPipeReader(stateDir)
		if err != nil {
			panic(err)
		}

		log.Debug("wait for start on pipe")
		err = pipe.WaitForStart()
		if err != nil {
			panic(err)
		}

		log.Debug("received start")
		waitForStart <- struct{}{}
	}()

	log.Debug("prepare namespaces")
	namespaces, err := namespace.From(runtime, spec)
	if err != nil {
		return err
	}

	log.Debug("create namespaces")
	for i, ns := range namespaces {
		if !ns.IsSupported() {
			log.Warn("namespace is not supported", zap.Any("type", ns.Type()))
			continue
		}
		if ns.Type() == specs.PIDNamespace {
			log.Debug("skipping pid namespace because it's already unshared (happened in runtime proc)")
			continue
		}

		log.Debug("init namespace", zap.Int("i", i), zap.Any("ns", ns.Type()))
		err = namespace.Unshare(ns)
		if err != nil {
			return err
		}

		log.Debug("finalizing namespace in rootfs", zap.Int("i", i), zap.Any("ns", ns.Type()))
		err = ns.Finalize(spec)
		if err != nil {
			return err
		}
	}

	log.Debug("prepare rootfs", zap.String("rootfs", rootfsPath))
	err = rootfs.FinalizeRootfs(rootfsPath, &spec)
	if err != nil {
		return err
	}

	log.Debug("notify runtime that container is ready")
	err = runtime.SendReady()
	if err != nil {
		return err
	}

	log.Debug("wait for runtime start signal")
	<-waitForStart

	log.Debug("exec container entrypoint")
	return execEntrypoint(spec)
}

func Entrypoint(process *specs.Process) (bin string, args, env []string, err error) {
	name, err := exec.LookPath(process.Args[0])
	if err != nil {
		return "", nil, nil, err
	}
	process.Args[0] = name
	return process.Args[0], process.Args, process.Env, nil
}

func execEntrypoint(spec specs.Spec) (err error) {
	arg0, args, env, err := Entrypoint(spec.Process)
	if err != nil {
		return err
	}

	for {
		err = syscall.Exec(arg0, args, env)
		if !errors.Is(err, syscall.EINTR) {
			return err
		}
	}
}
