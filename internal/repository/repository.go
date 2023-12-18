package repository

import (
	"github.com/craftdome/nymgraph/internal/entity"
	"github.com/craftdome/nymgraph/internal/repository/sqlite/contacts"
	"github.com/craftdome/nymgraph/internal/repository/sqlite/pseudonyms"
	"github.com/craftdome/nymgraph/internal/repository/sqlite/received"
	"github.com/craftdome/nymgraph/internal/repository/sqlite/replies"
	"github.com/craftdome/nymgraph/internal/repository/sqlite/sent"
	"github.com/craftdome/nymgraph/pkg/client"
)

type Pseudonyms interface {
	Create(dto pseudonyms.CreateDTO) (*entity.Pseudonym, error)
	Update(dto pseudonyms.UpdateDTO) (*entity.Pseudonym, error)
	Delete(dto pseudonyms.DeleteDTO) (*entity.Pseudonym, error)
	Get(dto pseudonyms.GetDTO) (*entity.Pseudonym, error)
	GetAll() ([]*entity.Pseudonym, error)
}

type Contacts interface {
	Create(dto contacts.CreateDTO) (*entity.Contact, error)
	Update(dto contacts.UpdateDTO) (*entity.Contact, error)
	Delete(dto contacts.DeleteDTO) (*entity.Contact, error)
	Get(dto contacts.GetDTO) (*entity.Contact, error)
	GetAll(dto contacts.GetAllDTO) ([]*entity.Contact, error)
}

type Received interface {
	Create(dto received.CreateDTO) (*entity.Received, error)
	Delete(dto received.DeleteDTO) (*entity.Received, error)
	Get(dto received.GetDTO) (*entity.Received, error)
	GetAll(dto received.GetAllDTO) ([]*entity.Received, error)
	Truncate() error
}

type Sent interface {
	Create(dto sent.CreateDTO) (*entity.Sent, error)
	Delete(dto sent.DeleteDTO) (*entity.Sent, error)
	GetAll(dto sent.GetAllDTO) ([]*entity.Sent, error)
	Truncate() error
}

type Replies interface {
	Create(dto replies.CreateDTO) (*entity.Reply, error)
	Delete(dto replies.DeleteDTO) (*entity.Reply, error)
	GetAll(dto replies.GetAllDTO) ([]*entity.Reply, error)
	Truncate() error
}

type Repository struct {
	Pseudonyms
	Contacts
	Received
	Sent
	Replies
}

func NewRepository(client client.Client) *Repository {
	return &Repository{
		Pseudonyms: pseudonyms.NewRepo(client),
		Contacts:   contacts.NewRepo(client),
		Received:   received.NewRepo(client),
		Sent:       sent.NewRepo(client),
		Replies:    replies.NewRepo(client),
	}
}
