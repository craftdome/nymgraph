package nym_client

import (
	"github.com/Tyz3/go-nym"
	"github.com/Tyz3/go-nym/response"
	"github.com/Tyz3/go-nym/tags"
	"github.com/pkg/errors"
)

type ConnectionState int

const (
	Closed ConnectionState = iota
	Connected
	Listening
)

type ClientConnect struct {
	client *nym.Client

	selfAddress     string
	dsn             string
	connectionState ConnectionState

	OnErrorCallback       func(error *response.Error)
	OnSelfAddressCallback func(error *response.SelfAddress)
	OnReceiveCallback     func(error *response.Received)
	OnCloseCallback       func()
}

func NewClientConnect(dsn string) *ClientConnect {
	c := &ClientConnect{
		dsn: dsn,

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
