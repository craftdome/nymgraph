package logger

import (
	"fmt"
	"github.com/Tyz3/nymgraph/pkg/utils"
	"io"
	"strings"
	"time"
)

type Memorandum interface {
	New(prefix ...any) MemoMessage
}

type memorandum struct {
	writer io.Writer
}

func (m *memorandum) New(prefix ...any) MemoMessage {
	date := time.Now().Format(DateTimeLayout)
	pref := anyJoinToStr(" ", prefix)

	memo := new(memoMessage)
	memo.Timestamps.start = time.Now().UnixMilli()
	memo.Timestamps.last = memo.Timestamps.start

	memo.message.WriteString(date)
	memo.message.WriteString(" ")
	memo.message.WriteString(utils.WHITEBLACK)
	memo.message.WriteString(" ")
	memo.message.WriteString(pref)
	memo.message.WriteString(" ")
	memo.message.WriteString(utils.RESET)

	return memo
}

func (m *memorandum) SetWriter(w io.Writer) {
	m.writer = w
}

type MemoMessage interface {
	Debug(format string, v ...any) MemoMessage
	Info(format string, v ...any) MemoMessage
	Warn(format string, v ...any) MemoMessage
	Error(format string, v ...any) MemoMessage
	Timestamp() MemoMessage
	Print()
	Ok()
	Failed()
	HiddenMode()
}

type memoMessage struct {
	Timestamps struct {
		start int64
		last  int64
	}
	message    strings.Builder
	hiddenMode bool
}

func (m *memoMessage) appendBlock(color string, spec string, format string, v ...any) MemoMessage {
	m.message.WriteString(" -> ")
	m.message.WriteString(color)
	m.message.WriteString(spec)
	m.message.WriteString("[")
	m.message.WriteString(utils.RESET)
	m.message.WriteString(fmt.Sprintf(format, v...))
	m.message.WriteString(color)
	m.message.WriteString("]")
	m.message.WriteString(utils.RESET)
	return m
}

func (m *memoMessage) Info(format string, v ...any) MemoMessage {
	return m.appendBlock(utils.GREEN, "", format, v...)
}

func (m *memoMessage) Debug(format string, v ...any) MemoMessage {
	return m.appendBlock(utils.CYAN, "~", format, v...)
}

func (m *memoMessage) Warn(format string, v ...any) MemoMessage {
	return m.appendBlock(utils.YELLOW, "!", format, v...)
}

func (m *memoMessage) Error(format string, v ...any) MemoMessage {
	return m.appendBlock(utils.RED, "?", format, v...)
}

func (m *memoMessage) Timestamp() MemoMessage {
	format := fmt.Sprintf("%.3fs", float64(time.Now().UnixMilli()-m.Timestamps.last)/1000)

	m.message.WriteString(" ")
	m.message.WriteString(utils.DARKGREY)
	m.message.WriteString(format)
	m.message.WriteString(utils.RESET)

	m.Timestamps.last = time.Now().UnixMilli()
	return m
}

func (m *memoMessage) appendStatus(color string, status string) {
	m.message.WriteString(" ")
	m.message.WriteString(color)
	m.message.WriteString(status)
	m.message.WriteString(utils.RESET)
}

func (m *memoMessage) Print() {
	fmt.Fprintln(Memo.writer, m.message.String())
}

func (m *memoMessage) Ok() {
	m.appendStatus(utils.GREEN, "OK")
	m.Print()
}

func (m *memoMessage) Failed() {
	m.appendStatus(utils.RED, "FAILED")
	m.Print()
}

func (m *memoMessage) HiddenMode() {
	m.hiddenMode = true
}
