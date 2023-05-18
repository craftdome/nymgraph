package contacts

import "database/sql"

type CreateDTO struct {
	PseudonymID int
	Address     string
	Alias       sql.NullString
}

type UpdateDTO struct {
	ID    int
	Alias sql.NullString
}

type DeleteDTO struct {
	ID int
}

type GetDTO struct {
	ID int
}
