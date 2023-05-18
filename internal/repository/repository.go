package repository

import (
	"github.com/Tyz3/nymgraph/internal/entity"
	"github.com/Tyz3/nymgraph/internal/repository/sqlite/contacts"
	"github.com/Tyz3/nymgraph/internal/repository/sqlite/incoming_messages"
	"github.com/Tyz3/nymgraph/internal/repository/sqlite/outcoming_messages"
	"github.com/Tyz3/nymgraph/internal/repository/sqlite/pseudonyms"
	"github.com/Tyz3/nymgraph/pkg/client"
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
	GetAll() ([]*entity.Contact, error)
}

type IncomingMessages interface {
	Create(dto incoming_messages.CreateDTO) (*entity.IncomingMessage, error)
	Delete(dto incoming_messages.DeleteDTO) (*entity.IncomingMessage, error)
	Get(dto incoming_messages.GetDTO) (*entity.IncomingMessage, error)
	GetAll() ([]*entity.IncomingMessage, error)
}

type OutcomingMessages interface {
	Create(dto outcoming_messages.CreateDTO) (*entity.OutcomingMessage, error)
	Delete(dto outcoming_messages.DeleteDTO) (*entity.OutcomingMessage, error)
	Get(dto outcoming_messages.GetDTO) (*entity.OutcomingMessage, error)
	GetAll() ([]*entity.OutcomingMessage, error)
}

type Repository struct {
	Pseudonyms
	Contacts
	IncomingMessages
	OutcomingMessages
}

func NewRepository(client client.Client) *Repository {
	return &Repository{
		Pseudonyms:        pseudonyms.NewRepo(client),
		Contacts:          contacts.NewRepo(client),
		IncomingMessages:  incoming_messages.NewRepo(client),
		OutcomingMessages: outcoming_messages.NewRepo(client),
	}
}
