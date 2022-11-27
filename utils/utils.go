package utils

import (
	"os"
	"syscall"
	"unsafe"
)

type winsize struct {
	rows uint16
	cols uint16
	x    uint16
	y    uint16
}

func GetTerminalSize() (x, y int, err error) {
	ws := winsize{}

	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		os.Stdout.Fd(),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)),
	)

	if errno != 0 {
		err = errno
		return
	}

	return int(ws.cols), int(ws.rows), nil
}
