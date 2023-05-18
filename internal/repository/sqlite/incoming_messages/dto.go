package incoming_messages

import (
	"database/sql"
)

type CreateDTO struct {
	Text      string
	SenderTag sql.NullString
}

type DeleteDTO struct {
	ID int
}

type GetDTO struct {
	ID int
}
