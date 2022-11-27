//go:build linux

package info

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

func int8ToStr(arr []int8) string {
	b := make([]byte, 0, len(arr))
	for _, v := range arr {
		if v == 0x00 {
			break
		}
		b = append(b, byte(v))
	}
	return string(b)
}

func GetKernel() (string, error) {
	var uname syscall.Utsname
	if err := syscall.Uname(&uname); err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Kernel: %v %v\n",
		int8ToStr(uname.Sysname[:]),
		int8ToStr(uname.Release[:]),
	), nil
}

func GetUptime() (string, error) {
	var sysinfo syscall.Sysinfo_t

	if err := syscall.Sysinfo(&sysinfo); err != nil {
		return "", err
	}

	return fmt.Sprintf("Uptime: %v minutes\n", int(sysinfo.Uptime/60)), nil
}

func GetShell() string {
	return fmt.Sprintf("Shell: %v\n", filepath.Base(os.Getenv("SHELL")))
}

func GetMemory() (string, error) {
	var sysinfo syscall.Sysinfo_t

	if err := syscall.Sysinfo(&sysinfo); err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Memory: %v/%v MB",
		int((sysinfo.Totalram-sysinfo.Freeram)/1000000),
		int(sysinfo.Totalram/1000000),
	), nil
}
