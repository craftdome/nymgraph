package service

import (
	"github.com/Tyz3/nymgraph/internal/entity"
	"github.com/Tyz3/nymgraph/internal/repository"
	"github.com/Tyz3/nymgraph/internal/repository/sqlite/contacts"
	"github.com/Tyz3/nymgraph/internal/state"
	"github.com/pkg/errors"
)

type ContactsService struct {
	repo  *repository.Repository
	state *state.State
}

func NewContactsController(repo *repository.Repository, state *state.State) *ContactsService {
	return &ContactsService{
		repo:  repo,
		state: state,
	}
}

func (s *ContactsService) Create(pseudonymID int, addr, alias string) (*entity.Contact, error) {
	dto := contacts.CreateDTO{
		PseudonymID: pseudonymID,
		Address:     addr,
		Alias:       alias,
	}
	created, err := s.repo.Contacts.Create(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Contacts.Create %+v", dto)
	}
	return created, nil
}

func (s *ContactsService) Update(id int, addr, alias string) (*entity.Contact, error) {
	dto := contacts.UpdateDTO{
		ID:      id,
		Address: addr,
		Alias:   alias,
	}
	updated, err := s.repo.Contacts.Update(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Contacts.Update %+v", dto)
	}
	return updated, nil
}

func (s *ContactsService) Delete(id int) (*entity.Contact, error) {
	dto := contacts.DeleteDTO{
		ID: id,
	}
	updated, err := s.repo.Contacts.Delete(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Contacts.Delete %+v", dto)
	}
	return updated, nil
}

func (s *ContactsService) GetAll(pseudonymID int) ([]*entity.Contact, error) {
	dto := contacts.GetAllDTO{PseudonymID: pseudonymID}
	all, err := s.repo.Contacts.GetAll(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Contacts.GetAll %+v", dto)
	}

	return all, nil
}
