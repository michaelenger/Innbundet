package main

import (
	"github.com/michaelenger/innbundet/db"
	"github.com/michaelenger/innbundet/server"
)

func main() {
	err := db.Init()
	if err != nil {
		panic(err)
	}

	server.RunServer()
}
