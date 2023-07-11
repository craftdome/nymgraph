package service

import (
	"github.com/Tyz3/nymgraph/internal/model"
	"github.com/Tyz3/nymgraph/internal/repository"
	"github.com/Tyz3/nymgraph/internal/repository/sqlite/pseudonyms"
	"github.com/Tyz3/nymgraph/internal/repository/sqlite/received"
	"github.com/Tyz3/nymgraph/internal/repository/sqlite/replies"
	"github.com/Tyz3/nymgraph/internal/state"
	"github.com/pkg/errors"
)

type RepliesService struct {
	repo  *repository.Repository
	state *state.State
}

func NewRepliesService(repo *repository.Repository, state *state.State) *RepliesService {
	return &RepliesService{
		repo:  repo,
		state: state,
	}
}

func (s *RepliesService) Create(receivedID int, text string) (*model.Reply, error) {
	dto := replies.CreateDTO{
		ReceivedID: receivedID,
		Text:       text,
	}

	created, err := s.repo.Replies.Create(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Replies.Create %+v", dto)
	}

	dto2 := received.GetDTO{ID: receivedID}
	receive, err := s.repo.Received.Get(dto2)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Received.Get %+v", dto2)
	}

	dto3 := pseudonyms.GetDTO{ID: receive.PseudonymID}
	pseudonym, err := s.repo.Pseudonyms.Get(dto3)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Pseudonyms.Get %+v", dto3)
	}

	return &model.Reply{
		Reply: created,
		Received: &model.Received{
			Received:  receive,
			Pseudonym: pseudonym,
		},
	}, nil
}

func (s *RepliesService) Delete(id int) (*model.Reply, error) {
	dto := replies.DeleteDTO{
		ID: id,
	}
	deleted, err := s.repo.Replies.Delete(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Replies.Delete %+v", dto)
	}

	dto2 := received.GetDTO{ID: deleted.ReceivedID}
	receive, err := s.repo.Received.Get(dto2)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Received.Get %+v", dto2)
	}

	dto3 := pseudonyms.GetDTO{ID: receive.PseudonymID}
	pseudonym, err := s.repo.Pseudonyms.Get(dto3)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Pseudonyms.Get %+v", dto3)
	}

	return &model.Reply{
		Reply: deleted,
		Received: &model.Received{
			Received:  receive,
			Pseudonym: pseudonym,
		},
	}, nil
}

func (s *RepliesService) GetAll(receivedID int) ([]*model.Reply, error) {
	dto := replies.GetAllDTO{ReceivedID: receivedID}
	all, err := s.repo.Replies.GetAll(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Replies.GetAll %+v", dto)
	}

	dto2 := received.GetDTO{ID: receivedID}
	receive, err := s.repo.Received.Get(dto2)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Received.Get %+v", dto2)
	}

	dto3 := pseudonyms.GetDTO{ID: receive.PseudonymID}
	pseudonym, err := s.repo.Pseudonyms.Get(dto3)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Pseudonyms.Get %+v", dto3)
	}

	models := make([]*model.Reply, 0, len(all))
	for _, e := range all {
		models = append(models, &model.Reply{
			Reply: e,
			Received: &model.Received{
				Received:  receive,
				Pseudonym: pseudonym,
			},
		})
	}

	return models, nil
}

func (s *RepliesService) Truncate() error {
	if err := s.repo.Received.Truncate(); err != nil {
		return errors.Wrap(err, "repo.Replies.Truncate")
	}
	return nil
}
