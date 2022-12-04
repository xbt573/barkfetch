package info

import (
	_ "embed"
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

	//go:embed logos/void.txt
	_void string

	//go:embed logos/gentoo.txt
	_gentoo string

	//go:embed logos/arch.txt
	_arch string

	//go:embed logos/linuxmint.txt
	_linuxmint string

	//go:embed logos/manjaro.txt
	_manjaro string

	//go:embed logos/mxlinux.txt
	_mxlinux string

	//go:embed logos/nixos.txt
	_nixos string

	//go:embed logos/opensuse.txt
	_opensuse string

	//go:embed logos/parabola.txt
	_parabola string

	//go:embed logos/popos.txt
	_popos string

	//go:embed logos/postmarketos.txt
	_postmarketos string

	//go:embed logos/pureos.txt
	_pureos string

	//go:embed logos/slackware.txt
	_slackware string

	//go:embed logos/ubuntu.txt
	_ubuntu string
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

	case "arch":
		logoText = _arch

	case "linuxmint":
		logoText = _linuxmint

	case "manjaro":
		logoText = _manjaro

	case "mxlinux":
		logoText = _mxlinux

	case "nixos":
		logoText = _nixos

	case "opensuse":
		logoText = _opensuse

	case "suse":
		logoText = _opensuse

	case "parabola":
		logoText = _parabola

	case "pop-os":
		logoText = _popos

	case "popos":
		logoText = _popos

	case "postmarketos":
		logoText = _postmarketos

	case "pureos":
		logoText = _pureos

	case "slackware":
		logoText = _slackware

	case "ubuntu":
		logoText = _ubuntu

	default:
		logoText = _default
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
