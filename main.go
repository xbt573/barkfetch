package main

import "github.com/xbt573/barkfetch/cmd"

func main() {
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
