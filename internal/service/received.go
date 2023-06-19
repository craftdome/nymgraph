package service

import (
	"database/sql"
	"github.com/Tyz3/nymgraph/internal/entity"
	"github.com/Tyz3/nymgraph/internal/repository"
	"github.com/Tyz3/nymgraph/internal/repository/sqlite/received"
	"github.com/Tyz3/nymgraph/internal/state"
	"github.com/pkg/errors"
)

type ReceivedService struct {
	repo  *repository.Repository
	state *state.State
}

func NewReceivedService(repo *repository.Repository, state *state.State) *ReceivedService {
	return &ReceivedService{
		repo:  repo,
		state: state,
	}
}

func (s *ReceivedService) Create(pseudonymID int, text, senderTag string) (*entity.Received, error) {
	dto := received.CreateDTO{
		PseudonymID: pseudonymID,
		Text:        text,
	}

	if senderTag != "" {
		dto.SenderTag = sql.NullString{
			String: senderTag,
			Valid:  true,
		}
	}

	created, err := s.repo.Received.Create(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Received.Create %+v", dto)
	}
	return created, nil
}

func (s *ReceivedService) Delete(id int) (*entity.Received, error) {
	dto := received.DeleteDTO{
		ID: id,
	}
	updated, err := s.repo.Received.Delete(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Received.Delete %+v", dto)
	}
	return updated, nil
}

func (s *ReceivedService) GetAll(pseudonymID int) ([]*entity.Received, error) {
	dto := received.GetAllDTO{PseudonymID: pseudonymID}
	all, err := s.repo.Received.GetAll(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Received.GetAll %+v", dto)
	}

	return all, nil
}

func (s *ReceivedService) Truncate() error {
	if err := s.repo.Received.Truncate(); err != nil {
		return errors.Wrap(err, "repo.Received.Truncate")
	}
	return nil
}
