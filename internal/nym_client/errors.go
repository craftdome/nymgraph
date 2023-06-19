package nym_client

import "github.com/pkg/errors"

var (
	ErrEstablishConnectionFirst = errors.New("first you need to establish a connection")
	ErrStateAlreadyListening    = errors.New("nym-client already listening")
	ErrStateAlreadyClosed       = errors.New("nym-client already listening")
	ErrStateNotListening        = errors.New("nym-client not listening")
	//ErrMessagesChannelWasClosed = errors.New("nym-client closed the messages channel")
	ErrNoSenderTag = errors.New("no sender tag")
)
