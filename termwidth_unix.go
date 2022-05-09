// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package alog

import (
	"io"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

// getTermWidth returns the dimensions of the given terminal.
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
	var fd uintptr
	if writer == os.Stdout {
		fd = uintptr(syscall.Stdout)
	} else {
		// For custom writers, just use the width we get for stderr. This might not be true in some
		// cases (and for those cases, we should add an option to explicitly set width), but it will
		// be true in most cases.
		fd = uintptr(syscall.Stderr)
	}
	var dimensions [4]uint16
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&dimensions)), 0, 0, 0); err != 0 {
		// Fall back to a width of 200
		return 200
	}
	if int(dimensions[1]) == 0 {
		// Seen inside a rkt container, where the syscall returns a width of 0, which isn't helpful.
		return 200
	}
	return int(dimensions[1])
}
