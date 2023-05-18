package utils

import "strings"

const (
	GREEN      = "\u001B[1;32m"
	YELLOW     = "\u001B[1;33m"
	RED        = "\u001B[1;31m"
	CYAN       = "\u001B[1;36m"
	DARKGREY   = "\u001B[2;37m"
	WHITEBLACK = "\u001B[2;30;47m"

	RESET = "\u001B[0m"
)

var (
	COLORS = []string{GREEN, YELLOW, RED, CYAN, DARKGREY, WHITEBLACK, RESET}

	colorsReplacer = strings.NewReplacer(
		GREEN, "",
		YELLOW, "",
		RED, "",
		CYAN, "",
		DARKGREY, "",
		WHITEBLACK, "",
		RESET, "",
	)
)

func RemoveColors(s string) string {
	return colorsReplacer.Replace(s)
}
