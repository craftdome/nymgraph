package entity

import "time"

type OutcomingMessage struct {
	ID        int
	ContactID int
	CreateAt  time.Time
	Text      string
}
