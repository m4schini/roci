package libcontainer

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"path"
	"path/filepath"
	"roci/pkg/libcontainer/initp"
	"roci/pkg/model"
	"roci/pkg/util"
)

type Container struct {
	id     string
	state  *StateManager
	config specs.Spec
	initp  *initp.Process
}

func (c *Container) Init() (pid int, err error) {
	return c.initp.Start()
}

func (c *Container) State() specs.State {
	return c.state.State()
}

func CreateContainer(fs *FS, id, bundle string) (c *Container, err error) {
	var spec specs.Spec
	if err = util.ReadJsonFile(path.Join(bundle, model.OciSpecFileName), &spec); err != nil {
		return nil, err
	}

	PrepareSpec(&spec, bundle)

	c, err = fs.Create(id, bundle, spec)
	if err != nil {
		return nil, err
	}

	pid, err := c.Init()
	if err != nil {
		return nil, err
	}

	c.state.SetPid(pid)
	c.state.SetStatus(specs.StateCreated)
	if err = c.state.UpdateState(); err != nil {
		return nil, err
	}

	return c, nil
}

func PrepareSpec(spec *specs.Spec, bundle string) {
	var rootfsPath = spec.Root.Path
	if !filepath.IsAbs(rootfsPath) {
		rootfsPath = filepath.Join(bundle, rootfsPath)
	}
	spec.Root.Path = rootfsPath
}
