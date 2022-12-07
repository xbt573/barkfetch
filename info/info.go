package info

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Possible options, to make output sorted independent of config/cmd
var possibleOptions = []string{"logo", "userline", "userunderline", "os",
	"kernel", "uptime", "shell", "resolution", "cpu", "gpu", "memory", "colors"}

// Regexp matching empty lines, useful to make output more pretty
var emptyLinesRegex = regexp.MustCompile(`(?m)\n$`)

// Helper function, chains fmt.Sprintf and os.Expand(..., ColorExpand)
func formatAndColor(format string, args ...any) string {
	return os.Expand(
		fmt.Sprintf(format, args...),
		ColorExpand,
	)
}

// Returns processed info for pretty output
func GetInfoString(options map[string]string) string {
	// out string
	var output string

	// offset for printing labels
	var offset int

	// Info lines count, for calibrating newlines at the end
	var logolines, lines int

	for _, possibleOption := range possibleOptions {
		value, exists := options[possibleOption]
		if !exists {
			continue
		}

		if value == "false" {
			continue
		}

		switch possibleOption {
		case "logo":
			logo := getLogo(value)

			output += os.Expand(logo.Logo, ColorExpand) +
				strings.Repeat("\x1b[F", logo.Lines-1)
			offset = logo.MaxLength + 2
			Colors["caccent"] = Colors[logo.AccentColor]
			logolines = logo.Lines

		case "userline":
			username := getRawUser()
			hostname := getRawHostname()

			output += formatAndColor(
				"\x1b[%vG${caccent}%v${creset}@${caccent}%v${creset}\n",
				offset,
				username,
				hostname,
			)
			lines++

		case "userunderline":
			username := getRawUser()
			hostname := getRawHostname()

			output += formatAndColor(
				"\x1b[%vG%v\n",
				offset,
				strings.Repeat("-", len(username)+len(hostname)),
			)
			lines++

		case "os":
			os := getRawPrettyName()
			arch := getRawArchitecture()

			output += formatAndColor(
				"\x1b[%vG${caccent}OS${creset}: %v %v\n",
				offset,
				os,
				arch,
			)
			lines++

		case "kernel":
			kernel := getRawKernel()

			output += formatAndColor(
				"\x1b[%vG${caccent}Kernel${creset}: %v\n",
				offset,
				kernel,
			)
			lines++

		case "uptime":
			uptime := getRawUptime()

			if uptime <= 0 {
				output += formatAndColor(
					"\x1b[%vG${caccent}Uptime:${creset}: n/a\n",
					offset,
				)
			}

			if uptime > 0 && uptime <= 60 {
				output += formatAndColor(
					"\x1b[%vG${caccent}Uptime${creset}: %v s\n",
					offset,
					int(uptime),
				)
			}

			if uptime > 60 && uptime <= 3600 {
				output += formatAndColor(
					"\x1b[%vG${caccent}Uptime${creset}: %v m, %v s\n",
					offset,
					int(uptime/60),
					int(uptime%60),
				)
			}

			if uptime > 3600 && uptime <= 86400 {
				output += formatAndColor(
					"\x1b[%vG${caccent}Uptime${creset}: %v h, %v m, %v s\n",
					offset,
					int(uptime/3600),
					int((uptime%3600)/60),
					int((uptime%3600)%60),
				)
			}

			if uptime > 86400 {
				output += formatAndColor(
					"\x1b[%vG${caccent}Uptime${creset}: %v d, %v h, %v m, %v s\n",
					offset,
					int(uptime/86400),
					int((uptime%86400)/3600),
					int(((uptime%86400)%3600)/60),
					int(((uptime%86400)%3600)%60),
				)
			}

			lines++

		case "shell":
			shell := getRawShell()

			output += formatAndColor(
				"\x1b[%vG${caccent}Shell${creset}: %v\n",
				offset,
				shell,
			)
			lines++

		case "resolution":
			resolution := getRawScreenResolution()

			output += formatAndColor(
				"\x1b[%vG${caccent}Resolution${creset}: %v\n",
				offset,
				resolution,
			)
			lines++

		case "cpu":
			cpu := getRawCpu()

			output += formatAndColor(
				"\x1b[%vG${caccent}CPU${creset}: %v\n",
				offset,
				cpu,
			)
			lines++

		case "gpu":
			gpus := getRawGpus()

			if len(gpus) == 0 {
				output += formatAndColor(
					"\x1b[%vG${caccent}GPU${creset}: n/a\n",
					offset,
				)
			}

			for _, gpu := range gpus {
				output += formatAndColor(
					"\x1b[%vG${caccent}GPU${creset}: %v\n",
					offset,
					gpu,
				)
				lines++
			}

		case "memory":
			used, total := getRawMemory()

			if used <= 0 || total <= 0 {
				output += formatAndColor(
					"\x1b[%vG${caccent}Memory${creset}: n/a\n",
					offset,
				)
			} else {
				output += formatAndColor(
					"\x1b[%vG${caccent}Memory${creset}: %v / %v Mb (%v%%)\n",
					offset,
					used,
					total,
					int(float32(used)/float32(total)*100.0),
				)
			}

			lines++

		case "colors":
			colors := getRawColors()
			i := 0

			output += fmt.Sprintf("\x1b[%vG", offset)

			for _, color := range colors {
				if i == 8 {
					output += fmt.Sprintf("\n\x1b[%vG", offset)
					lines++
					i = 0
				}

				output += formatAndColor("%v", color)
				i++
			}
			lines++
		}
	}

	output = emptyLinesRegex.ReplaceAllString(output, "")

	if lines < logolines {
		output += strings.Repeat("\n", logolines-lines)
	}

	return output
}
