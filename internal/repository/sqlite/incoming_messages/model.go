package incoming_messages

import (
	"database/sql"
	"github.com/Tyz3/nymgraph/internal/entity"
	"time"
)

type Entity struct {
	ID        int
	CreateAt  time.Time
	Text      string
	SenderTag sql.NullString
}

func (e *Entity) ToDomain() *entity.IncomingMessage {
	return &entity.IncomingMessage{
		ID:        e.ID,
		CreateAt:  e.CreateAt,
		Text:      e.Text,
		SenderTag: e.SenderTag.String,
	}
}
