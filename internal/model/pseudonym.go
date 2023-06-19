package model

import (
	"fyne.io/fyne/v2"
	"github.com/Tyz3/nymgraph/internal/entity"
	"github.com/Tyz3/nymgraph/internal/nym_client"
)

type Pseudonym struct {
	Pseudonym *entity.Pseudonym
	NymClient *nym_client.ClientConnect
	MenuItem  *fyne.MenuItem
}
