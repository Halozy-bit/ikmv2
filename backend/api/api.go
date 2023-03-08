package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ikmv2/backend/pkg/cache"
	"github.com/ikmv2/backend/pkg/repository"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/labstack/echo/v4"
)

type JsonMap map[string]interface{}

type Api struct {
	server  *echo.Echo
	service Service
}

func NewEndpoint(repo repository.Repository) Api {
	if cache.Get(cache.TopCatalog).(string) == "" || cache.Get(cache.BottomCatalog).(string) == "" {
		first, err := repo.FirstItem()
		if err != nil {
			log.Fatalln("Error get first item")
		}

		last, err := repo.LastItem()
		if err != nil {
			log.Fatalln("Error get last item")
		}

		cache.Store(cache.TopCatalog, first.Id.Hex())
		cache.Store(cache.BottomCatalog, last.Id.Hex())
	}

	return Api{
		server:  echo.New(),
		service: Service{repo: repo},
	}
}

func (a Api) ExposeRoute() {
	a.server.GET("/ping", a.Pong)

	a.server.GET("/catalog/:page", a.GetCatalog)
}

func (a Api) Pong(c echo.Context) error {
	return c.JSON(http.StatusOK, JsonMap{"reponse": "pong"})
}

// TODO
// Return multiple error from binding
func (a Api) GetCatalog(c echo.Context) error {
	req := new(CatalogGetProduct)
	var err error

	req.Page, err = strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.JSON(http.StatusNotFound, JsonMap{"message": "page not exist"})
	}

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, JsonMap{"message": err.Error()})
	}

	log.Println(req)

	var catalog []repository.DocCatalog
	var cErr error
	if len(req.Category) > 3 {
		catalog, cErr = a.service.CatalogListByCategory(
			c.Request().Context(),
			req.Page, req.Category,
			req.LastID,
		)
	} else {
		catalog, cErr = a.service.CatalogList(
			c.Request().Context(),
			req.Page,
			req.LastID,
		)
	}

	switch cErr {
	case nil:
		err = c.JSON(http.StatusOK, JsonMap{"catalog": catalog})
	case mongo.ErrNoDocuments:
		err = c.JSON(http.StatusNoContent, JsonMap{"message": "no content"})
	default:
		err = c.JSON(http.StatusInternalServerError, JsonMap{"message": err.Error()})
	}

	return err
}

func (a Api) Server() *echo.Echo {
	return a.server
}

func (a Api) StartServer(address string) {
	a.server.Start(address)
}
