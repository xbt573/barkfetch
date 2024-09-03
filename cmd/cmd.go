package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/xbt573/barkfetch/info"
)

// Config not found error
var ErrConfigNotFound = errors.New("not found config")

// Command-line arguments
var (
	_logo          = flag.String("logo", "auto", "Selects which logo is displayed")
	_userline      = flag.Bool("userline", true, "Display username and hostname")
	_userunderline = flag.Bool("userunderline", true, "Display fancy line of - under userline")
	_os            = flag.Bool("os", true, "Display host os and architecture")
	_kernel        = flag.Bool("kernel", true, "Display system kernel type and version")
	_uptime        = flag.Bool("uptime", true, "Display system uptime")
	_shell         = flag.Bool("shell", true, "Display current shell")
	_resolution    = flag.Bool("resolution", true, "Display screen resolution")
	_cpu           = flag.Bool("cpu", true, "Display CPU model")
	_gpu           = flag.Bool("gpu", true, "Display GPU manufacturer and model")
	_memory        = flag.Bool("memory", true, "Display used and total memory in megabytes")
	_localip       = flag.Bool("localip", true, "Display local IP")
	_remoteip      = flag.Bool("remoteip", true, "Display remote IP")
	_colors        = flag.Bool("colors", true, "Display colors")
)

var configCommonPaths = [3]string{"./barkfetch.config", os.ExpandEnv("$HOME/.config/barkfetch/config"), "/etc/barkfetch.config"}

// Helper function, returns true if flag was given at command-line
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		found = f.Name == name
		if found {
			return
		}
	})

	return found
}

// Converts true to "true" and false to "false"
func boolToString(input bool) string {
	return fmt.Sprintf("%t", input)
}

// Load config and return map of options
func loadConfig() (map[string]string, error) {
	var content []byte
	for _, configPath := range configCommonPaths {
		f, err := os.ReadFile(configPath)
		if err == nil {
			content = f
			goto configChoosed
		}
	}

	return map[string]string{}, ErrConfigNotFound

configChoosed:
	contents := string(content)
	config := parseConfig(contents)

	if isFlagPassed("logo") {
		config["logo"] = *_logo
	}

	if isFlagPassed("userline") {
		config["userline"] = boolToString(*_userline)
	}

	if isFlagPassed("userunderline") {
		config["userunderline"] = boolToString(*_userunderline)
	}

	if isFlagPassed("os") {
		config["os"] = boolToString(*_os)
	}

	if isFlagPassed("kernel") {
		config["kernel"] = boolToString(*_kernel)
	}

	if isFlagPassed("uptime") {
		config["uptime"] = boolToString(*_uptime)
	}

	if isFlagPassed("shell") {
		config["shell"] = boolToString(*_shell)
	}

	if isFlagPassed("resolution") {
		config["resolution"] = boolToString(*_resolution)
	}

	if isFlagPassed("cpu") {
		config["cpu"] = boolToString(*_cpu)
	}

	if isFlagPassed("gpu") {
		config["gpu"] = boolToString(*_gpu)
	}

	if isFlagPassed("memory") {
		config["memory"] = boolToString(*_memory)
	}

	if isFlagPassed("localip") {
		config["localip"] = boolToString(*_localip)
	}

	if isFlagPassed("remoteip") {
		config["remoteip"] = boolToString(*_remoteip)
	}

	if isFlagPassed("colors") {
		config["colors"] = boolToString(*_colors)
	}

	return config, nil
}

// Parse simple config
func parseConfig(config string) map[string]string {
	options := make(map[string]string)

	for _, line := range strings.Split(config, "\n") {
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		kv := strings.Split(line, "=")

		if len(kv) != 2 {
			continue
		}

		options[kv[0]] = kv[1]
	}

	return options
}

// Run cmd-related stuff and return non-nil error if something is wrong
func Run() error {
	flag.Parse()
	config, err := loadConfig()
	if err != nil {
		return err
	}

	sysinfo := info.GetInfoString(config)
	fmt.Println(sysinfo)

	return nil
}
