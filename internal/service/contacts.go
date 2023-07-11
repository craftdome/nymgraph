package service

import (
	"github.com/Tyz3/nymgraph/internal/model"
	"github.com/Tyz3/nymgraph/internal/repository"
	"github.com/Tyz3/nymgraph/internal/repository/sqlite/contacts"
	"github.com/Tyz3/nymgraph/internal/repository/sqlite/pseudonyms"
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

func (s *ContactsService) Create(pseudonymID int, addr, alias string) (*model.Contact, error) {
	dto := contacts.CreateDTO{
		PseudonymID: pseudonymID,
		Address:     addr,
		Alias:       alias,
	}
	created, err := s.repo.Contacts.Create(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Contacts.Create %+v", dto)
	}

	dto2 := pseudonyms.GetDTO{ID: pseudonymID}
	pseudonym, err := s.repo.Pseudonyms.Get(dto2)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Pseudonyms.Get %+v", dto2)
	}

	return &model.Contact{
		Contact:   created,
		Pseudonym: pseudonym,
	}, nil
}

func (s *ContactsService) Update(id int, addr, alias string) (*model.Contact, error) {
	dto := contacts.UpdateDTO{
		ID:      id,
		Address: addr,
		Alias:   alias,
	}
	updated, err := s.repo.Contacts.Update(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Contacts.Update %+v", dto)
	}

	dto2 := pseudonyms.GetDTO{ID: updated.PseudonymID}
	pseudonym, err := s.repo.Pseudonyms.Get(dto2)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Pseudonyms.Get %+v", dto2)
	}

	return &model.Contact{
		Contact:   updated,
		Pseudonym: pseudonym,
	}, nil
}

func (s *ContactsService) Delete(id int) (*model.Contact, error) {
	dto := contacts.DeleteDTO{
		ID: id,
	}
	deleted, err := s.repo.Contacts.Delete(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Contacts.Delete %+v", dto)
	}

	dto2 := pseudonyms.GetDTO{ID: deleted.PseudonymID}
	pseudonym, err := s.repo.Pseudonyms.Get(dto2)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Pseudonyms.Get %+v", dto2)
	}

	return &model.Contact{
		Contact:   deleted,
		Pseudonym: pseudonym,
	}, nil
}

func (s *ContactsService) GetAll(pseudonymID int) ([]*model.Contact, error) {
	dto := contacts.GetAllDTO{PseudonymID: pseudonymID}
	all, err := s.repo.Contacts.GetAll(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Contacts.GetAll %+v", dto)
	}

	dto2 := pseudonyms.GetDTO{ID: pseudonymID}
	pseudonym, err := s.repo.Pseudonyms.Get(dto2)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Pseudonyms.Get %+v", dto2)
	}

	models := make([]*model.Contact, 0, len(all))
	for _, e := range all {
		models = append(models, &model.Contact{
			Contact:   e,
			Pseudonym: pseudonym,
		})
	}

	return models, nil
}
