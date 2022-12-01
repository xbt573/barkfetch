package main

import (
	"barkfetch/config"
	"barkfetch/logos"
	"barkfetch/routines"
	"fmt"
	"os"
	"strings"
)

func main() {
	var logo logos.Logo
	var lines int
	for _, opt := range config.Config {
		something := routines.Routines[opt]()

		switch obj := something.(type) {
		case logos.Logo:
			logo = obj
			logoText := os.Expand(logo.Logo, logos.ColorExpand)

			fmt.Printf("%v%v", logoText, strings.Repeat("\x1b[F", logo.Lines-1))
		default:
			fmt.Printf("%v%v", fmt.Sprintf("\x1b[%vG", logo.MaxLength+2), obj.(string))
			lines++
		}
	}

	if lines < logo.Lines {
		fmt.Println(strings.Repeat("\n", logo.Lines-lines))
	}
}
