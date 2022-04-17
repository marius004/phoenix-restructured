package main

import (
	"log"

	"github.com/marius004/phoenix/internal"
)

func main() {
	config := internal.NewConfig()
	db, err := internal.ConnectToPSQL(internal.GenerateDatabaseDSN(config))

	if err != nil {
		log.Fatalln("Could not connect to the database", err)
	}

	server := NewServer(db, config, &log.Logger{})
	server.Serve()
}
