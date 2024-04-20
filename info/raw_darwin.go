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
	extractWiredMemoryRegex      = regexp.MustCompile(`Pages wired down:\s+(\d+)\.`)
	extractActiveMemoryRegex     = regexp.MustCompile(`Pages active:\s+(\d+)\.`)
	extractCompressedMemoryRegex = regexp.MustCompile(`Pages occupied by compressor:\s+(\d+)\.`)
)

// Extract GPU model from "system_profiler SPDisplaysDataType"
var extractChipsetModelRegex = regexp.MustCompile(`Chipset Model: (.*)`)

// Extract resolution from "system_profiler SPDisplaysDataType"
var extractResolutionRegex = regexp.MustCompile(`(?m)^\s+Resolution: (\d+) x (\d+)$`)

// Returns OS kernel type and it's version
func getRawKernel() string {
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "n/a"
	}

	return string(out)
}

// Returns system uptime in seconds
func getRawUptime() int64 {
	out, err := exec.Command("sysctl", "kern.boottime").Output()
	if err != nil {
		return -1
	}

	match := extractBoottimeRegex.FindStringSubmatch(string(out))
	seconds, err := strconv.ParseInt(match[1], 10, 64)
	if len(match) == 0 || err != nil {
		return -1
	}

	return time.Now().Add(time.Duration(-seconds * int64(time.Second))).Unix()
}

// Returns used shell
func getRawShell() string {
	return os.Getenv("SHELL")
}

// Returns used and total memory in megabytes
func getRawMemory() (used, total int) {
	totalMemoryString, err := exec.Command("sysctl", "-n", "hw.memsize").Output()
	if err != nil {
		return
	}

	// newline hack
	totalMemoryString = totalMemoryString[:len(totalMemoryString)-1]
	totalMemory, err := strconv.ParseInt(string(totalMemoryString), 10, 64)
	if err != nil {
		return
	}

	vmStat, err := exec.Command("vm_stat").Output()
	if err != nil {
		return
	}

	match := extractWiredMemoryRegex.FindStringSubmatch(string(vmStat))
	wired, err := strconv.ParseInt(match[1], 10, 64)
	if len(match) == 0 || err != nil {
		return
	}

	match = extractActiveMemoryRegex.FindStringSubmatch(string(vmStat))
	active, err := strconv.ParseInt(match[1], 10, 64)
	if len(match) == 0 || err != nil {
		return
	}

	match = extractCompressedMemoryRegex.FindStringSubmatch(string(vmStat))
	compressed, err := strconv.ParseInt(match[1], 10, 64)
	if len(match) == 0 || err != nil {
		return
	}

	total, used = int(totalMemory/1000000), int((wired+active+compressed)*4/1024)

	return
}

// Returns CPU model (currently first)
func getRawCpu() string {
	out, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output()
	if err != nil {
		return "n/a"
	}

	return string(out[:len(out)-1])
}

// Returns GPU manufacturer and model
func getRawGpus() []string {
	out, err := exec.Command("system_profiler", "SPDisplaysDataType").Output()

	contents := string(out)
	match := extractChipsetModelRegex.FindStringSubmatch(contents)

	if err != nil || len(match) == 0 {
		return []string{}
	}
	return []string{match[1]}
}

// Returns main screen resolution
func getRawScreenResolutions() []string {
	out, err := exec.Command("system_profiler", "SPDisplaysDataType").Output()

	contents := string(out)
	match := extractResolutionRegex.FindStringSubmatch(contents)
	if err != nil || len(match) == 0 {
		return []string{}
	}

	return []string{match[1]}
}

// Returns OS pretty name
func getRawPrettyName() string {
	out, err := exec.Command("sw_vers", "-productVersion").Output()
	if err != nil {
		return "n/a"
	}

	// newline hack
	out = out[:len(out)-1]
	return fmt.Sprintf("macOS %v", string(out))
}

// Returns "mac" always lol
func guessDistro() string {
	return "mac"
}
