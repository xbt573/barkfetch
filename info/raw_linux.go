//go:build linux

package info

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"syscall"
)

// Regexes used for extraction memory info from /proc/meminfo
var (
	getTotalMemRegex     = regexp.MustCompile(`MemTotal:\s+(\d+) kB`)
	getAvailableMemRegex = regexp.MustCompile(`MemAvailable:\s+(\d+) kB`)
	getIdRegex           = regexp.MustCompile(`(?m)^ID=\"?([^\"]*?)\"?$`)
	getPrettyNameRegex   = regexp.MustCompile(`(?m)^PRETTY_NAME=\"?([^\"]*?)\"?$`)
)

// Convert array of int8 to string
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

// Returns OS kernel type and it's version
func getRawKernel() (string, error) {
	var info syscall.Utsname

	err := syscall.Uname(&info)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Linux %v", int8ToStr(info.Release[:])), nil
}

// Returns system uptime in seconds
func getRawUptime() (int64, error) {
	var info syscall.Sysinfo_t

	err := syscall.Sysinfo(&info)
	if err != nil {
		return -1, err
	}

	return info.Uptime, nil
}

// Returns used shell
func getRawShell() string {
	return os.Getenv("SHELL")
}

// Returns used and total memory in megabytes
func getRawMemory() (used, total int, err error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return
	}

	raw, err := io.ReadAll(f)
	if err != nil {
		return
	}

	contents := string(raw)
	match := getTotalMemRegex.FindStringSubmatch(contents)

	totalMem, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	match = getAvailableMemRegex.FindStringSubmatch(contents)
	availableMem, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	total = totalMem / 1024
	used = (totalMem - availableMem) / 1024

	return
}

// Returns OS pretty name
func getRawPrettyName() (string, error) {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return "", err
	}
	defer f.Close()

	raw, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	contents := string(raw)
	match := getPrettyNameRegex.FindStringSubmatch(contents)
	return match[1], nil
}

// Guesses distro by /etc/os-release values
func guessDistro() string {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return ""
	}
	defer f.Close()

	raw, err := io.ReadAll(f)
	if err != nil {
		return ""
	}

	contents := string(raw)
	match := getIdRegex.FindStringSubmatch(contents)

	return match[1]
}
