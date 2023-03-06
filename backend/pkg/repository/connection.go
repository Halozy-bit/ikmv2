package repository

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ikmv2/backend/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoBuilder struct {
	url    string
	dbName string
	opt    *options.ServerAPIOptions
}

const mongoNative = "mongodb"
const mongoAtlasFreeTier = "mongodb+srv"

func buildConfig(mongoCfg config.MongoConfig) (mongoBuilder, error) {
	log.Println("building configuration")
	builder := mongoBuilder{}

	driver := strings.ToLower(mongoCfg.MongoDriver)

	switch driver {
	case mongoAtlasFreeTier:
		builder.url = "/?retryWrites=true&w=majority"
		// builder.opt = options.ServerAPI(options.ServerAPIVersion1)
	case mongoNative:
		builder.url = ":27017"
	default:
		return builder, fmt.Errorf("")
	}

	url := fmt.Sprintf("%s://%s:%s@%s", driver, mongoCfg.User, mongoCfg.Password, mongoCfg.Address)
	builder.url = url + builder.url
	builder.dbName = mongoCfg.DbName

	return builder, nil
}

func ConnectDatabase(mongoCfg config.MongoConfig) (*mongo.Database, error) {
	bl, err := buildConfig(mongoCfg)
	if err != nil {
		return nil, err
	}

	log.Println(bl.url)
	log.Println("setup environment database")
	var clientOptions = options.Client()
	clientOptions.ApplyURI(bl.url)
	clientOptions.SetServerAPIOptions(bl.opt)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("connecting to database")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	log.Println("database connected")

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	connection := client.Database(bl.dbName)

	return connection, nil
}
