package service

import (
	"github.com/craftdome/nymgraph/internal/entity"
	"github.com/craftdome/nymgraph/internal/repository"
	"github.com/craftdome/nymgraph/internal/repository/sqlite/pseudonyms"
	"github.com/craftdome/nymgraph/internal/state"
	"github.com/pkg/errors"
)

type PseudonymsService struct {
	repo  *repository.Repository
	state *state.State
}

func NewPseudonymsController(repo *repository.Repository, state *state.State) *PseudonymsService {
	return &PseudonymsService{
		repo:  repo,
		state: state,
	}
}

func (s *PseudonymsService) Create(name, server string) (*entity.Pseudonym, error) {
	dto := pseudonyms.CreateDTO{
		Name:   name,
		Server: server,
	}
	created, err := s.repo.Pseudonyms.Create(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Pseudonyms.Create %+v", dto)
	}
	return created, nil
}

func (s *PseudonymsService) Update(id int, name, server string) (*entity.Pseudonym, error) {
	dto := pseudonyms.UpdateDTO{
		ID:     id,
		Name:   name,
		Server: server,
	}
	updated, err := s.repo.Pseudonyms.Update(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Pseudonyms.Update %+v", dto)
	}
	return updated, nil
}

func (s *PseudonymsService) Delete(id int) (*entity.Pseudonym, error) {
	dto := pseudonyms.DeleteDTO{
		ID: id,
	}
	updated, err := s.repo.Pseudonyms.Delete(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Pseudonyms.Delete %+v", dto)
	}
	return updated, nil
}

func (s *PseudonymsService) GetAll() ([]*entity.Pseudonym, error) {
	all, err := s.repo.Pseudonyms.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "repo.Pseudonyms.GetAll")
	}

	return all, nil
}
