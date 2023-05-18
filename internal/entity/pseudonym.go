package entity

import "fmt"

type Pseudonym struct {
	ID     int
	Name   string `validate:"required,min=1,max=16"`
	Server string `validate:"required,hostname_port"`
}

func (e *Pseudonym) DSN() string {
	return fmt.Sprintf("ws://%s", e.Server)
}

func (e *Pseudonym) DataChanged() {

}
