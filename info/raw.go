package info

import (
	"embed"
	_ "embed"
	"fmt"
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
func getRawHostname() (string, error) {
	return os.Hostname()
}

// Gets OS architecture
func getRawArchitecture() string {
	return runtime.GOARCH
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
func getLogo(distro string) (Logo, error) {
	logoText := ""

	if distro == "auto" {
		return getLogo(guessDistro())
	}

	bytes, err := logos.ReadFile(fmt.Sprintf("logos/%v.txt", distro))
	if err != nil {
		logoText = _default
	} else {
		logoText = string(bytes)
	}

	var logo Logo

	match := getLogoRegex.FindStringSubmatch(logoText)
	logo.Logo = match[1]

	match = getAccentRegex.FindStringSubmatch(logoText)
	logo.AccentColor = match[1]

	directiveFreeLogoText := replaceDirectivesRegex.ReplaceAllString(logoText, "")
	logo.Lines = len(strings.Split(directiveFreeLogoText, "\n"))

	max := 0

	for _, line := range strings.Split(directiveFreeLogoText, "\n") {
		if len(line) > max {
			max = len(line)
		}
	}

	logo.MaxLength = max

	return logo, nil
}
