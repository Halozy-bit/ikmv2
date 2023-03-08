package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/ikmv2/backend/api"
	"github.com/ikmv2/backend/config"
	"github.com/ikmv2/backend/pkg/cache"
	"github.com/ikmv2/backend/pkg/repository"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var repo repository.Repository
var conn *mongo.Client
var category [3]string
var sumPerCategory [3]int
var DummyData []primitive.ObjectID
var initData int = 40

func TestMain(m *testing.M) {

	for i := range category {
		category[i] = repository.RandName(true)
	}

	cfg := config.MongoConfig{
		MongoDriver: "mongodb",
		User:        "user",
		Password:    "secret",
		Address:     "127.0.0.1",
		DbName:      "ikm-project",
	}

	log.Println("connecting to database")
	db, err := repository.ConnectDatabase(cfg)
	if err != nil {
		log.Fatal()
	}
	log.Println("database connected")

	log.Println("fill database arg dan set dummy data")
	conn = db.Client()
	repo = repository.NewRepository(db)
	seedData()

	log.Println("Running tests!")

	exitVal := m.Run()

	db.Collection("catalog").DeleteMany(context.TODO(), bson.D{})
	conn.Disconnect(context.TODO())
	log.Println("wiping data")

	log.Println("exiting tests!")
	os.Exit(exitVal)
}

func seedData() {
	insertParam := make([]repository.DocCatalog, initData)
	for i := range insertParam {
		Category := repository.RandInt(0, 3)
		sumPerCategory[Category]++
		insertParam[i] = repository.DocCatalog{
			Name:        repository.RandName(),
			Category:    category[Category],
			Description: repository.RandString(15),
			Owner:       primitive.NewObjectID().Hex(),
			Foto:        primitive.NewObjectID().Hex() + ".jpg",
		}
	}

	for i := range insertParam {
		insrd, err := repo.Insert(context.TODO(), insertParam[i])
		if err != nil {
			continue
		}

		id := insrd.(primitive.ObjectID)

		primID := id
		log.Print(primID)
		DummyData = append(DummyData, primID)
	}
}

func TestGetPagination(t *testing.T) {
	defer func(t *testing.T) {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			t.Fail()
		}
	}(t)
	maxPage := len(DummyData) / api.MaxProductPerPage
	if maxPage > 0 {
		maxPage += 1
	}

	log.Println("load server")
	node := api.NewEndpoint(repo)
	server := node.Server()
	var last_id = ""
	for i := 1; i <= maxPage; i++ {
		path := fmt.Sprintf("/catalog/%d", i)
		body, err := EncodeID(last_id)
		if !assert.NoError(t, err) {
			break
		}

		tp := TestSetParam{Name: "page", Value: fmt.Sprint(i)}
		c, rec := CreateRequestContext(server, path, body, tp)

		t.Log("request address: ", c.Path())
		err = node.GetCatalog(c)
		if !assert.NoError(t, err) {
			t.Log(err)
			break
		}

		if !assert.Equal(t, http.StatusOK, rec.Code) {
			t.Log(rec.Body.String())
			break
		}

		var js = struct {
			Catalog []repository.DocCatalog `json:"catalog"`
		}{}

		err = json.Unmarshal(rec.Body.Bytes(), &js)

		if !assert.NoError(t, err) {
			break
		}

		ctlg := js.Catalog
		t.Log("get response length: ", len(ctlg))

		last_id = verifyOutput(t, i, maxPage, ctlg, DummyData)
	}
}

// TODO
// Pagination per category
// Find error case

func TestGetPaginationNextID(t *testing.T) {
	defer func(t *testing.T) {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			t.Fail()
		}
	}(t)

	top := 27
	cache.Store(cache.TopCatalog, DummyData[top].Hex())
	cache.Store(cache.BottomCatalog, DummyData[initData-1].Hex())

	maxPage := len(DummyData) / api.MaxProductPerPage
	if maxPage > 0 {
		maxPage += 1
	}

	log.Println("load server")
	node := api.NewEndpoint(repo)
	server := node.Server()
	var last_id = ""
	var trait = int(top)

	for page := 1; page <= maxPage; page++ {
		path := fmt.Sprintf("/catalog/%d", page)
		body, err := EncodeID(last_id)
		if !assert.NoError(t, err) {
			break
		}

		tp := TestSetParam{Name: "page", Value: fmt.Sprint(page)}
		c, rec := CreateRequestContext(server, path, body, tp)

		t.Log("request address: ", c.Path())
		err = node.GetCatalog(c)
		if !assert.NoError(t, err) {
			t.Log(err)
			break
		}

		if !assert.Equal(t, http.StatusOK, rec.Code) {
			t.Log(rec.Body.String())
			break
		}

		var js = struct {
			Catalog []repository.DocCatalog `json:"catalog"`
		}{}

		err = json.Unmarshal(rec.Body.Bytes(), &js)

		if !assert.NoError(t, err) {
			break
		}

		ctlg := js.Catalog
		t.Log("get response length: ", len(ctlg))

		if len(ctlg) > api.MaxProductPerPage {
			panic("too many return")
		}

		trait, last_id = verifyOutputNextID(
			outputNextID{
				t: t, page: page, maxPage: maxPage, trait: trait, initData: initData,
				last_id: last_id, ctlg: ctlg, DummyData: DummyData,
			},
		)
	}
}
