//go:build linux

package info

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

// Regexes used for extraction memory info from /proc/meminfo
var (
	getTotalMemRegex     = regexp.MustCompile(`MemTotal:\s+(\d+) kB`)
	getFreeMemRegex      = regexp.MustCompile(`MemFree:\s+(\d+) kB`)
	getShMemRegex        = regexp.MustCompile(`Shmem:\s+(\d+) kB`)
	getBuffersRegex      = regexp.MustCompile(`Buffers:\s+(\d+) kB`)
	getCachedRegex       = regexp.MustCompile(`Cached:\s+(\d+) kB`)
	getSReclaimableRegex = regexp.MustCompile(`SReclaimable:\s+(\d+) kB`)
	getIdRegex           = regexp.MustCompile(`(?m)^ID=\"?([^\"]*?)\"?$`)
	getPrettyNameRegex   = regexp.MustCompile(`(?m)^PRETTY_NAME=\"?([^\"]*?)\"?$`)
)

// Regex used to extract CPU model from /proc/cpuinfo
var getCpuModelRegex = regexp.MustCompile(`(?m)^model name\s+: (.*)$`)

// Regex used to extract raw GPU manufacturer and model from "lspci" output
var getRawGpuManufacturerAndModelRegex = regexp.MustCompile(`.*"(?:Display|3D|VGA).*?" "(.*?)" "(.*?)"`)

// Regexes used to extract pretty GPU manufacturer and model from raw input
var (
	getGpuModelRegex = regexp.MustCompile(`.*\[(.*)\]`)
)

// Regex used to extract resolutions from "xrandr --nograb --current"
var getResolutionRegex = regexp.MustCompile(`connected(?: primary)? (\d+)x(\d+)`)

// Regex used to remove too much spaces between words
var removeExtraSpacesRegex = regexp.MustCompile(`\s+`)

// Convert array of (u)int8 to string
func int8ToStr(arr []int) string {
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
func getRawKernel() string {
	var info syscall.Utsname

	err := syscall.Uname(&info)
	if err != nil {
		return "n/a"
	}

	intArr := []int{}

	for _, val := range info.Release {
		intArr = append(intArr, int(val))
	}

	return fmt.Sprintf("Linux %v", int8ToStr(intArr))
}

// Returns system uptime in seconds
func getRawUptime() int {
	var info syscall.Sysinfo_t

	err := syscall.Sysinfo(&info)
	if err != nil {
		return -1
	}

	return int(info.Uptime)
}

// Returns used shell
func getRawShell() string {
	return os.Getenv("SHELL")
}

// Returns used and total memory in megabytes
func getRawMemory() (used, total int) {
	raw, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return
	}

	contents := string(raw)

	match := getTotalMemRegex.FindStringSubmatch(contents)
	if len(match) == 0 {
		return
	}
	totalMem, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	match = getShMemRegex.FindStringSubmatch(contents)
	if len(match) == 0 {
		return
	}
	shMem, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	match = getFreeMemRegex.FindStringSubmatch(contents)
	if len(match) == 0 {
		return
	}
	freeMem, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	match = getBuffersRegex.FindStringSubmatch(contents)
	if len(match) == 0 {
		return
	}
	buffers, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	match = getCachedRegex.FindStringSubmatch(contents)
	if len(match) == 0 {
		return
	}
	cached, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	match = getSReclaimableRegex.FindStringSubmatch(contents)
	if len(match) == 0 {
		return
	}
	sReclaimable, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	// match = getAvailableMemRegex.FindStringSubmatch(contents)
	// availableMem, err := strconv.Atoi(match[1])
	// if err != nil {
	// 	return
	// }

	total = totalMem / 1000
	used = (totalMem + shMem - freeMem - buffers - cached - sReclaimable) / 1000

	return
}

// Returns CPU model (currently first)
func getRawCpu() string {
	raw, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "n/a"
	}

	contents := string(raw)
	match := getCpuModelRegex.FindStringSubmatch(contents)
	if len(match) == 0 {
		return "n/a"
	}

	return removeExtraSpacesRegex.ReplaceAllString(match[1], " ")
}

// Returns GPU manufacturer and model
func getRawGpus() []string {
	out, err := exec.Command("lspci", "-mm").Output()
	if err != nil {
		return []string{}
	}

	contents := string(out)
	match := getRawGpuManufacturerAndModelRegex.FindAllStringSubmatch(contents, -1)
	if len(match) == 0 {
		return []string{}
	}

	gpus := []string{}
	for _, line := range match {
		var manufacturer string

		if strings.Contains(line[1], "Intel") {
			manufacturer = "Intel"
		}

		if strings.Contains(line[1], "NVIDIA") {
			manufacturer = "NVIDIA"
		}

		if strings.Contains(line[1], "AMD") {
			manufacturer = "AMD"
		}

		modelMatch := getGpuModelRegex.FindStringSubmatch(line[2])
		if len(modelMatch) > 0 {
			gpus = append(gpus, fmt.Sprintf("%v %v", manufacturer, modelMatch[1]))
		}
	}

	return gpus
}

// Returns main screen resolution
func getRawScreenResolutions() []string {
	out, err := exec.Command("xrandr", "--nograb", "--current").Output()
	if err != nil {
		return []string{}
	}

	contents := string(out)
	match := getResolutionRegex.FindAllStringSubmatch(contents, -1)
	if len(match) == 0 {
		return []string{}
	}

	resolutions := []string{}

	for _, resMatch := range match {
		resolutions = append(resolutions, fmt.Sprintf("%vx%v", resMatch[1], resMatch[2]))
	}

	return resolutions
}

// Returns OS pretty name
func getRawPrettyName() string {
	raw, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "n/a"
	}

	contents := string(raw)
	match := getPrettyNameRegex.FindStringSubmatch(contents)
	if len(match) == 0 {
		return "n/a"
	}

	return match[1]
}

// Guesses distro by /etc/os-release values
func guessDistro() string {
	raw, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "n/a"
	}

	contents := string(raw)
	match := getIdRegex.FindStringSubmatch(contents)
	if len(match) == 0 {
		return "n/a"
	}

	return match[1]
}
