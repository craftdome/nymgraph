package state

import (
	"github.com/craftdome/nymgraph/cmd/app/config"
)

type State struct {
	cfg *config.Config
}

func NewState(cfg *config.Config) *State {
	return &State{
		cfg: cfg,
	}
}

func (s *State) GetConfig() *config.Config {
	return s.cfg
}
