package logger

import (
	"fmt"
	"github.com/Tyz3/nymgraph/pkg/utils"
	"io"
	"strings"
	"time"
)

type LogMessage interface {
	Printf(format string, v ...any)
}

type info struct {
	prefix string
}

type warn struct {
	prefix string
}

type err struct {
	prefix string
}

type log struct {
	INFO   LogMessage
	WARN   LogMessage
	ERROR  LogMessage
	writer io.Writer
}

func (l *log) Ok() {
	fmt.Println(utils.GREEN + "OK" + utils.RESET)
}

func (l *log) Failed() {
	fmt.Println(utils.RED + "FAILED" + utils.RESET)
}

func (l *log) SetWriter(w io.Writer) {
	l.writer = w
}

func prnt(prefix, format string, v ...any) {
	date := time.Now().Format(DateTimeLayout)
	var sb strings.Builder
	sb.WriteString(date)
	sb.WriteString(" [")
	sb.WriteString(prefix)
	sb.WriteString("] ")
	msg := strings.ReplaceAll(fmt.Sprintf(format, v...), "\n", "\t\n")
	sb.WriteString(msg)
	fmt.Fprint(Log.writer, sb.String())
}

func (l *info) Printf(format string, v ...any) {
	prnt(l.prefix, format, v...)
}

func (l *warn) Printf(format string, v ...any) {
	prnt(l.prefix, format, v...)
}

func (l *err) Printf(format string, v ...any) {
	prnt(l.prefix, format, v...)
}
