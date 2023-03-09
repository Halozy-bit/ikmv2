package main

import (
	"log"

	"github.com/ikmv2/backend/api"

	"github.com/ikmv2/backend/pkg/repository"

	"github.com/ikmv2/backend/config"
)

func main() {
	cfg := config.MongoConfig{
		MongoDriver: "mongodb",
		User:        "user",
		Password:    "secret",
		Address:     "127.0.0.1",
		DbName:      "ikm-project",
	}

	db, err := repository.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	repo := repository.NewRepository(db)

	node := api.NewEndpoint(repo)
	node.StartSideJob(db)
	node.ExposeRoute()
	node.StartServer(":8081")
}
