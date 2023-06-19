package model

import "github.com/Tyz3/nymgraph/internal/entity"

type Reply struct {
	Reply    *entity.Reply
	Received *Received
}
