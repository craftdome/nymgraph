package replies

type CreateDTO struct {
	ReceivedID int
	Text       string
}

type DeleteDTO struct {
	ID int
}

type GetAllDTO struct {
	ReceivedID int
}
