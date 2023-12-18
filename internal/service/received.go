package service

import (
	"database/sql"
	"github.com/craftdome/nymgraph/internal/model"
	"github.com/craftdome/nymgraph/internal/repository"
	"github.com/craftdome/nymgraph/internal/repository/sqlite/pseudonyms"
	"github.com/craftdome/nymgraph/internal/repository/sqlite/received"
	"github.com/craftdome/nymgraph/internal/state"
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

func (s *ReceivedService) Create(pseudonymID int, text, senderTag string) (*model.Received, error) {
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

	dto2 := pseudonyms.GetDTO{ID: pseudonymID}
	pseudonym, err := s.repo.Pseudonyms.Get(dto2)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Pseudonyms.Get %+v", dto2)
	}

	return &model.Received{
		Received:  created,
		Pseudonym: pseudonym,
	}, nil
}

func (s *ReceivedService) Delete(id int) (*model.Received, error) {
	dto := received.DeleteDTO{
		ID: id,
	}
	deleted, err := s.repo.Received.Delete(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Received.Delete %+v", dto)
	}

	dto2 := pseudonyms.GetDTO{ID: deleted.PseudonymID}
	pseudonym, err := s.repo.Pseudonyms.Get(dto2)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Pseudonyms.Get %+v", dto2)
	}

	return &model.Received{
		Received:  deleted,
		Pseudonym: pseudonym,
	}, nil
}

func (s *ReceivedService) GetAll(pseudonymID int) ([]*model.Received, error) {
	dto := received.GetAllDTO{PseudonymID: pseudonymID}
	all, err := s.repo.Received.GetAll(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Received.GetAll %+v", dto)
	}

	dto2 := pseudonyms.GetDTO{ID: pseudonymID}
	pseudonym, err := s.repo.Pseudonyms.Get(dto2)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Pseudonyms.Get %+v", dto2)
	}

	models := make([]*model.Received, 0, len(all))
	for _, e := range all {
		models = append(models, &model.Received{
			Received:  e,
			Pseudonym: pseudonym,
		})
	}

	return models, nil
}

func (s *ReceivedService) Truncate() error {
	if err := s.repo.Received.Truncate(); err != nil {
		return errors.Wrap(err, "repo.Received.Truncate")
	}
	return nil
}
