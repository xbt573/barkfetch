//go:build darwin

package info

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

// Extracts boottime from "sysctl kern.boottime"
var extractBoottimeRegex = regexp.MustCompile(`kern\.boottime: \{ sec = (\d+)`)

// Various regexes to extract info from "vm_stat"
var (
	extractWiredMemoryRegex  = regexp.MustCompile(`Pages wired down:\s+(\d+)\.`)
	extractActiveMemoryRegex = regexp.MustCompile(`Pages active:\s+(\d+)\.`)
	extractCompressedMemory  = regexp.MustCompile(`Pages occupied by compressor:\s+(\d+)\.`)
)

// Extract GPU model from "system_profiler SPDisplaysDataType"
var extractChipsetModel = regexp.MustCompile(`Chipset Model: (.*)`)

// Returns OS kernel type and it's version
func getRawKernel() (string, error) {
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// Returns system uptime in seconds
func getRawUptime() (int64, error) {
	out, err := exec.Command("sysctl", "kern.boottime").Output()
	if err != nil {
		return -1, err
	}

	match := extractBoottimeRegex.FindStringSubmatch(string(out))
	seconds, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return -1, err
	}

	return time.Now().Add(time.Duration(-seconds * int64(time.Second))).Unix(), nil
}

// Returns used shell
func getRawShell() string {
	return os.Getenv("SHELL")
}

// Returns used and total memory in megabytes
func getRawMemory() (used, total int, err error) {
	totalMemoryString, err := exec.Command("sysctl", "-n", "hw.memsize").Output()
	if err != nil {
		return -1, -1, err
	}

	// newline hack
	totalMemoryString = totalMemoryString[:len(totalMemoryString)-1]
	totalMemory, err := strconv.ParseInt(string(totalMemoryString), 10, 64)
	if err != nil {
		return -1, -1, err
	}
	total = int(totalMemory / 1000000)

	vmStat, err := exec.Command("vm_stat").Output()
	if err != nil {
		return -1, -1, err
	}

	match := extractWiredMemoryRegex.FindStringSubmatch(string(vmStat))
	wired, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return -1, -1, err
	}

	match = extractActiveMemoryRegex.FindStringSubmatch(string(vmStat))
	active, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return -1, -1, err
	}

	match = extractCompressedMemory.FindStringSubmatch(string(vmStat))
	compressed, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return -1, -1, err
	}

	used = int((wired + active + compressed) * 4 / 1024)
	return
}

// Returns CPU model (currently first)
func getRawCpu() (string, error) {
	out, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output()
	if err != nil {
		return "", err
	}

	return string(out[:len(out)-1]), nil
}

// Returns GPU manufacturer and model
func getRawGpu() ([]string, error) {
	out, err := exec.Command("system_profiler", "SPDisplaysDataType").Output()
	if err != nil {
		return []string{}, err
	}

	contents := string(out)
	match := extractChipsetModel.FindStringSubmatch(contents)
	return []string{match[1]}, nil
}

// Returns OS pretty name
func getRawPrettyName() (string, error) {
	out, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return "", nil
	}

	// newline hack
	out = out[:len(out)-1]
	return fmt.Sprintf("macOS %v", string(out)), nil
}

// Returns "mac" always lol
func guessDistro() string {
	return "mac"
}
