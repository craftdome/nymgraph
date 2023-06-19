package nym_client

import (
	"github.com/Tyz3/go-nym"
	"github.com/Tyz3/go-nym/response"
	"github.com/Tyz3/go-nym/tags"
	"github.com/Tyz3/nymgraph/internal/repository"
	"github.com/pkg/errors"
)

type ConnectionState int

const (
	Closed ConnectionState = iota
	Connected
	Listening
)

type ClientConnect struct {
	repo   *repository.Repository
	client *nym.Client

	selfAddress     string
	dsn             string
	connectionState ConnectionState

	OnErrorCallback       func(error *response.Error)
	OnSelfAddressCallback func(error *response.SelfAddress)
	OnReceiveCallback     func(error *response.Received)
	OnCloseCallback       func()
}

func NewClientConnect(repo *repository.Repository, dsn string) *ClientConnect {
	c := &ClientConnect{
		repo: repo,
		dsn:  dsn,

		OnErrorCallback:       func(*response.Error) {},
		OnSelfAddressCallback: func(*response.SelfAddress) {},
		OnReceiveCallback:     func(*response.Received) {},
		OnCloseCallback:       func() {},
	}

	return c
}

func (c *ClientConnect) IsOnline() bool {
	return c.connectionState != Closed
}

func (c *ClientConnect) Dial() error {
	c.client = nym.NewClient(c.dsn)
	if err := c.client.Dial(); err != nil {
		return errors.Wrap(err, "client.Dial")
	} else {
		c.connectionState = Connected
		return nil
	}
}

func (c *ClientConnect) ListenAndServe() error {
	switch c.connectionState {
	case Closed:
		return ErrEstablishConnectionFirst
	case Listening:
		return ErrStateAlreadyListening
	default:
		c.connectionState = Listening
	}

	go c.client.ListenAndServe()

	if err := c.client.SendRequestAsText(nym.NewGetSelfAddress()); err != nil {
		c.connectionState = Connected
		return errors.Wrap(err, "client.SendRequestAsText")
	}

	go func() {
		for msg := range c.client.Messages() {
			switch msg.Type() {
			case tags.Error:
				res := msg.(*response.Error)
				c.OnErrorCallback(res)
			case tags.Received:
				res := msg.(*response.Received)
				c.OnReceiveCallback(res)
			case tags.SelfAddress:
				res := msg.(*response.SelfAddress)
				c.selfAddress = res.Address
				c.OnSelfAddressCallback(res)
			}
		}

		c.connectionState = Closed
		c.OnCloseCallback()
	}()

	return nil
}

func (c *ClientConnect) SelfAddress() string {
	return c.selfAddress
}

func (c *ClientConnect) Close() error {
	switch c.connectionState {
	case Closed:
		return ErrStateAlreadyClosed
	default:
		c.connectionState = Closed
	}

	if err := c.client.Close(); err != nil {
		return errors.Wrap(err, "client.Close")
	}

	c.connectionState = Closed
	c.OnCloseCallback()

	return nil
}

func (c *ClientConnect) SendMessage(text, recipient string, replySURBs int) error {
	var req nym.Request
	if replySURBs != 0 {
		req = nym.NewSendAnonymous(text, recipient, replySURBs)
	} else {
		req = nym.NewSend(text, recipient)
	}

	// Отправляем сообщение в nym-client
	if err := c.client.SendRequestAsText(req); err != nil {
		return errors.Wrap(err, "client.SendRequestAsText")
	}

	return nil
}

func (c *ClientConnect) SendReply(text, senderTag string) error {
	req := nym.NewReply(text, senderTag)

	// Отправляем сообщение в nym-client
	if err := c.client.SendRequestAsText(req); err != nil {
		return errors.Wrap(err, "client.SendRequestAsText")
	}

	return nil
}

//func (c *ClientConnect) SendReply(receivedID int, text string) (*model.Reply, error) {
//	dto := received.GetDTO{ID: receivedID}
//	rec, err := c.repo.Received.Get(dto)
//	if err != nil {
//		return nil, errors.Wrapf(err, "repo.Received.Get %+v", dto)
//	}
//
//	if rec.SenderTag == "" {
//		return nil, errors.Wrapf(ErrNoSenderTag, "receivedID %d", receivedID)
//	}
//
//	req := nym.NewReply(text, rec.SenderTag)
//	// Отправляем сообщение в nym-client
//	if err := c.client.SendRequestAsText(req); err != nil {
//		return nil, errors.Wrap(err, "client.SendRequestAsText")
//	}
//
//	dto2 := replies.CreateDTO{
//		ReceivedID: receivedID,
//		Text:       text,
//	}
//
//	if reply, err := c.repo.Replies.Create(dto2); err != nil {
//		return nil, errors.Wrapf(err, "repo.Replies.Create %+v", dto2)
//	} else {
//		dto3 := pseudonyms.GetDTO{ID: rec.PseudonymID}
//		ps, err := c.repo.Pseudonyms.Get(dto3)
//		if err != nil {
//			return nil, errors.Wrapf(err, "repo.Pseudonyms.Get %+v", dto3)
//		}
//
//		return &model.Reply{
//			Reply: reply,
//			Received: &model.Received{
//				Received:  rec,
//				Pseudonym: ps,
//			},
//		}, nil
//	}
//}
//
//func (c *ClientConnect) SendMessageToContact(contactID int, text string, replySURBs int) (*model.Sent, error) {
//	dto := contacts.GetDTO{ID: contactID}
//	contact, err := c.repo.Contacts.Get(dto)
//	if err != nil {
//		return nil, errors.Wrapf(err, "repo.Contacts.Get %+v", dto)
//	}
//
//	var req nym.Request
//	if replySURBs != 0 {
//		req = nym.NewSendAnonymous(text, contact.Address, replySURBs)
//	} else {
//		req = nym.NewSend(text, contact.Address)
//	}
//
//	// Отправляем сообщение в nym-client
//	if err := c.client.SendRequestAsText(req); err != nil {
//		return nil, errors.Wrap(err, "client.SendRequestAsText")
//	}
//
//	dto2 := sent.CreateDTO{
//		Text:      text,
//		ContactID: contactID,
//	}
//
//	// Добавляем сообщение в БД (контакт ID не указывается)
//	if msg, err := c.repo.Sent.Create(dto2); err != nil {
//		return nil, errors.Wrapf(err, "repo.OutcomingMessages.Create %+v", dto)
//	} else {
//		return &model.Sent{
//			Sent:    msg,
//			Contact: contact,
//		}, nil
//	}
//}
