package received

import (
	"database/sql"
	"github.com/Tyz3/nymgraph/internal/entity"
	"time"
)

type Entity struct {
	ID          int
	PseudonymID int
	CreateAt    time.Time
	Text        string
	SenderTag   sql.NullString
}

func (e *Entity) ToDomain() *entity.Received {
	return &entity.Received{
		ID:          e.ID,
		PseudonymID: e.PseudonymID,
		CreateAt:    e.CreateAt,
		Text:        e.Text,
		SenderTag:   e.SenderTag.String,
	}
}
