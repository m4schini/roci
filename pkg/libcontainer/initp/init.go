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

//func InitFromStateDir(stateDir string) (err error) {
//	log := zap.L().Named("init")
//	configPath := path.Join(stateDir, model.OciSpecFileName)
//	var spec specs.Spec
//	log.Debug("reading config", zap.String("path", configPath))
//	err = util.ReadJsonFile(configPath, &spec)
//	if err != nil {
//		log.Warn("failed to read spec", zap.Error(err), zap.String("configPath", configPath))
//		return err
//	}
//	log.Debug("read config", zap.Any("config", spec))
//	return Init(stateDir, spec)
//}
//
//func Init(stateDir string, spec specs.Spec) (err error) {
//	var (
//		log     = zap.L().Named("init")
//		process = spec.Process
//		root    = spec.Root
//		state   = state{dir: stateDir}
//	)
//
//	if err = initDevSymLinks(*root); err != nil {
//		return err
//	}
//
//	wait, err := initWaitSocket(stateDir)
//	if err != nil {
//		return err
//	}
//
//	log.Debug("preparing namespaces")
//	namespaces, err := namespace.From(spec)
//	if err != nil {
//		return err
//	}
//	log.Debug("prepared namespaces", zap.Int("count", len(namespaces)))
//	err = namespace.FinalizeAll(namespaces...)
//	if err != nil {
//		return err
//	}
//	log.Debug("applied namespaces")
//
//	log.Debug("init hostname", zap.Any("hostname", spec.Hostname))
//	if err = initHostname(spec.Hostname); err != nil {
//		log.Debug("failed to init hostname", zap.Error(err))
//		return err
//	}
//
//	log.Debug("init domainame", zap.Any("domainname", spec.Domainname))
//	if err = initDomainname(spec.Domainname); err != nil {
//		log.Debug("failed to init domainname", zap.Error(err))
//		return err
//	}
//
//	log.Debug("init cwd", zap.Any("wd", process.Cwd))
//	if err = initCWD(process.Cwd); err != nil {
//		log.Debug("failed to init cwd", zap.Error(err))
//		return err
//	}
//
//	// chroot needs to happen before namespace because it takes
//	//if err = initRoot(root.Path, root.Readonly); err != nil {
//	//	log.Debug("failed to init root", zap.Error(err), zap.String("path", root.Path), zap.Bool("readonly", root.Readonly))
//	//	return err
//	//}
//
//	log.Debug("init exec", zap.Any("process", process))
//	if err = initExec(process); err != nil {
//		log.Debug("failed to init exec", zap.Error(err))
//		return err
//	}
//
//	if err = state.LoadState(); err != nil {
//		return err
//	}
//
//	state.state = state.WithStatus(specs.StateCreated)
//	if err = state.UpdateState(); err != nil {
//		return err
//	}
//
//	<-wait
//	log.Debug("received start signal")
//
//	defer func() {
//		state.state = state.WithStatus(specs.StateRunning)
//		state.UpdateState()
//	}()
//
//	log.Debug("executing container process")
//	return doexec(process.Args, os.Environ())
//}
//
//// Before namespaces and chroot
//func initDevSymLinks(root specs.Root) (err error) {
//	sl := SpecDevSymlinks.WithTargetBase(root.Path)
//	for _, symlink := range sl {
//		if err = os.Symlink(symlink.Source, symlink.Target); err != nil {
//			//TODO os exist check?
//			return err
//		}
//	}
//	return nil
//}

func initExec(process *specs.Process) (err error) {
	for _, s := range process.Env {
		kv := strings.Split(s, "=")
		err = os.Setenv(kv[0], kv[1])
		if err != nil {
			return err
		}
	}

	name, err := exec.LookPath(process.Args[0])
	if err != nil {
		return err
	}
	process.Args[0] = name
	return nil
}
