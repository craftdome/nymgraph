package sent

type CreateDTO struct {
	ContactID  int
	Text       string
	ReplySurbs int
}

type DeleteDTO struct {
	ID int
}

type GetAllDTO struct {
	ContactID int
}
