package received

import (
	"database/sql"
)

type CreateDTO struct {
	PseudonymID int
	Text        string
	SenderTag   sql.NullString
}

type DeleteDTO struct {
	ID int
}

type GetDTO struct {
	ID int
}

type GetAllDTO struct {
	PseudonymID int
}
