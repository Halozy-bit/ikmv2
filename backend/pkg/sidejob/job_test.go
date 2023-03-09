package sidejob

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ikmv2/backend/config"
	asynctask "github.com/ikmv2/backend/pkg/async_task"
	"github.com/ikmv2/backend/pkg/cache"
	"github.com/ikmv2/backend/pkg/repository"
	testhelper "github.com/ikmv2/backend/pkg/test_helper"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database

func TestMain(m *testing.M) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	cfg := config.MongoConfig{
		MongoDriver: "mongodb",
		User:        "user",
		Password:    "secret",
		Address:     "127.0.0.1",
		DbName:      "ikm-project",
	}
	var err error

	db, err = repository.ConnectDatabase(cfg)
	if err != nil {
		log.Fatal()
	}

	log.Println("database connected")

	testhelper.SeedCatalog(40, 3, repository.NewRepository(db))
	log.Println("Running tests!")

	exitVal := m.Run()

	log.Println("exiting tests!")
	log.Println("wiping data")

	db.Collection("catalog").DeleteMany(context.TODO(), bson.D{})
	db.Client().Disconnect(context.TODO())
	os.Exit(exitVal)
}

func TestTask1(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			t.Fail()
		}
	}()

	rcp := RefreshCatalogPage{
		Db: db,
		TaskIdentifier: asynctask.TaskIdentifier{
			Name:     "refresh catalog page",
			Interval: time.Second,
		},
	}

	rcp.Run()
	run1 := cache.Pagination.Page(1)
	t.Log(run1)

	rcp.Run()
	run2 := cache.Pagination.Page(1)
	t.Log(run2)

	assert.NotEqual(t, run1, run2)
}

func TestAsync(t *testing.T) {
	err := asynctask.AddTask(&RefreshCatalogPage{
		Db: db,
		TaskIdentifier: asynctask.TaskIdentifier{
			Name:     "refresh catalog page",
			Interval: time.Minute * 2,
		},
	})
	assert.NoError(t, err)
	log.Println("preparing side job")
	err = asynctask.Start(time.Second * 5)
	assert.NoError(t, err)
	time.Sleep(time.Minute * 3)
	log.Print(cache.Pagination.Page(1))
	log.Println("sending stop signal")
	time.Sleep(time.Minute * 3)
	asynctask.Stop()
}
