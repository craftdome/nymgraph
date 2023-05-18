package state

import (
	"github.com/Tyz3/nymgraph/cmd/app/config"
	"github.com/Tyz3/nymgraph/internal/entity"
)

type State struct {
	cfg *config.Config

	// RAM fields
	SelfAddress    string
	SelectedClient *entity.Pseudonym
}

func NewState(cfg *config.Config) *State {
	return &State{
		cfg: cfg,
	}
}
