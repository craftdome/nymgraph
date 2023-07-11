package model

import (
	"fmt"
	"github.com/Tyz3/nymgraph/internal/entity"
)

type Contact struct {
	Contact   *entity.Contact
	Pseudonym *entity.Pseudonym
}

func (m *Contact) Pretty() string {
	return fmt.Sprintf("(%s) %s...", m.Contact.Alias, m.Contact.Address[:7])
}
