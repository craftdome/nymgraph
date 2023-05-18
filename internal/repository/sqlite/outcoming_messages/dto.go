package outcoming_messages

type CreateDTO struct {
	ContactID int
	Text      string
}

type DeleteDTO struct {
	ID int
}

type GetDTO struct {
	ID int
}
