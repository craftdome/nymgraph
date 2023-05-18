package controller

import (
	"github.com/Tyz3/nymgraph/internal/entity"
	"github.com/Tyz3/nymgraph/internal/repository"
	"github.com/Tyz3/nymgraph/internal/repository/sqlite/pseudonyms"
	"github.com/Tyz3/nymgraph/internal/state"
	"github.com/pkg/errors"
)

type PseudonymsController struct {
	repo  *repository.Repository
	state *state.State
}

func NewPseudonymsController(repo *repository.Repository, state *state.State) Pseudonyms {
	return &PseudonymsController{
		repo:  repo,
		state: state,
	}
}

func (c *PseudonymsController) Create(name, server string) (*entity.Pseudonym, error) {
	dto := pseudonyms.CreateDTO{
		Name:   name,
		Server: server,
	}
	created, err := c.repo.Pseudonyms.Create(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Pseudonyms.Create %+v", dto)
	}
	return created, nil
}

func (c *PseudonymsController) Update(id int, name, server string) (*entity.Pseudonym, error) {
	dto := pseudonyms.UpdateDTO{
		ID:     id,
		Name:   name,
		Server: server,
	}
	updated, err := c.repo.Pseudonyms.Update(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Pseudonyms.Update %+v", dto)
	}
	return updated, nil
}

func (c *PseudonymsController) Delete(id int) (*entity.Pseudonym, error) {
	dto := pseudonyms.DeleteDTO{
		ID: id,
	}
	updated, err := c.repo.Pseudonyms.Delete(dto)
	if err != nil {
		return nil, errors.Wrapf(err, "repo.Pseudonyms.Delete %+v", dto)
	}
	return updated, nil
}

func (c *PseudonymsController) GetAll() ([]*entity.Pseudonym, error) {
	pseudonyms, err := c.repo.Pseudonyms.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "repo.Pseudonyms.GetAll")
	}

	return pseudonyms, nil
}
