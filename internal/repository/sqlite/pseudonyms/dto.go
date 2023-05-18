package pseudonyms

type CreateDTO struct {
	Name   string
	Server string
}

type UpdateDTO struct {
	ID     int
	Name   string
	Server string
}

type DeleteDTO struct {
	ID int
}

type GetDTO struct {
	ID int
}
