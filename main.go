package main

import (
	"log"

	"github.com/ikmv2/backend/api"
	"github.com/ikmv2/backend/config"
	"github.com/ikmv2/backend/pkg/repository"
)

func main() {
	cfg := config.MongoConfig{
		Driver:   "mongodb",
		User:     "user",
		Password: "secret",
		Address:  "127.0.0.1",
		DbName:   "ikm-project-test",
	}

	config.AutoEnv(&cfg)

	db, err := repository.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	repo := repository.NewRepository(db)

	node := api.NewEndpoint(repo)
	node.StartSideJob(db)
	node.ExposeRoute()
	node.StartServer(":8080")
}
