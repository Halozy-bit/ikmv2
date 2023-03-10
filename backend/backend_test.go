package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ikmv2/backend/api"
	"github.com/ikmv2/backend/config"
	asynctask "github.com/ikmv2/backend/pkg/async_task"
	"github.com/ikmv2/backend/pkg/helper"
	"github.com/ikmv2/backend/pkg/repository"
	"github.com/ikmv2/backend/pkg/sidejob"
	testhelper "github.com/ikmv2/backend/pkg/test_helper"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// NOTE
// the program is running properly
// but the logic test is still not correct

var repo repository.Repository
var dummyCtlg testhelper.CatalogDummy
var db *mongo.Database
var job sidejob.RefreshCatalogPage
var initData int = 40
var jobRunCount int

func TestMain(m *testing.M) {
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

	job = sidejob.RefreshCatalogPage{
		Db: db,
		TaskIdentifier: asynctask.TaskIdentifier{
			Name:     "refresh catalog page",
			Interval: time.Minute * 3,
		},
	}
	log.Print("sidejob initialized")

	log.Println("fill database arg dan set dummy data")
	repo = repository.NewRepository(db)
	dummyCtlg = testhelper.SeedCatalog(initData, repo)

	log.Println("Running tests!")

	exitVal := m.Run()

	db.Collection("catalog").DeleteMany(context.TODO(), bson.D{})
	db.Client().Disconnect(context.TODO())
	log.Println("wiping data")

	log.Println("exiting tests!")
	os.Exit(exitVal)
}

func TestGetPagination(t *testing.T) {
	defer func(t *testing.T) {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			t.Fail()
		}
	}(t)
	log.Println("load server")
	node := api.NewEndpoint(repo)
	job.Run()
	jobRunCount++

	maxPage := helper.MaxPage(helper.MaxProductPerPage, len(dummyCtlg.Dummy))
	server := node.Server()

	for i := 1; i <= maxPage; i++ {
		path := fmt.Sprintf("/catalog/%d", i)

		tp := testhelper.TestSetParam{Name: "page", Value: fmt.Sprint(i)}
		c, rec := testhelper.CreateRequestContext(server, path, nil, tp)

		t.Log("request address: ", c.Path())
		err := node.GetCatalog(c)
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

		testhelper.VerifyOutput(t, i, maxPage, ctlg, dummyCtlg.Dummy)
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

	job.Run()
	jobRunCount++
	job.Run()
	jobRunCount++

	maxPage := helper.MaxPage(helper.MaxProductPerPage, len(dummyCtlg.Dummy))

	log.Println("load server")
	node := api.NewEndpoint(repo)
	server := node.Server()
	var expectFirst int = 6
	var expectLast int = expectFirst

	for page := 1; page <= maxPage; page++ {
		path := fmt.Sprintf("/catalog/%d", page)

		tp := testhelper.TestSetParam{Name: "page", Value: fmt.Sprint(page)}
		c, rec := testhelper.CreateRequestContext(server, path, nil, tp)

		t.Log("request address: ", c.Path())
		err := node.GetCatalog(c)
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

		t.Log("get response length: ", len(js.Catalog))

		if len(js.Catalog) > helper.MaxProductPerPage {
			panic("too many return")
		}

		expectFirst, expectLast = testhelper.VerifyOutputNextID(
			testhelper.OutputNextID{
				T: t, Page: page, ExpectedFirst: expectFirst, ExpectedLast: expectLast,
				InitData: initData, Ctlg: js.Catalog, DummyData: dummyCtlg.Dummy,
			},
		)
	}
}
