package cmd

import (
	"barkfetch/info"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Config not found error
var ErrConfigNotFound = errors.New("not found config")

func loadConfig() (map[string]string, error) {
	f, err := os.Open("./barkfetch.config")
	if err == nil {
		goto configChosed
	}

	f, err = os.Open("~/.config/barkfetch")
	if err == nil {
		goto configChosed
	}

	f, err = os.Open("/etc/barkfetch.config")
	if err == nil {
		goto configChosed
	}

	return map[string]string{}, ErrConfigNotFound

configChosed:
	defer f.Close()

	raw, err := io.ReadAll(f)
	if err != nil {
		return map[string]string{}, err
	}

	contents := string(raw)
	return parseConfig(contents), nil
}

func parseConfig(config string) map[string]string {
	options := make(map[string]string)

	for _, line := range strings.Split(config, "\n") {
		if len(line) == 0 {
			continue
		}

		if line[0] == '#' {
			continue
		}

		kv := strings.Split(line, "=")

		if len(kv) < 2 || len(kv) > 2 {
			continue
		}

		options[kv[0]] = kv[1]
	}

	return options
}

func Run() error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	// "logo", "userline", "userunderline", "kernel", "uptime", "shell", "memory"
	sysinfo, err := info.GetInfoString(config)

	if err != nil {
		return err
	}

	fmt.Println(sysinfo)

	return nil
}
