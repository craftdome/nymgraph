package sent

import (
	"github.com/Tyz3/nymgraph/internal/entity"
	"time"
)

type Entity struct {
	ID        int
	ContactID int
	CreateAt  int64
	Text      string
}

func (e *Entity) ToDomain() *entity.Sent {
	return &entity.Sent{
		ID:        e.ID,
		ContactID: e.ContactID,
		CreateAt:  time.Unix(e.CreateAt, 0),
		Text:      e.Text,
	}
}
