package controller

import (
	"github.com/Tyz3/go-nym"
	"github.com/Tyz3/go-nym/response"
	"github.com/Tyz3/nymgraph/internal/entity"
	"github.com/Tyz3/nymgraph/internal/repository"
	"github.com/Tyz3/nymgraph/internal/state"
)

type NymClient interface {
	Dial(*entity.Pseudonym) error
	ListenAndServe() error
	Close() error
	SetErrorHandler(f func(*response.Error))
	SetReceivedHandler(f func(*response.Received))
	SetSelfAddressHandler(f func(*response.SelfAddress))
	SetCloseChannelHandler(f func())
	SendRequestAsText(r nym.Request) error
	SendRequestAsBinary(r nym.Request) error
}

type Contacts interface {
	Create(pseudonymID int, address, alias string) (*entity.Contact, error)
	Update(id int, address, alias string) (*entity.Contact, error)
	Delete(id int) (*entity.Contact, error)
	GetAll() ([]*entity.Contact, error)
}

type Pseudonyms interface {
	Create(name, server string) (*entity.Pseudonym, error)
	Update(id int, name, server string) (*entity.Pseudonym, error)
	Delete(id int) (*entity.Pseudonym, error)
	GetAll() ([]*entity.Pseudonym, error)
}

type IncomingMessages interface {
}

type OutcomingMessages interface {
}

type Controller struct {
	NymClient
	Contacts
	Pseudonyms
}

func NewController(repo *repository.Repository, state *state.State) *Controller {
	return &Controller{
		NymClient:  NewNymClientController(repo, state),
		Pseudonyms: NewPseudonymsController(repo, state),
	}
}
