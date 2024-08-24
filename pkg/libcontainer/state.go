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
	state    *specs.State
	mu       sync.Mutex
	stateDir string
}

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

func (s *StateManager) State() specs.State {
	return *s.state
}

func (s *StateManager) SetPid(pid int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.Pid = pid
}

func (s *StateManager) SetBundle(bundle string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.Bundle = bundle
}

func (s *StateManager) SetStatus(state specs.ContainerState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.Status = state
}

func (s *StateManager) UpdateState() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return util.WriteJsonFile(path.Join(s.stateDir, model.OciStateFileName), s.state)
}

func (s *StateManager) LoadState() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	err = util.ReadJsonFile(path.Join(s.stateDir, model.OciStateFileName), s.state)
	switch {
	case os.IsNotExist(err):
		return model.ErrNotExist
	case err != nil:
		return err
	default:
		return nil
	}
}
