package outcoming_messages

import (
	"github.com/Tyz3/nymgraph/internal/entity"
	"time"
)

type Entity struct {
	ID        int
	ContactID int
	CreateAt  time.Time
	Text      string
}

func (e *Entity) ToDomain() *entity.OutcomingMessage {
	return &entity.OutcomingMessage{
		ID:        e.ID,
		ContactID: e.ContactID,
		CreateAt:  e.CreateAt,
		Text:      e.Text,
	}
}
