package info

import (
	"barkfetch/logos"
	"fmt"
	"os"
	"strings"
)

func getUsername() string {
	return os.Getenv("USER")
}

func getHostname() (string, error) {
	return os.Hostname()
}

func GetUserline() (string, error) {
	username := getUsername()
	hostname, err := getHostname()
	if err != nil {
		return "", err
	}

	status := fmt.Sprintf(
		"${caccent}%v${creset}@${caccent}%v${creset}\n",
		username,
		hostname,
	)

	colored := os.Expand(status, logos.ColorExpand)

	return colored, nil
}

func GetUserUnderline() (string, error) {
	username := getUsername()
	hostname, err := getHostname()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v\n", strings.Repeat("=", len(username)+len(hostname))), nil
}

func GetLogo() (logos.Logo, error) {
	logo, err := logos.GuessLogo()
	if err != nil {
		return logos.Logo{}, err
	}

	return logo, nil
}
