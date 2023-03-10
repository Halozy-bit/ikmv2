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
	"github.com/ikmv2/backend/pkg/helper"
	"github.com/ikmv2/backend/pkg/repository"
	testhelper "github.com/ikmv2/backend/pkg/test_helper"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database
var ctldDummy testhelper.CatalogDummy

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

	ctldDummy = testhelper.SeedCatalog(40, repository.NewRepository(db))
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

	var before primitive.ObjectID
	for i := 0; i < 5; i++ {
		rcp.Run()
		run1 := cache.Pagination.CategoryPage(helper.CategoryAvail[0], 1)
		log.Println(run1)
		before = run1
		assert.NotEqual(t, before, run1)
	}
}

func TestAsyncTask1(t *testing.T) {
	defer func(t *testing.T) {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			t.Fail()
		}
	}(t)
	err := asynctask.AddTask(&RefreshCatalogPage{
		Db: db,
		TaskIdentifier: asynctask.TaskIdentifier{
			Name:     "refresh catalog page",
			Interval: time.Second * 2,
		},
	})

	assert.NoError(t, err)

	log.Println("preparing side job")
	err = asynctask.Start(time.Second * 2)
	assert.NoError(t, err)

	time.Sleep(time.Second * 20)
	log.Print(cache.Pagination.Page(1))

	log.Println("sending stop signal")
	asynctask.Stop()
}

func TestTask2(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			t.Fail()
		}
	}()

	rcp := RefreshCatalogCategoryPage{
		Db: db,
		TaskIdentifier: asynctask.TaskIdentifier{
			Name:     "refresh catalog page per category",
			Interval: time.Second,
		},
	}

	log.Println(ctldDummy.CounCategory[0])
	log.Println(ctldDummy.CounCategory[1])

	var before primitive.ObjectID
	for i := 0; i < 5; i++ {
		rcp.Run()
		run1 := cache.Pagination.CategoryPage(helper.CategoryAvail[0], 1)
		log.Println(run1)
		before = run1
		assert.NotEqual(t, before, run1)
	}
}
