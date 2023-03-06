package main

import (
	"log"

	"github.com/ikmv2/backend/api"

	"github.com/ikmv2/backend/pkg/repository"

	"github.com/ikmv2/backend/config"
)

func main() {
	cfg := config.MongoConfig{}

	db, err := repository.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	repo := repository.NewRepository(db)

	node := api.NewEndpoint(repo)
	node.ExposeRoute()
	node.StartServer()
}
