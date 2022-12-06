//go:build linux

package info

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"syscall"
)

// Regexes used for extraction memory info from /proc/meminfo
var (
	getTotalMemRegex     = regexp.MustCompile(`MemTotal:\s+(\d+) kB`)
	getFreeMemRegex      = regexp.MustCompile(`MemFree:\s+(\d+) kB`)
	getAvailableMemRegex = regexp.MustCompile(`MemAvailable:\s+(\d+) kB`)
	getShMemRegex        = regexp.MustCompile(`Shmem:\s+(\d+) kB`)
	getBuffersRegex      = regexp.MustCompile(`Buffers:\s+(\d+) kB`)
	getCachedRegex       = regexp.MustCompile(`Cached:\s+(\d+) kB`)
	getSReclaimableRegex = regexp.MustCompile(`SReclaimable:\s+(\d+) kB`)
	getIdRegex           = regexp.MustCompile(`(?m)^ID=\"?([^\"]*?)\"?$`)
	getPrettyNameRegex   = regexp.MustCompile(`(?m)^PRETTY_NAME=\"?([^\"]*?)\"?$`)
)

// Regex used to extract CPU model from /proc/cpuinfo
var getCpuModelRegex = regexp.MustCompile(`(?m)^model name\s+: (.*)$`)

// Regex used to extract GPU manufacturer and model from "lspci" output
var getGpuManufacturerAndModel = regexp.MustCompile(`.*"(?:Display|3D|VGA).*?" ".*?\[(.*)\]" ".*?\[(.*?)\]"`)

// Regex used to remove too much spaces between words
var removeExtraSpacesRegex = regexp.MustCompile(`\s+`)

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

	match = getShMemRegex.FindStringSubmatch(contents)
	shMem, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	match = getFreeMemRegex.FindStringSubmatch(contents)
	freeMem, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	match = getBuffersRegex.FindStringSubmatch(contents)
	buffers, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	match = getCachedRegex.FindStringSubmatch(contents)
	cached, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	match = getSReclaimableRegex.FindStringSubmatch(contents)
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
func getRawCpu() (string, error) {
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "", err
	}
	defer f.Close()

	raw, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	contents := string(raw)
	match := getCpuModelRegex.FindStringSubmatch(contents)
	return removeExtraSpacesRegex.ReplaceAllString(match[1], " "), nil
}

// Returns GPU manufacturer and model
func getRawGpu() (string, error) {
	out, err := exec.Command("lspci", "-mm").Output()
	if err != nil {
		return "", err
	}

	contents := string(out)
	match := getGpuManufacturerAndModel.FindStringSubmatch(contents)
	return fmt.Sprintf("%v %v", match[1], match[2]), nil
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
