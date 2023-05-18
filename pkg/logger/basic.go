package logger

import (
	"fmt"
	"os"
)

type LogLevel int

const (
	Minimum LogLevel = iota
	Critical
	Error
	Warning
	Info
	Debug
)

func (l LogLevel) Prefix() string {
	switch l {
	case Minimum:
		return "MINIMUM"
	case Critical:
		return "CRITICAL"
	case Error:
		return "ERROR"
	case Warning:
		return "WARN"
	case Info:
		return "INFO"
	case Debug:
		return "DEBUG"
	default:
		return ""
	}
}

var (
	Mode = Info
)

func SetLogLevel(name string) {
	switch name {
	case "MINIMUM":
		Mode = Minimum
	case "CRITICAL":
		Mode = Critical
	case "ERROR":
		Mode = Error
	case "WARNING":
		Mode = Warning
	case "INFO":
		Mode = Info
	case "DEBUG":
		Mode = Debug
	default:
		Mode = Info
	}
}

func Printf(l LogLevel, format string, a ...any) {
	if Mode < l {
		return
	}

	switch l {
	case Critical:
		fallthrough
	case Error:
		fmt.Fprintf(os.Stderr, "[%s] %s\n", l.Prefix(), fmt.Sprintf(format, a...))
	default:
		fmt.Fprintf(os.Stdout, "[%s] %s\n", l.Prefix(), fmt.Sprintf(format, a...))
	}

}
