package service

import (
	"github.com/Tyz3/nymgraph/internal/model"
	"github.com/Tyz3/nymgraph/internal/repository"
	"github.com/Tyz3/nymgraph/internal/repository/sqlite/contacts"
	"github.com/Tyz3/nymgraph/internal/repository/sqlite/sent"
	"github.com/Tyz3/nymgraph/internal/state"
	"github.com/pkg/errors"
)

type SentService struct {
	repo  *repository.Repository
	state *state.State
}

func NewSentService(repo *repository.Repository, state *state.State) *SentService {
	return &SentService{
		repo:  repo,
		state: state,
	}
}

func (s *SentService) Create(contactID int, text string, replySurbs int) (*model.Sent, error) {
	dto := sent.CreateDTO{
		ContactID:  contactID,
		Text:       text,
		ReplySurbs: replySurbs,
	}
	created, err := s.repo.Sent.Create(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Sent.Create %+v", dto)
	}

	dto2 := contacts.GetDTO{ID: contactID}
	contact, err := s.repo.Contacts.Get(dto2)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Contacts.Get %+v", dto2)
	}

	return &model.Sent{
		Sent:    created,
		Contact: contact,
	}, nil
}

func (s *SentService) Delete(id int) (*model.Sent, error) {
	dto := sent.DeleteDTO{
		ID: id,
	}
	deleted, err := s.repo.Sent.Delete(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Sent.Delete %+v", dto)
	}

	dto2 := contacts.GetDTO{ID: deleted.ContactID}
	contact, err := s.repo.Contacts.Get(dto2)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Contacts.Get %+v", dto2)
	}

	return &model.Sent{
		Sent:    deleted,
		Contact: contact,
	}, nil
}

func (s *SentService) GetAll(contactID int) ([]*model.Sent, error) {
	dto := sent.GetAllDTO{ContactID: contactID}
	all, err := s.repo.Sent.GetAll(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Sent.GetAll %+v", dto)
	}

	dto2 := contacts.GetDTO{ID: contactID}
	contact, err := s.repo.Contacts.Get(dto2)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Contacts.Get %+v", dto2)
	}

	models := make([]*model.Sent, 0, len(all))
	for _, e := range all {
		models = append(models, &model.Sent{
			Sent:    e,
			Contact: contact,
		})
	}

	return models, nil
}

func (s *SentService) Truncate() error {
	if err := s.repo.Sent.Truncate(); err != nil {
		return errors.Wrap(err, "repo.Sent.Truncate")
	}
	return nil
}
