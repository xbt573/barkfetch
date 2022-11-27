package info

import (
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

	return fmt.Sprintf("%v@%v\n", username, hostname), nil
}

func GetUserUnderline() (string, error) {
	userline, err := GetUserline()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v\n", strings.Repeat("=", len(userline))), nil
}

func GetLogo() string {
	return `  /\
 /  \
/\/\/\`
}

func GetLogoLines() int {
	return len(strings.Split(GetLogo(), "\n"))
}

func GetLogoMaxLength() int {
	logo := GetLogo()

	max := 0
	for _, line := range strings.Split(logo, "\n") {
		if len(line) > max {
			max = len(line)
		}
	}

	return max
}
