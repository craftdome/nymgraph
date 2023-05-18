package pseudonyms

import (
	"github.com/Tyz3/nymgraph/internal/entity"
)

type Entity struct {
	ID     int
	Name   string
	Server string
}

func (e *Entity) ToDomain() *entity.Pseudonym {
	return &entity.Pseudonym{
		ID:     e.ID,
		Name:   e.Name,
		Server: e.Server,
	}
}
