package logger

import (
	"fmt"
	"strings"
)

type list []string

type Tree struct {
	sb      strings.Builder
	content map[int]list
}

func NewTree() *Tree {
	return &Tree{content: make(map[int]list, 10)}
}

func (t *Tree) Append(level int, format string, a ...any) *Tree {
	if level == 0 {
		t.sb.WriteString("━ ")
	} else {
		t.sb.WriteString(strings.Repeat(" ", 2*level))
		t.sb.WriteString("└ ")
	}
	t.sb.WriteString(fmt.Sprintf(format, a...))
	t.sb.WriteString("\n")
	return t
}

func (t *Tree) Print() {
	fmt.Print(t.sb.String())
}

//func Print(format string, a ...any) {
//	fmt.Printf("━ %s\n", fmt.Sprintf(format, a...))
//}
