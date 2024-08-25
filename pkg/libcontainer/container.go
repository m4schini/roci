package libcontainer

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"path"
	"path/filepath"
	"roci/pkg/libcontainer/initp"
	"roci/pkg/libcontainer/oci"
	"roci/pkg/model"
	"roci/pkg/util"
)

// Container represents a container with its associated state and configuration.
type Container struct {
	id     string
	state  *StateManager
	config specs.Spec
	initp  *initp.Process
}

// Init starts the init process.
// It returns the process ID (pid) of the init process.
func (c *Container) Init() (pid int, err error) {
	return c.initp.Start()
}

// State returns the current state of the container.
func (c *Container) State() specs.State {
	return c.state.State()
}

// CreateContainer creates a new container using the container filesystem, id, and bundle path.
// It reads the container's specification, prepares it, creates the container, initializes it, and updates its state.
// Returns the created container
func CreateContainer(fs *FS, id, bundle string) (c *Container, err error) {
	var spec specs.Spec
	if err = util.ReadJsonFile(path.Join(bundle, model.OciSpecFileName), &spec); err != nil {
		return nil, err
	}

	PrepareSpec(&spec, bundle)

	// Create the container using the confs, ID, bundle, and prepared specification
	c, err = fs.Create(id, bundle, spec)
	if err != nil {
		return nil, err
	}

	// Initialize the container and retrieve its process ID
	pid, err := c.Init()
	if err != nil {
		return nil, err
	}

	// Set the container's process ID and update its state to "Created"
	c.state.SetPid(pid)
	c.state.SetStatus(specs.StateCreated)
	if err = c.state.UpdateState(); err != nil {
		return nil, err
	}

	err = oci.InvokeHooks(spec.Hooks, oci.HookCreateRuntime)
	if err != nil {
		return c, err
	}

	err = oci.InvokeHooks(spec.Hooks, oci.HookCreateContainer)
	if err != nil {
		return c, err
	}

	return c, nil
}

// PrepareSpec ensures that the root filesystem path in the container specification is absolute.
// If the path is relative, it will be joined with the provided bundle path.
func PrepareSpec(spec *specs.Spec, bundle string) {
	var rootfsPath = spec.Root.Path
	if !filepath.IsAbs(rootfsPath) {
		rootfsPath = filepath.Join(bundle, rootfsPath)
	}
	spec.Root.Path = rootfsPath
}
