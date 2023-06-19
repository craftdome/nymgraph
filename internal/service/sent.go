package service

import (
	"github.com/Tyz3/nymgraph/internal/entity"
	"github.com/Tyz3/nymgraph/internal/repository"
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

func (s *SentService) Create(contactID int, text string) (*entity.Sent, error) {
	dto := sent.CreateDTO{
		ContactID: contactID,
		Text:      text,
	}

	created, err := s.repo.Sent.Create(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Sent.Create %+v", dto)
	}
	return created, nil
}

func (s *SentService) Delete(id int) (*entity.Sent, error) {
	dto := sent.DeleteDTO{
		ID: id,
	}
	updated, err := s.repo.Sent.Delete(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Sent.Delete %+v", dto)
	}
	return updated, nil
}

func (s *SentService) GetAll(contactID int) ([]*entity.Sent, error) {
	dto := sent.GetAllDTO{ContactID: contactID}
	all, err := s.repo.Sent.GetAll(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Sent.GetAll %+v", dto)
	}

	return all, nil
}

func (s *SentService) Truncate() error {
	if err := s.repo.Sent.Truncate(); err != nil {
		return errors.Wrap(err, "repo.Sent.Truncate")
	}
	return nil
}
