package logos

import (
	_ "embed"
	"regexp"
	"strings"
)

// Logo section
var (
	emptyRegex = regexp.MustCompile(`\$\{c(\d{1,2}|reset)\}`)
	logoRegex  = regexp.MustCompile(`#accent (c\d{1,2})\n((.|\n)*)`)
	//go:embed default.txt
	_default string
)

func GuessLogo() (Logo, error) {
	// Currently not supported, returns default
	return LogoByName("default"), nil
}

func LogoByName(name string) Logo {
	var logo Logo

	switch name {
	default:
		match := logoRegex.FindStringSubmatch(_default)
		logo.Logo = match[2]
		logo.Lines = getLines(logo.Logo)
		logo.MaxLength = getMaxLength(logo.Logo)
		logo.AccentColor = match[1]
	}

	return logo
}

func getLines(logo string) int {
	emptyLogo := emptyRegex.ReplaceAllString(logo, "")
	return len(strings.Split(emptyLogo, "\n"))
}

func getMaxLength(logo string) int {
	emptyLogo := emptyRegex.ReplaceAllString(logo, "")
	max := 0
	for _, line := range strings.Split(emptyLogo, "\n") {
		if len(line) > max {
			max = len(line)
		}
	}

	return max
}
