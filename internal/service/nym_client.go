package service

import (
	"github.com/craftdome/nymgraph/internal/entity"
	"github.com/craftdome/nymgraph/internal/nym_client"
	"github.com/craftdome/nymgraph/internal/repository"
	"github.com/craftdome/nymgraph/internal/state"
)

type NymClientService struct {
	repo  *repository.Repository
	state *state.State
}

func NewNymClientController(repo *repository.Repository, state *state.State) *NymClientService {
	return &NymClientService{
		repo:  repo,
		state: state,
	}
}

func (s *NymClientService) New(pseudonym *entity.Pseudonym) *nym_client.ClientConnect {
	return nym_client.NewClientConnect(pseudonym.DSN())
}
