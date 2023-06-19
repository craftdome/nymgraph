package contacts

type CreateDTO struct {
	PseudonymID int
	Address     string
	Alias       string
}

type UpdateDTO struct {
	ID      int
	Address string
	Alias   string
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
