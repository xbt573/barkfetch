package info

import (
	_ "embed"
	"os"
	"regexp"
	"strings"
)

// Useful regexes
var (
	replaceDirectivesRegexp = regexp.MustCompile(`(?:#accent c\d{1,2}\n|\$\{c(?:\d{1,2}|reset)\})`)
	getAccentRegexp         = regexp.MustCompile(`#accent (c\d{1,2})`)
	getLogoRegexp           = regexp.MustCompile(`#accent c\d{1,2}\n([\S\s]*)`)
)

// Gets username from program environment
func getRawUser() string {
	return os.Getenv("USER")
}

// Gets system hostname
func getRawHostname() (string, error) {
	return os.Hostname()
}

var (
	//go:embed logos/default.txt
	_default string

	//go:embed logos/void.txt
	_void string

	//go:embed logos/gentoo.txt
	_gentoo string
)

// Returns distro logo by name, or guesses if arg is "auto"
// Currently return small count of logos, mostly default
func getLogo(distro string) (Logo, error) {
	logoText := ""

	switch distro {
	case "auto":
		return getLogo(guessDistro())

	case "void":
		logoText = _void

	case "gentoo":
		logoText = _gentoo

	default:
		logoText = _default
	}

	var logo Logo

	match := getLogoRegexp.FindStringSubmatch(logoText)
	logo.Logo = match[1]

	match = getAccentRegexp.FindStringSubmatch(logoText)
	logo.AccentColor = match[1]

	directiveFreeLogoText := replaceDirectivesRegexp.ReplaceAllString(logoText, "")
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
