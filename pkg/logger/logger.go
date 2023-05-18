package logger

import (
	"fmt"
	"github.com/Tyz3/nymgraph/pkg/utils"
	"os"
	"strings"
)

var (
	Log = &log{
		INFO:  &info{prefix: utils.GREEN + "INFO" + utils.RESET},
		WARN:  &warn{prefix: utils.YELLOW + "WARN" + utils.RESET},
		ERROR: &err{prefix: utils.RED + "ERROR" + utils.RESET},

		writer: os.Stdout,
	}
	Memo = &memorandum{
		writer: os.Stdout,
	}
)

const DateTimeLayout = "2006-01-02 15:04:05 Z07"

func anyJoinToStr(sep string, a ...any) string {
	var sb strings.Builder
	for i := 0; i < len(a); i++ {
		sb.WriteString(fmt.Sprint(a[i]))
		if i != len(a)-1 {
			sb.WriteString(sep)
		}
	}

	return sb.String()
}
