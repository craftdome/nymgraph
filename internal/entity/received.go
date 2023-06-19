package entity

import "time"

type Received struct {
	ID          int
	CreateAt    time.Time
	Text        string
	SenderTag   string
	PseudonymID int
}
