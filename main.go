package main

import (
	"context"
	"log"

	"github.com/ikmv2/backend/api"

	"github.com/ikmv2/backend/pkg/repository"

	"github.com/ikmv2/backend/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := config.MongoConfig{}
	db, err := repository.ConnectDatabase(ctx, cfg)
	if err != nil {
		log.Fatalln(err)
	}

	repo := repository.NewRepository(db)

	node := api.NewEndpoint(repo)
	node.ExposeRoute()
	node.StartServer()
}
