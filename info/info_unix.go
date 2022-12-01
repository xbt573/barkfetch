//go:build linux

package info

import (
	"barkfetch/logos"
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

	status := fmt.Sprintf(
		"${caccent}Kernel${creset}: %v %v\n",
		int8ToStr(uname.Sysname[:]),
		int8ToStr(uname.Release[:]),
	)

	colored := os.Expand(status, logos.ColorExpand)

	return colored, nil
}

func GetUptime() (string, error) {
	var sysinfo syscall.Sysinfo_t

	if err := syscall.Sysinfo(&sysinfo); err != nil {
		return "", err
	}

	status := fmt.Sprintf(
		"${caccent}Uptime${creset}: %v minutes\n",
		int(sysinfo.Uptime/60),
	)

	colored := os.Expand(status, logos.ColorExpand)

	return colored, nil
}

func GetShell() string {
	status := fmt.Sprintf(
		"${caccent}Shell${creset}: %v\n",
		filepath.Base(os.Getenv("SHELL")),
	)

	colored := os.Expand(status, logos.ColorExpand)

	return colored
}

func GetMemory() (string, error) {
	var sysinfo syscall.Sysinfo_t

	if err := syscall.Sysinfo(&sysinfo); err != nil {
		return "", err
	}

	status := fmt.Sprintf(
		"${caccent}Memory${creset}: %v/%v MiB\n",
		int((sysinfo.Totalram-sysinfo.Freeram)*uint64(sysinfo.Unit)/1000000),
		int(sysinfo.Totalram*uint64(sysinfo.Unit)/1000000),
	)

	colored := os.Expand(status, logos.ColorExpand)

	return colored, nil
}
