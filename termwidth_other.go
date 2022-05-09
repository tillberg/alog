// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd,!solaris

package alog

import (
	"io"
	"os"
	"strconv"
)

func getTermWidth(writer io.Writer) int {
	envColumns := os.Getenv("COLUMNS")
	if envColumns != "" {
		num, _ := strconv.Atoi(envColumns)
		if num != 0 {
			return num
		}
	}
	ws := getWriterState(writer)
	if ws.termWidth != 0 {
		return ws.termWidth
	}
	return 200
}
