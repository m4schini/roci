package libcontainer

import (
	"github.com/opencontainers/runtime-spec/specs-go"
	"os"
	"path"
	"roci/pkg/model"
	"roci/pkg/util"
	"sync"
)

type StateManager struct {
	// state is a pointer to the container's state as defined by the OCI runtime spec
	state *specs.State

	// stateDir is the directory where the state file is stored
	stateDir string

	mu sync.Mutex
}

// NewStateManager creates a new StateManager instance and updates the state file.
// It takes the state directory and an initial state as input.
// Returns a pointer to the StateManager and an error if updating the state fails.
func NewStateManager(stateDir string, state *specs.State) (*StateManager, error) {
	sm := &StateManager{
		state:    state,
		stateDir: stateDir,
	}
	err := sm.UpdateState()
	if err != nil {
		return nil, err
	}
	return sm, nil
}

// LoadStateManager loads an existing state from the state directory.
// It initializes a new StateManager and attempts to load the state from the state file.
// Returns a pointer to the StateManager and an error if loading the state fails.
func LoadStateManager(stateDir string) (*StateManager, error) {
	sm := &StateManager{
		state:    new(specs.State),
		stateDir: stateDir,
	}
	err := sm.LoadState()
	if err != nil {
		return nil, err
	}
	return sm, nil
}

// State returns a copy of the current container state.
func (s *StateManager) State() specs.State {
	return *s.state
}

// SetPid sets the PID of the container process
func (s *StateManager) SetPid(pid int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.Pid = pid
}

// SetBundle sets the bundle path of the container
func (s *StateManager) SetBundle(bundle string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.Bundle = bundle
}

// SetStatus sets the container status
func (s *StateManager) SetStatus(state specs.ContainerState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.Status = state
}

// UpdateState writes the current state to the state file.
// Returns an error if the file writing fails.
func (s *StateManager) UpdateState() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return util.WriteJsonFile(path.Join(s.stateDir, model.OciStateFileName), s.state)
}

// LoadState reads the state from the state file.
// Returns an error if the file reading fails, or if the file does not exist.
func (s *StateManager) LoadState() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	err = util.ReadJsonFile(path.Join(s.stateDir, model.OciStateFileName), s.state)
	switch {
	case os.IsNotExist(err):
		// if the state file is cannot be found it is assumed the container doesn't exist
		return model.ErrNotExist
	case err != nil:
		return err
	default:
		return nil
	}
}
