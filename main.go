package main

import (
	"github.com/michaelenger/innbundet/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
