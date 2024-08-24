package libcontainer

import (
	"fmt"
	"github.com/opencontainers/runtime-spec/specs-go"
	"go.uber.org/zap"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"roci/pkg/libcontainer/initp"
	"roci/pkg/libcontainer/ipc"
	"roci/pkg/libcontainer/oci"
	"roci/pkg/libcontainer/rootfs"
	"roci/pkg/logger"
	"roci/pkg/model"
	"roci/pkg/util"
	"syscall"
)

var (
	validIdRx = regexp.MustCompile(`[0-9A-Za-z_+\-.]+`)
)

type ContainerFS interface {
	/* Standard container operations */
	Create(id, bundle string, spec specs.Spec) (container *Container, err error)
	Start(id string) (err error)
	Kill(id string, signal syscall.Signal) (err error)
	Remove(id string) (err error)
	State(id string) (state specs.State, err error)

	/* Additional container operations */
	List() (containers []specs.State, err error)
}

type FS struct {
	dir string
}

func NewContainerFS(rootDir string) (*FS, error) {
	if rootDir == "" {
		return nil, fmt.Errorf("rootDir is empty")
	}
	if err := os.MkdirAll(rootDir, 0o700); err != nil {
		return nil, err
	}
	return &FS{dir: rootDir}, nil
}

func (r *FS) Create(id, bundle string, spec specs.Spec) (*Container, error) {
	stateDir, err := r.validateId(id)
	if err != nil {
		return nil, err
	}

	if err := os.Mkdir(stateDir, 0o711); err != nil {
		return nil, err
	}

	err = ipc.CreateRuntimePipe(stateDir)
	if err != nil {
		return nil, err
	}
	err = ipc.CreateInitPipe(stateDir)
	if err != nil {
		return nil, err
	}

	err = util.WriteJsonFile(path.Join(stateDir, model.OciSpecFileName), &spec)
	if err != nil {
		return nil, err
	}

	state, err := NewStateManager(stateDir, &specs.State{
		Version:     oci.Version,
		ID:          id,
		Status:      specs.StateCreating,
		Pid:         0,
		Bundle:      bundle,
		Annotations: nil,
	})
	if err != nil {
		return nil, err
	}

	c := &Container{
		id:     id,
		state:  state,
		config: spec,
		initp:  initp.NewInitProcess(spec.Root.Path, stateDir, spec.Hooks),
	}
	return c, nil
}

func (r *FS) Load(id string) (*Container, error) {
	if !validateIdFormat(id) {
		return nil, model.ErrInvalidID
	}
	state, err := LoadStateManager(r.stateDir(id))
	if err != nil {
		return nil, err
	}

	spec, err := r.loadSpec(id)
	if err != nil {
		return nil, err
	}

	c := &Container{
		id:     id,
		state:  state,
		config: *spec,
	}
	return c, nil
}

func (r *FS) Start(id string) (err error) {
	if err = r.assertContainerExists(id); err != nil {
		return err
	}
	log := logger.Log().Named("start")
	stateDir := r.stateDir(id)

	log.Debug("new init pipe writer")
	pipe, err := ipc.NewInitPipeWriter(stateDir)
	if err != nil {
		return err
	}

	log.Debug("new state manager")
	state, err := LoadStateManager(stateDir)
	if err != nil {
		return err
	}

	//// TODO this is probably slow
	//log.Debug("load spec")
	//spec, err := r.loadSpec(id)
	//if err != nil {
	//	return err
	//}

	//log.Debug("invoking hooks HookStartContainer")
	//err = oci.InvokeHooks(spec.Hooks, oci.HookStartContainer)
	//if err != nil {
	//	return err
	//}

	log.Debug("sending start")
	err = pipe.SendStart()
	if err != nil {
		return err
	}

	//log.Debug("invoking hooks HookPostStart")
	//err = oci.InvokeHooks(spec.Hooks, oci.HookPostStart)
	//if err != nil {
	//	return err
	//}

	state.SetStatus(specs.StateRunning)
	return state.UpdateState()
}

func (r *FS) Kill(id string, signal syscall.Signal) (err error) {
	state, err := r.State(id)
	if err != nil {
		return err
	}
	if state.Status != specs.StateRunning && state.Status != specs.StateCreated {
		return model.ErrNotRunning
	}

	err = syscall.Kill(state.Pid, signal)
	if err != nil {
		return err
	}

	return nil
}

func (r *FS) Remove(id string) (err error) {
	state, err := r.State(id)
	if err != nil {
		return err
	}
	if state.Status == specs.StateRunning {
		return model.ErrRunning
	}

	spec, err := r.loadSpec(id)
	if err != nil {
		return err
	}
	//defer func() {
	//	_ = oci.InvokeHooks(spec.Hooks, oci.HookPostStop)
	//}()

	err = rootfs.CleanRootfs(oci.Rootfs(*spec.Root), spec)
	if err != nil {
		return err
	}

	return os.RemoveAll(r.stateDir(id))
}

func (r *FS) State(id string) (state specs.State, err error) {
	if err = r.assertContainerExists(id); err != nil {
		return state, err
	}

	sm, err := LoadStateManager(r.stateDir(id))
	if err != nil {
		return state, err
	}
	state = sm.State()

	if state.Status == specs.StateRunning && !IsProcessRunning(state.Pid) {
		sm.SetStatus(specs.StateStopped)
		err = sm.UpdateState()
		if err != nil {
			return state, err
		}
	}

	return state, nil
}

func (r *FS) List() (containers []specs.State, err error) {
	files, err := os.ReadDir(r.dir)
	if err != nil {
		return nil, err
	}

	containers = make([]specs.State, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			containerId := filepath.Base(file.Name())
			state, err := r.State(containerId)
			if err != nil {
				continue
			}
			containers = append(containers, state)
		}
	}
	return containers, nil
}

func (r *FS) loadSpec(id string) (spec *specs.Spec, err error) {
	spec = new(specs.Spec)
	p := path.Join(r.stateDir(id), model.OciSpecFileName)
	logger.Log().Debug("loading spec", zap.String("path", p))
	err = util.ReadJsonFile(p, spec)
	return spec, err
}

func (r *FS) stateDir(id string) string {
	return path.Join(r.dir, id)
}

func (r *FS) validateId(id string) (stateDir string, err error) {
	if !validateIdFormat(id) {
		return "", fmt.Errorf("invalid id format")
	}

	return r.stateDir(id), r.assertContainerNotExists(id)
}

func (r *FS) assertContainerExists(id string) error {
	_, err := os.Stat(r.stateDir(id))
	switch {
	case os.IsNotExist(err):
		return model.ErrNotExist
	case err != nil:
		return err
	default:
		return nil
	}
}

func (r *FS) assertContainerNotExists(id string) error {
	_, err := os.Stat(r.stateDir(id))
	switch {
	case os.IsNotExist(err):
		return nil
	case err != nil:
		return err
	default:
		return model.ErrExist
	}
}

func validateIdFormat(id string) bool {
	return validIdRx.MatchString(id)
}

// IsProcessRunning checks if a process with the given PID is running.
func IsProcessRunning(pid int) bool {
	// Send a signal 0 to the process to check if it's running
	// syscall.Kill returns nil if the process is running, or an error otherwise
	err := syscall.Kill(pid, 0)
	if err == nil {
		return true
	}

	// Check if the error is due to the process not existing
	if err == syscall.ESRCH {
		return false
	}

	// For other errors, it's safer to assume the process is not running
	return false
}
