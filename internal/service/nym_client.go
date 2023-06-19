package service

import (
	"github.com/Tyz3/nymgraph/internal/entity"
	"github.com/Tyz3/nymgraph/internal/nym_client"
	"github.com/Tyz3/nymgraph/internal/repository"
	"github.com/Tyz3/nymgraph/internal/state"
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
	return nym_client.NewClientConnect(s.repo, pseudonym.DSN())
}
