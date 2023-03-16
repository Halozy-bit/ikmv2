package testhelper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/ikmv2/backend/pkg/helper"
	"github.com/ikmv2/backend/pkg/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CatalogDummy struct {
	Category     []string
	CounCategory []int
	Dummy        []primitive.ObjectID
}

func SeedCatalog(total int, repo repository.Repository) CatalogDummy {
	cd := CatalogDummy{
		Category:     helper.CategoryAvail,
		CounCategory: make([]int, len(helper.CategoryAvail)),
	}

	cTotal := int32(len(helper.CategoryAvail))
	for i := 0; i < total; i++ {
		Category := repository.RandInt(0, cTotal)
		param := repository.Product{
			Name:        repository.RandName(),
			Category:    cd.Category[Category],
			Description: repository.RandString(15),
			Owner:       primitive.NewObjectID().Hex(),
			Foto: repository.Foto{
				Cover:   primitive.NewObjectID().Hex() + ".jpg",
				Detail1: primitive.NewObjectID().Hex() + ".jpg",
				Detail2: primitive.NewObjectID().Hex() + ".jpg",
			},
			Weight: []string{
				fmt.Sprintf("%d gr", repository.RandInt(100, 500)),
				fmt.Sprintf("%d gr", repository.RandInt(100, 500)),
			},
			Variant: []string{
				repository.RandName(true),
				repository.RandName(true),
			},
			Composition: []string{
				repository.RandName(true),
				repository.RandName(true),
			},
		}
		insrd, err := repo.InsertCatalog(context.TODO(), repository.ProductToDocument(param))
		if err != nil {
			panic(err)
		}

		cd.CounCategory[Category]++
		id := insrd.(primitive.ObjectID)
		log.Print(id)
		cd.Dummy = append(cd.Dummy, id)
	}

	return cd
}

type TestSetParam struct {
	Name  string
	Value string
}

func CreateRequestContext(server *echo.Echo, reqPath string, body io.Reader, params ...TestSetParam) (echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest("GET", reqPath, body)
	req.Header["Content-Type"] = append(req.Header["Content-Type"], "application/json")
	rec := httptest.NewRecorder()

	c := server.NewContext(req, rec)
	c.SetPath(req.URL.Path)
	for _, param := range params {
		c.SetParamNames(param.Name)
		c.SetParamValues(param.Value)
	}
	return c, rec
}

func EncodeID(last_id string) (io.Reader, error) {
	if last_id != "" {
		var lId = struct {
			LastID string `json:"last_id"`
		}{LastID: last_id}

		mr, err := json.Marshal(lId)
		if err != nil {
			return nil, err
		}

		return bytes.NewReader(mr), nil
	}
	return nil, nil
}

// return -1 if error
func VerifyOutput(t *testing.T, page int, maxPage int, ctlg []repository.CatalogDisplay, DummyData []primitive.ObjectID) {
	if len(ctlg) > helper.MaxProductPerPage {
		t.Fatal()
		return
	}

	lastIdx := len(ctlg) - 1
	trait := (page - 1) * helper.MaxProductPerPage
	t.Log("first index: ", ctlg[0].Id, ", last index: ", ctlg[lastIdx].Id)
	assert.Equal(t, ctlg[0].Id, DummyData[trait])

	trait += lastIdx

	log.Print("trait: ", trait)
	assert.Equal(t, ctlg[lastIdx].Id, DummyData[trait])

}

type OutputNextID struct {
	T             *testing.T
	Page          int
	ExpectedFirst int
	ExpectedLast  int
	InitData      int
	Ctlg          []repository.CatalogDisplay
	DummyData     []primitive.ObjectID
}

// return trait and last_id
func VerifyOutputNextID(oid OutputNextID) (expectFirst int, expectLast int) {
	expectFirst, expectLast = oid.ExpectedFirst, oid.ExpectedLast
	lastIdx := len(oid.Ctlg) - 1
	log.Print("first index: ", oid.Ctlg[0].Id, ", last index: ", oid.Ctlg[lastIdx].Id)
	log.Print("expect first: ", oid.DummyData[expectFirst])
	assert.Equal(oid.T, oid.Ctlg[0].Id, oid.DummyData[expectFirst])

	dummyLeng := len(oid.DummyData)
	expectLast += len(oid.Ctlg)
	if expectLast > dummyLeng {
		expectLast -= dummyLeng
	}

	log.Print("expect last: ", oid.DummyData[expectLast])

	assert.Equal(oid.T, oid.Ctlg[lastIdx].Id, oid.DummyData[expectLast])
	expectFirst += len(oid.Ctlg)

	if expectFirst > dummyLeng {
		expectFirst -= dummyLeng
	}

	return
}
