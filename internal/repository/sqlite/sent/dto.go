package sent

type CreateDTO struct {
	ContactID int
	Text      string
}

type DeleteDTO struct {
	ID int
}

type GetAllDTO struct {
	ContactID int
}
