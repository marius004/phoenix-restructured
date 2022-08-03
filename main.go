package main

import (
	"log"

	"github.com/marius004/phoenix-/internal"
)

var evalConfigPath = "eval.config.json"

func main() {
	config := internal.NewConfig()
	evalConfig := internal.NewEvalConfig(evalConfigPath)

	db, err := internal.ConnectToPSQL(internal.GenerateDatabaseDSN(config))
	if err != nil {
		log.Fatalln("Could not connect to the database", err)
	}

	server := NewServer(db, config, evalConfig)
	server.Serve()
}
