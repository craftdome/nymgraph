package service

import (
	"github.com/Tyz3/nymgraph/internal/entity"
	"github.com/Tyz3/nymgraph/internal/nym_client"
	"github.com/Tyz3/nymgraph/internal/repository"
	"github.com/Tyz3/nymgraph/internal/state"
)

type NymClient interface {
	New(pseudonym *entity.Pseudonym) *nym_client.ClientConnect
}

type Contacts interface {
	Create(pseudonymID int, addr, alias string) (*entity.Contact, error)
	Update(id int, addr, alias string) (*entity.Contact, error)
	Delete(id int) (*entity.Contact, error)
	GetAll(pseudonymID int) ([]*entity.Contact, error)
}

type Pseudonyms interface {
	Create(name, server string) (*entity.Pseudonym, error)
	Update(id int, name, server string) (*entity.Pseudonym, error)
	Delete(id int) (*entity.Pseudonym, error)
	GetAll() ([]*entity.Pseudonym, error)
}

type Sent interface {
	Create(contactID int, text string) (*entity.Sent, error)
	Delete(id int) (*entity.Sent, error)
	GetAll(contactID int) ([]*entity.Sent, error)
	Truncate() error
}

type Replies interface {
	Create(receivedID int, text string) (*entity.Reply, error)
	Delete(id int) (*entity.Reply, error)
	GetAll(receivedID int) ([]*entity.Reply, error)
	Truncate() error
}

type Received interface {
	Create(pseudonymID int, text, senderTag string) (*entity.Received, error)
	Delete(id int) (*entity.Received, error)
	GetAll(pseudonymID int) ([]*entity.Received, error)
	Truncate() error
}

type Service struct {
	NymClient
	Contacts
	Pseudonyms
	Sent
	Replies
	Received
}

func NewService(repo *repository.Repository, state *state.State) *Service {
	return &Service{
		NymClient:  NewNymClientController(repo, state),
		Pseudonyms: NewPseudonymsController(repo, state),
		Contacts:   NewContactsController(repo, state),
		Sent:       NewSentService(repo, state),
		Replies:    NewRepliesService(repo, state),
		Received:   NewReceivedService(repo, state),
	}
}
