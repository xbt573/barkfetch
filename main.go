package main

import "barkfetch/cmd"

func main() {
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
