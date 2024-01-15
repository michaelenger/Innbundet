package main

import (
	"github.com/michaelenger/innbundet/db"
	"github.com/michaelenger/innbundet/server"
	"log"
)

func main() {
	logger := log.Default()

	// Database
	err := db.Init()
	if err != nil {
		panic(err)
	}

	// Server
	serv, err := server.Init()
	if err != nil {
		panic(err)
	}

	logger.Fatal(serv.Start(":8080"))
}
