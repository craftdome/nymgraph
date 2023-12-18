package service

import (
	"github.com/craftdome/nymgraph/internal/entity"
	"github.com/craftdome/nymgraph/internal/model"
	"github.com/craftdome/nymgraph/internal/nym_client"
	"github.com/craftdome/nymgraph/internal/repository"
	"github.com/craftdome/nymgraph/internal/state"
)

type NymClient interface {
	New(pseudonym *entity.Pseudonym) *nym_client.ClientConnect
}

type Contacts interface {
	Create(pseudonymID int, addr, alias string) (*model.Contact, error)
	Update(id int, addr, alias string) (*model.Contact, error)
	Delete(id int) (*model.Contact, error)
	GetAll(pseudonymID int) ([]*model.Contact, error)
}

type Pseudonyms interface {
	Create(name, server string) (*entity.Pseudonym, error)
	Update(id int, name, server string) (*entity.Pseudonym, error)
	Delete(id int) (*entity.Pseudonym, error)
	GetAll() ([]*entity.Pseudonym, error)
}

type Sent interface {
	Create(contactID int, text string, replySurbs int) (*model.Sent, error)
	Delete(id int) (*model.Sent, error)
	GetAll(contactID int) ([]*model.Sent, error)
	Truncate() error
}

type Replies interface {
	Create(receivedID int, text string) (*model.Reply, error)
	Delete(id int) (*model.Reply, error)
	GetAll(receivedID int) ([]*model.Reply, error)
	Truncate() error
}

type Received interface {
	Create(pseudonymID int, text, senderTag string) (*model.Received, error)
	Delete(id int) (*model.Received, error)
	GetAll(pseudonymID int) ([]*model.Received, error)
	Truncate() error
}

type Config interface {
	DeleteHistoryAfterQuit() bool
	SetDeleteHistoryAfterQuit(b bool)
	UseProxy(b bool)
	SetProxy(s string)
	UsingProxy() bool
	GetProxy() string
}

type Service struct {
	NymClient
	Contacts
	Pseudonyms
	Sent
	Replies
	Received
	Config
}

func NewService(repo *repository.Repository, state *state.State) *Service {
	return &Service{
		NymClient:  NewNymClientController(repo, state),
		Pseudonyms: NewPseudonymsController(repo, state),
		Contacts:   NewContactsController(repo, state),
		Sent:       NewSentService(repo, state),
		Replies:    NewRepliesService(repo, state),
		Received:   NewReceivedService(repo, state),
		Config:     NewConfigService(state),
	}
}
