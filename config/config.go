package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var ErrNotFoundConfig = errors.New("config not found")

var possibleOptions = map[string]bool{
	"logo":          false,
	"userline":      false,
	"userunderline": false,
	"kernel":        false,
	"uptime":        false,
	"shell":         false,
	"memory":        false,
}

var Config []string

func init() {
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	options := handleConfig(config)

	opts := []string{}
	for _, option := range options {
		_, exists := possibleOptions[option]
		if !exists {
			panic(fmt.Sprintf("unknown option %v", option))
		}

		opts = append(opts, option)
	}

	Config = opts
}

func loadConfig() (string, error) {
	configLocation := ""

	_, err := os.Stat(os.ExpandEnv("/home/$USER/.config/barkfetch/config"))
	if err == nil {
		configLocation = os.ExpandEnv("/home/$USER/.config/barkfetch/config")
		goto configChosed
	}

	_, err = os.Stat("/etc/barkfetch.config")
	if err == nil {
		configLocation = "/etc/barkfetch.config"
		goto configChosed
	}

	_, err = os.Stat("./barkfetch.config")
	if err == nil {
		configLocation = "./barkfetch.config"
		goto configChosed
	}

	return "", ErrNotFoundConfig

configChosed:
	file, err := os.Open(configLocation)
	if err != nil {
		return "", err
	}
	defer file.Close()

	config, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(config), nil
}

func handleConfig(config string) []string {
	options := []string{}
	for _, line := range strings.Split(config, "\n") {
		if len(line) == 0 {
			continue
		}

		if []rune(line)[0] == '#' {
			continue
		}

		options = append(options, strings.TrimSpace(line))
	}

	return options
}
