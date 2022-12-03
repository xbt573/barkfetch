package cmd

import (
	"barkfetch/info"
	"fmt"
)

func Run() error {
	// "logo", "userline", "userunderline", "kernel", "uptime", "shell", "memory"
	sysinfo, err := info.GetInfoString(map[string]string{
		"logo":          "auto",
		"userline":      "true",
		"userunderline": "true",
		"kernel":        "true",
		"uptime":        "true",
		"shell":         "true",
		"memory":        "true",
	})

	if err != nil {
		return err
	}

	fmt.Println(sysinfo)

	return nil
}
