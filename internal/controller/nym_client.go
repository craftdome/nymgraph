package controller

import (
	"fmt"
	"github.com/Tyz3/go-nym"
	"github.com/Tyz3/go-nym/response"
	"github.com/Tyz3/go-nym/tags"
	"github.com/Tyz3/nymgraph/internal/entity"
	"github.com/Tyz3/nymgraph/internal/repository"
	"github.com/Tyz3/nymgraph/internal/state"
	"github.com/pkg/errors"
)

type ConnectionState int

const (
	Closed ConnectionState = iota
	Connected
	Listening
)

type NymClientController struct {
	repo   *repository.Repository
	state  *state.State
	client *nym.Client

	connectionState ConnectionState

	stopped chan struct{}

	onErrorResponse       func(error *response.Error)
	onSelfAddressResponse func(error *response.SelfAddress)
	onReceiveResponse     func(error *response.Received)
	onCloseChannel        func()
}

func NewNymClientController(repo *repository.Repository, state *state.State) NymClient {
	return &NymClientController{
		repo:    repo,
		state:   state,
		stopped: make(chan struct{}),
	}
}

func (c *NymClientController) Dial(client *entity.Pseudonym) error {
	c.state.SelectedClient = client
	c.client = nym.NewClient(c.state.SelectedClient.DSN())
	if err := c.client.Dial(); err != nil {
		return errors.Wrap(err, "client.Dial")
	} else {
		c.connectionState = Connected
		return nil
	}
}

func (c *NymClientController) ListenAndServe() error {
	if c.connectionState == Closed {
		return ErrEstablishConnectionFirst
	}

	if c.connectionState == Listening {
		return ErrStateAlreadyListening
	}

	c.connectionState = Listening

	go c.client.ListenAndServe()

	for {
		select {
		case <-c.stopped:
			c.connectionState = Connected
			return nil
		default:
			msg, ok := <-c.client.Messages()
			if !ok {
				c.client.Close()
				c.connectionState = Closed
				if c.onCloseChannel != nil {
					c.onCloseChannel()
				}
				return ErrMessagesChannelWasClosed
			}

			switch msg.Type() {
			case tags.Error:
				res := msg.(*response.Error)
				c.onErrorResponse(res)
			case tags.Received:
				res := msg.(*response.Received)
				c.onReceiveResponse(res)
			case tags.SelfAddress:
				res := msg.(*response.SelfAddress)
				c.onSelfAddressResponse(res)
			}
		}
	}
}

func (c *NymClientController) StopListening() error {
	if c.connectionState != Listening {
		return ErrStateNotListening
	}

	c.stopped <- struct{}{}
	return nil
}

func (c *NymClientController) Close() error {
	fmt.Println(c.connectionState)
	if c.connectionState == Closed {
		return ErrStateAlreadyClosed
	}

	if c.connectionState == Listening {
		if err := c.StopListening(); err != nil {
			return errors.Wrap(err, "StopListening")
		}
	}

	if c.connectionState == Connected {
		if err := c.client.Close(); err != nil {
			return errors.Wrap(err, "client.Close")
		}

		c.connectionState = Closed
		if c.onCloseChannel != nil {
			c.onCloseChannel()
		}
	}
	return nil
}

func (c *NymClientController) SetErrorHandler(f func(*response.Error)) {
	c.onErrorResponse = f
}

func (c *NymClientController) SetReceivedHandler(f func(*response.Received)) {
	c.onReceiveResponse = f
}

func (c *NymClientController) SetSelfAddressHandler(f func(*response.SelfAddress)) {
	c.onSelfAddressResponse = f
}

func (c *NymClientController) SetCloseChannelHandler(f func()) {
	c.onCloseChannel = f
}

func (c *NymClientController) SendRequestAsText(r nym.Request) error {
	return c.client.SendRequestAsText(r)
}

func (c *NymClientController) SendRequestAsBinary(r nym.Request) error {
	return c.client.SendRequestAsBinary(r)
}
