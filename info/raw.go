package info

import (
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

// Returns distro logo by name, or guesses if arg is "auto"
// Currently return default logo, unimplemented
func getLogo(distro string) (Logo, error) {
	text := `#accent c2
${c2}  /\
 /  \
/\/\/\${creset}`
	directiveFreeText := replaceDirectivesRegexp.ReplaceAllString(text, "")

	var logo Logo

	match := getLogoRegexp.FindStringSubmatch(text)
	logo.Logo = match[1]
	logo.Lines = len(strings.Split(directiveFreeText, "\n"))

	max := 0
	for _, line := range strings.Split(directiveFreeText, "\n") {
		if len(line) > max {
			max = len(line)
		}
	}

	logo.MaxLength = max

	match = getAccentRegexp.FindStringSubmatch(text)
	logo.AccentColor = match[1]

	return logo, nil
}
