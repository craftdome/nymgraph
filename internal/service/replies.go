package service

import (
	"github.com/Tyz3/nymgraph/internal/entity"
	"github.com/Tyz3/nymgraph/internal/repository"
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

func (s *RepliesService) Create(receivedID int, text string) (*entity.Reply, error) {
	dto := replies.CreateDTO{
		ReceivedID: receivedID,
		Text:       text,
	}

	created, err := s.repo.Replies.Create(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Replies.Create %+v", dto)
	}
	return created, nil
}

func (s *RepliesService) Delete(id int) (*entity.Reply, error) {
	dto := replies.DeleteDTO{
		ID: id,
	}
	updated, err := s.repo.Replies.Delete(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Replies.Delete %+v", dto)
	}
	return updated, nil
}

func (s *RepliesService) GetAll(receivedID int) ([]*entity.Reply, error) {
	dto := replies.GetAllDTO{ReceivedID: receivedID}
	all, err := s.repo.Replies.GetAll(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Replies.GetAll %+v", dto)
	}

	return all, nil
}

func (s *RepliesService) Truncate() error {
	if err := s.repo.Received.Truncate(); err != nil {
		return errors.Wrap(err, "repo.Replies.Truncate")
	}
	return nil
}
