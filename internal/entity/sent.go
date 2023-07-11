package entity

import "time"

type Sent struct {
	ID         int
	ContactID  int
	CreateAt   time.Time
	Text       string
	ReplySurbs int
}
