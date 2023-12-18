package received

import (
	"database/sql"
	"github.com/craftdome/nymgraph/internal/entity"
	"time"
)

type Entity struct {
	ID          int
	PseudonymID int
	CreateAt    int64
	Text        string
	SenderTag   sql.NullString
}

func (e *Entity) ToDomain() *entity.Received {
	return &entity.Received{
		ID:          e.ID,
		PseudonymID: e.PseudonymID,
		CreateAt:    time.Unix(e.CreateAt, 0),
		Text:        e.Text,
		SenderTag:   e.SenderTag.String,
	}
}
