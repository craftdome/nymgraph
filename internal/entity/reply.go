package entity

import "time"

type Reply struct {
	ID         int
	ReceivedID int
	CreateAt   time.Time
	Text       string
}
