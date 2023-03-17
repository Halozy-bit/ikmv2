package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ikmv2/backend/config"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var repo Repository
var inserted []Product
var conn *mongo.Database
var category [3]string
var sumPerCategory [3]int

func TestMain(m *testing.M) {
	log.Println("Running tests!")

	for i := range category {
		category[i] = RandName(true)
	}

	exitVal := m.Run()
	log.Println("exiting tests!")

	os.Exit(exitVal)
}

func TestConnection(t *testing.T) {
	cfg := config.MongoConfig{
		Driver:   mongoNative,
		User:     "user",
		Password: "secret",
		Address:  "127.0.0.1",
		DbName:   "ikm-project",
	}

	db, err := ConnectDatabase(cfg)
	if !assert.NoError(t, err) {
		t.Fatal()
	}

	conn = db
	repo = NewRepository(conn)
}

func TestInsert(t *testing.T) {
	insertParam := make([]Product, 10)
	for i := range insertParam {
		Category := RandInt(0, 3)
		sumPerCategory[Category]++
		insertParam[i] = Product{
			Name:        RandName(),
			Category:    category[Category],
			Description: RandString(15),
			Owner:       primitive.NewObjectID().Hex(),
			Foto: Foto{
				Cover:   primitive.NewObjectID().Hex() + ".jpg",
				Detail1: primitive.NewObjectID().Hex() + ".jpg",
				Detail2: primitive.NewObjectID().Hex() + ".jpg",
			},
			Weight: []string{
				fmt.Sprintf("%d gr", RandInt(100, 500)),
				fmt.Sprintf("%d gr", RandInt(100, 500)),
			},
			Variant: []string{
				RandName(true),
				RandName(true),
			},
			Composition: []string{
				RandName(true),
				RandName(true),
			},
		}
	}

	for i, val := range insertParam {
		id, err := repo.InsertCatalog(context.TODO(), val)
		assert.NoError(t, err)
		if err != nil {
			continue
		}

		primID := id.(primitive.ObjectID)
		insertParam[i].Id = primID

		t.Log("inserted: ", insertParam[i])
		inserted = append(inserted, insertParam[i])
	}
}

func TestTopAndBottom(t *testing.T) {
	first, err := repo.FirstItem()
	assert.NoError(t, err)
	last, err := repo.LastItem()
	assert.NoError(t, err)

	assert.NotEmpty(t, first)
	assert.NotEmpty(t, last)

	switch len(inserted) {
	case 1:
		assert.Equal(t, first.Id, last.Id)
	default:
		assert.NotEqual(t, first.Id, last.Id)
	}
}

func TestFirstPage(t *testing.T) {
	firstPage, err := repo.CatalogFirstLine(context.TODO(), int64(len(inserted)))
	assert.NoError(t, err)
	assert.Equal(t, len(firstPage), len(inserted))

	firstPage, err = repo.CatalogFirstLine(context.TODO(), int64(len(inserted)-2))
	assert.NoError(t, err)
	assert.Equal(t, len(firstPage), len(inserted)-2)
}

func TestGreater(t *testing.T) {
	gt, err := repo.CatalogGtId(context.TODO(), inserted[0].Id, 2)
	assert.NoError(t, err)
	assert.NotEmpty(t, gt)
	assert.Equal(t, 2, len(gt))
	assert.NotEqual(t, gt[0].Id, inserted[0].Id)

	gte, err := repo.CatalogGteId(context.TODO(), inserted[0].Id, 2)
	assert.NoError(t, err)
	assert.NotEmpty(t, gt)
	assert.Equal(t, 2, len(gt))
	assert.Equal(t, gte[0].Id, inserted[0].Id)
}

func TestDelete(t *testing.T) {
	defer conn.Client().Disconnect(context.TODO())

	for i := range category {
		t.Log(category[i], ": ", sumPerCategory[i])
	}

	c, err := repo.CountCatalog(context.TODO())
	count := int(c)
	assert.NoError(t, err)
	assert.Equal(t, len(inserted), count)

	reslt, err := conn.Collection("catalog").DeleteMany(context.TODO(), bson.D{})
	assert.NoError(t, err)
	t.Log("deleted count: ", reslt.DeletedCount)
	assert.Equal(t, len(inserted), int(reslt.DeletedCount))

	c, err = repo.CountCatalog(context.TODO())
	assert.NoError(t, err)
	assert.Zero(t, c)
}
