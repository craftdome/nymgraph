package replies

import (
	"github.com/craftdome/nymgraph/internal/entity"
	"time"
)

type Entity struct {
	ID         int
	ReceivedID int
	CreateAt   int64
	Text       string
}

func (e *Entity) ToDomain() *entity.Reply {
	return &entity.Reply{
		ID:         e.ID,
		ReceivedID: e.ReceivedID,
		CreateAt:   time.Unix(e.CreateAt, 0),
		Text:       e.Text,
	}
}
