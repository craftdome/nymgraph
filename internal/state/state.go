package state

import (
	"github.com/Tyz3/nymgraph/cmd/app/config"
)

type State struct {
	cfg *config.Config
}

func NewState(cfg *config.Config) *State {
	return &State{
		cfg: cfg,
	}
}
