package contacts

import (
	"database/sql"
	"github.com/Tyz3/nymgraph/internal/entity"
)

type Entity struct {
	ID          int
	PseudonymID int
	Address     string
	Alias       sql.NullString
}

func (e *Entity) ToDomain() *entity.Contact {
	return &entity.Contact{
		ID:          e.ID,
		PseudonymID: e.PseudonymID,
		Address:     e.Address,
		Alias:       e.Alias.String,
	}
}
