package info

import (
	"embed"
	_ "embed"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
)

// Useful regexes
var (
	replaceDirectivesRegex = regexp.MustCompile(`(?:#accent c\d{1,2}\n|\$\{c(?:\d{1,2}|reset)\})`)
	getAccentRegex         = regexp.MustCompile(`#accent (c\d{1,2})`)
	getLogoRegex           = regexp.MustCompile(`#accent c\d{1,2}\n([\S\s]*)`)
)

// Gets username from program environment
func getRawUser() string {
	return os.Getenv("USER")
}

// Gets system hostname
func getRawHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "n/a"
	}

	return hostname
}

// Get local ip
func getRawLocalIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "n/a"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// Get outbound ip
func getRawOutboundIp() string {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "n/a"
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "n/a"
	}

	return string(ip)
}

// Gets OS architecture
func getRawArchitecture() string {
	return runtime.GOARCH
}

// Get colors table
func getRawColors() []string {
	colors := []string{
		"0", "1", "2", "3", "4", "5", "6", "7",
		"8", "9", "10", "11", "12", "13", "14", "15",
	}

	colorArr := []string{}

	for _, color := range colors {
		colorArr = append(colorArr, fmt.Sprintf("${c%v}███", color))
	}

	return colorArr
}

// Logos
var (
	//go:embed logos/default.txt
	_default string

	//go:embed logos/*
	logos embed.FS
)

// Returns distro logo by name, or guesses if arg is "auto"
// Currently return small count of logos, mostly default
func getLogo(distro string) Logo {
	logoText := _default

	if distro == "auto" {
		return getLogo(guessDistro())
	}

	bytes, err := logos.ReadFile(fmt.Sprintf("logos/%v.txt", distro))
	if err == nil {
		logoText = string(bytes)
	}

	var logo Logo

	logo.Logo = getLogoRegex.FindStringSubmatch(logoText)[1]

	logo.AccentColor = getAccentRegex.FindStringSubmatch(logoText)[1]

	directiveFreeLogoText := replaceDirectivesRegex.ReplaceAllString(logoText, "")
	logo.Lines = len(strings.Split(directiveFreeLogoText, "\n"))

	max := 0

	for _, line := range strings.Split(directiveFreeLogoText, "\n") {
		if len(line) > max {
			max = len(line)
		}
	}

	logo.MaxLength = max

	return logo
}
