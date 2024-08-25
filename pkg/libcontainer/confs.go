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
	"roci/pkg/procfs"
	"roci/pkg/util"
	"syscall"
)

var (
	validIdRx = regexp.MustCompile(`[0-9A-Za-z_+\-.]+`)
)

// ContainerFS defines the interface for container file system operations.
type ContainerFS interface {
	/* Standard container operations */

	// Create initializes and creates a new container with the given ID, bundle path, and OCI runtime specification.
	// Returns the created container and any error encountered.
	Create(id, bundle string, spec specs.Spec) (container *Container, err error)

	// Start launches the container with the specified ID.
	// It returns any error encountered during the start process.
	Start(id string) (err error)

	// Kill sends a termination signal to the container with the specified ID.
	// It accepts a signal of type syscall.Signal and returns any error encountered.
	Kill(id string, signal syscall.Signal) (err error)

	// Remove deletes the container with the given ID.
	// It returns any error encountered during the removal process.
	Remove(id string) (err error)

	// State retrieves the current state of the container with the specified ID.
	// It returns the container's state and any error encountered.
	State(id string) (state specs.State, err error)

	/* Additional container operations */

	// List returns a list of all containers' states.
	// It returns the list of container states and any error encountered.
	List() (containers []specs.State, err error)
}

// FS represents a file system that manages containers.
type FS struct {
	dir string // Directory where container (state and configuration) are stored
}

// NewContainerFS creates a new FS instance with the specified container directory.
// It initializes the directory structure if it does not already exist.
func NewContainerFS(rootDir string) (*FS, error) {
	if rootDir == "" {
		return nil, fmt.Errorf("rootDir is empty")
	}
	if err := os.MkdirAll(rootDir, 0o700); err != nil {
		return nil, err
	}
	return &FS{dir: rootDir}, nil
}

// Create initializes and creates a new container with the given ID, bundle path, and OCI runtime specification.
// It returns the created container and any error encountered.
func (r *FS) Create(id, bundle string, spec specs.Spec) (*Container, error) {
	stateDir, err := r.validateId(id)
	if err != nil {
		return nil, err
	}

	if err := os.Mkdir(stateDir, 0o711); err != nil {
		return nil, err
	}

	// Create IPC pipes required for communication with the container
	err = ipc.CreateRuntimePipe(stateDir)
	if err != nil {
		return nil, err
	}
	err = ipc.CreateInitPipe(stateDir)
	if err != nil {
		return nil, err
	}

	// Copy the OCI runtime specification into the state dir
	err = util.WriteJsonFile(path.Join(stateDir, model.OciSpecFileName), &spec)
	if err != nil {
		return nil, err
	}

	// Initialize the state
	state, err := NewStateManager(stateDir, &specs.State{
		Version: oci.Version,
		ID:      id,
		Status:  specs.StateCreating,
		Bundle:  bundle,
	})
	if err != nil {
		return nil, err
	}

	// Create and return a new Container instance
	c := &Container{
		id:     id,
		state:  state,
		config: spec,
		initp:  initp.NewInitProcess(spec.Root.Path, stateDir, &spec),
	}
	return c, nil
}

// Start launches the container specified by the given ID.
// Uses the ipc pipes to send the start signal
// It returns any error encountered during the start process.
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

	log.Debug("load spec")
	spec, err := r.loadSpec(id)
	if err != nil {
		return err
	}

	log.Debug("invoking hooks HookStartContainer")
	err = oci.InvokeHooks(spec.Hooks, oci.HookStartContainer)
	if err != nil {
		return err
	}

	// Send a start signal to the container
	log.Debug("sending start")
	err = pipe.SendStart()
	if err != nil {
		return err
	}

	log.Debug("invoking hooks HookPostStart")
	err = oci.InvokeHooks(spec.Hooks, oci.HookPostStart)
	if err != nil {
		return err
	}

	// Update the container's state to "Running"
	state.SetStatus(specs.StateRunning)
	return state.UpdateState()
}

// Kill sends a termination signal to the container with the specified ID.
// It returns any error encountered during the process.
func (r *FS) Kill(id string, signal syscall.Signal) (err error) {
	sm, err := LoadStateManager(r.stateDir(id))
	if err != nil {
		return err
	}
	state := sm.State()

	if state.Status != specs.StateRunning && state.Status != specs.StateCreated {
		return model.ErrNotRunning
	}

	err = syscall.Kill(state.Pid, signal)
	switch {
	case err == syscall.ESRCH:
		return nil
	case err != nil:
		return err
	}

	_ = procfs.WaitForProcessStop(state.Pid)
	sm.SetStatus(specs.StateStopped)
	_ = sm.UpdateState()

	return nil
}

// Remove deletes the container with the specified ID.
// It returns any error encountered during the removal process.
func (r *FS) Remove(id string) (err error) {
	state, err := r.State(id)
	if err != nil {
		return err
	}
	if state.Status != specs.StateStopped {
		return model.ErrRunning
	}

	spec, err := r.loadSpec(id)
	if err != nil {
		return err
	}

	// Clean up the container's root filesystem
	err = rootfs.CleanRootfs(oci.Rootfs(spec.Root), spec)
	if err != nil {
		return err
	}

	// Remove the container's state directory
	err = os.RemoveAll(r.stateDir(id))
	if err != nil {
		return err
	}

	err = oci.InvokeHooks(spec.Hooks, oci.HookPostStop)
	if err != nil {
		return err
	}

	return nil
}

// State retrieves the current state of the container with the specified ID.
// It returns the container's state and any error encountered.
func (r *FS) State(id string) (state specs.State, err error) {
	if err = r.assertContainerExists(id); err != nil {
		return state, err
	}

	sm, err := LoadStateManager(r.stateDir(id))
	if err != nil {
		return state, err
	}
	state = sm.State()
	if state.Status == specs.StateStopped || state.Status == specs.StateCreated {
		return state, nil
	}

	// If the container is running but its process is no longer active, update the state
	if !procfs.IsProcessRunning(state.Pid) {
		sm.SetStatus(specs.StateStopped)
		err = sm.UpdateState()
		if err != nil {
			return state, err
		}
	}

	return state, nil
}

// List returns a list of all containers' states.
// It returns the list of container states and any error encountered.
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

// loadSpec loads the OCI runtime specification for the container with the given ID.
// It returns the specification and any error encountered.
func (r *FS) loadSpec(id string) (spec *specs.Spec, err error) {
	spec = new(specs.Spec)
	p := path.Join(r.stateDir(id), model.OciSpecFileName)
	logger.Log().Debug("loading spec", zap.String("path", p))
	err = util.ReadJsonFile(p, spec)
	return spec, err
}

// stateDir returns the directory path for the container with the given ID.
func (r *FS) stateDir(id string) string {
	return path.Join(r.dir, id)
}

// validateId validates the format of the container ID and ensures it does not already exist.
// It returns the state directory for the container and any error encountered.
func (r *FS) validateId(id string) (stateDir string, err error) {
	if !validateIdFormat(id) {
		return "", fmt.Errorf("invalid id format")
	}

	return r.stateDir(id), r.assertContainerNotExists(id)
}

// assertContainerExists checks if a container with the given ID exists.
// It returns an error if the container does not exist or if any other error occurs.
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

// assertContainerNotExists checks if a container with the given ID does not exist.
// It returns an error if the container already exists or if any other error occurs.
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

// validateIdFormat checks if the container ID is in a valid format using a regular expression.
func validateIdFormat(id string) bool {
	return validIdRx.MatchString(id)
}
