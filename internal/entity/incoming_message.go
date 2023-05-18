package entity

import "time"

type IncomingMessage struct {
	ID        int
	CreateAt  time.Time
	Text      string
	SenderTag string
}
