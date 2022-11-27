package main

import (
	"barkfetch/config"
	"barkfetch/console"
	"barkfetch/info"
	"barkfetch/routines"
	"barkfetch/utils"
)

func main() {
	_, width, err := utils.GetTerminalSize()
	if err != nil {
		panic(err)
	}

	var height int
	var logoEnabled bool

	for _, opt := range config.Config {
		if opt == "logo" {
			logoEnabled = true
			break
		}
	}

	if !logoEnabled {
		height = len(config.Config)
	} else {
		height = info.GetLogoLines()
		if len(config.Config) > height {
			height = len(config.Config)
		}
	}

	c := console.NewConsole(width, height)
	c.Clear()

	for _, opt := range config.Config {
		routines.Routines[opt](&c)
	}

	c.Render()
}
