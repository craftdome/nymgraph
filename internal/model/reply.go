package model

import "github.com/craftdome/nymgraph/internal/entity"

type Reply struct {
	Reply    *entity.Reply
	Received *Received
}
