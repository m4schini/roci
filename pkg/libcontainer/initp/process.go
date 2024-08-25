package initp

import (
	"context"
	"github.com/opencontainers/runtime-spec/specs-go"
	"os"
	"os/exec"
	"roci/pkg/libcontainer/ipc"
	"roci/pkg/libcontainer/oci"
	"roci/pkg/logger"
	"roci/pkg/procfs"
	"syscall"
)

type ProcessSpec = specs.Process

type Process struct {
	cmd      *exec.Cmd
	stateDir string
	hooks    *specs.Hooks
}

func NewInitProcess(rootfs, stateDir string, spec *specs.Spec) *Process {
	cmd, err := prepareCmd(stateDir)
	if err != nil {
		panic(err) //TODO
	}
	return &Process{
		stateDir: stateDir,
		cmd:      cmd,
		hooks:    spec.Hooks,
	}
}

func (i *Process) Start() (pid int, err error) {
	err = i.cmd.Start()
	if err != nil {
		return -1, err
	}

	waitForReady, pipe, err := ipc.NewRuntimePipeReader(context.Background(), i.stateDir, procfs.Root)
	if err != nil {
		return -1, err
	}
	defer pipe.Close()

	logger.Log().Debug("waiting for init process")
	select {
	case <-waitForReady:
		logger.Log().Debug("received ready")
		break
	case err := <-i.wait():
		logger.Log().Debug("init process stopped")
		if err != nil {
			return -1, err
		}
	}

	pid = i.cmd.Process.Pid
	err = oci.InvokeHooks(i.hooks, oci.HookCreateRuntime)
	if err != nil {
		return pid, err
	}

	err = oci.InvokeHooks(i.hooks, oci.HookCreateContainer)
	if err != nil {
		return pid, err
	}

	return pid, nil
}

func (i *Process) wait() <-chan error {
	ch := make(chan error, 1)
	go func() {
		defer close(ch)

		ch <- i.cmd.Wait()
	}()
	return ch
}

func prepareCmd(stateDir string) (*exec.Cmd, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(executablePath, "init", stateDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
		//pid namespace is unshared here because unshare can't move to current process to pid 1
		Cloneflags: syscall.CLONE_NEWPID,
	}
	return cmd, nil
}
