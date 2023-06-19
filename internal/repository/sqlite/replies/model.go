package replies

import (
	"github.com/Tyz3/nymgraph/internal/entity"
	"time"
)

type Entity struct {
	ID         int
	ReceivedID int
	CreateAt   time.Time
	Text       string
}

func (e *Entity) ToDomain() *entity.Reply {
	return &entity.Reply{
		ID:         e.ID,
		ReceivedID: e.ReceivedID,
		CreateAt:   e.CreateAt,
		Text:       e.Text,
	}
}
