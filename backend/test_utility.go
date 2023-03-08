package backend

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/ikmv2/backend/api"
	"github.com/ikmv2/backend/pkg/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
func verifyOutput(t *testing.T, page int, maxPage int, ctlg []repository.DocCatalog, DummyData []primitive.ObjectID) string {
	if len(ctlg) > api.MaxProductPerPage {
		t.Fatal()
		return ""
	}

	var last_id string
	lastIdx := len(ctlg) - 1
	trait := (page - 1) * api.MaxProductPerPage
	t.Log("first index: ", ctlg[0].Id, ", last index: ", ctlg[lastIdx].Id)
	assert.Equal(t, ctlg[0].Id, DummyData[trait])

	trait += lastIdx
	if page != maxPage {
		last_id = ctlg[lastIdx].Id.Hex()
	}

	log.Print("trait: ", trait)
	assert.Equal(t, ctlg[lastIdx].Id, DummyData[trait])
	return last_id

}
