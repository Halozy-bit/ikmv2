package api

import (
	"net/http"
	"strconv"
	"time"

	asynctask "github.com/ikmv2/backend/pkg/async_task"
	"github.com/ikmv2/backend/pkg/repository"
	"github.com/ikmv2/backend/pkg/sidejob"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/labstack/echo/v4"
)

type JsonMap map[string]interface{}

type Api struct {
	server  *echo.Echo
	service Service
}

func NewEndpoint(repo repository.Repository) Api {
	return Api{
		server:  echo.New(),
		service: &ServiceCirclePage{repo: repo},
	}
}

// TODO
// Check data updated
func (a Api) StartSideJob(db *mongo.Database) error {
	if err := asynctask.AddTask(&sidejob.RefreshCatalogPage{
		Db: db,
		TaskIdentifier: asynctask.TaskIdentifier{
			Name:     "refresh catalog page",
			Interval: time.Hour * 10,
		},
	}); err != nil {
		return err
	}

	if err := asynctask.AddTask(&sidejob.RefreshCatalogCategoryPage{
		Db: db,
		TaskIdentifier: asynctask.TaskIdentifier{
			Name:     "refresh catalog per category",
			Interval: time.Hour * 15,
		},
	}); err != nil {
		return err
	}

	err := asynctask.Start(time.Second * 5)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 5)
	return nil
}

func (a Api) ExposeRoute() {
	a.server.GET("/ping", a.Pong)

	a.server.GET("/catalog/:page", a.GetCatalog)
	a.server.GET("/catalog/:category/:page", a.GetCatalogByCategory)
	a.server.GET("/product/:id", a.GetProduct)
	a.server.GET("/page/count", a.TotalPage)
	a.server.GET("/page/count/:category", a.TotalPage)
	a.server.GET("/umkm/:id", a.GetUmkm)
}

func (a Api) Pong(c echo.Context) error {
	return c.JSON(http.StatusOK, JsonMap{"reponse": "pong"})
}

// TODO
// Return multiple error from binding
func (a Api) GetCatalog(c echo.Context) error {
	req := CatalogGetProduct{}
	var err error

	req.Page, err = strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.JSON(http.StatusNotFound, JsonMap{"message": "page not exist"})
	}

	catalog, cErr := a.service.CatalogList(
		c.Request().Context(),
		req.Page,
	)

	switch cErr {
	case nil:
		err = c.JSON(http.StatusOK, JsonMap{"catalog": catalog})
	case mongo.ErrNoDocuments:
		err = c.JSON(http.StatusNoContent, JsonMap{"message": "no content"})
	default:
		err = c.JSON(http.StatusInternalServerError, JsonMap{"message": cErr.Error()})
	}
	return err
}

func (a Api) GetCatalogByCategory(c echo.Context) error {
	req := CatalogGetProduct{}
	var err error

	req.Category = c.Param("category")

	req.Page, err = strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.JSON(http.StatusNotFound, JsonMap{"message": "page not exist"})
	}

	catalog, cErr := a.service.CatalogListByCategory(
		c.Request().Context(),
		req.Page, req.Category,
	)

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

func (a Api) GetProduct(c echo.Context) error {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.JSON(http.StatusNotFound, JsonMap{"message": "page not exist"})
	}

	p, err := a.service.Product(c.Request().Context(), id)

	switch err {
	case nil:
		err = c.JSON(http.StatusOK, JsonMap{"product": p})
	case mongo.ErrNoDocuments:
		err = c.JSON(http.StatusNoContent, JsonMap{"message": "no content"})
	default:
		err = c.JSON(http.StatusInternalServerError, JsonMap{"message": err.Error()})
	}

	return err
}

func (a Api) TotalPage(c echo.Context) error {
	req := c.Param("category")
	var ttl int
	if req == "" {
		ttl = a.service.TotalPage()
	} else {
		ttl = a.service.TotalPage(req)
	}
	return c.JSON(200, JsonMap{"total": ttl})
}

func (a Api) GetUmkm(c echo.Context) error {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.JSON(http.StatusNotFound, JsonMap{"message": "page not exist"})
	}
	u, err := a.service.Umkm(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, JsonMap{"message": "umkm not found"})
	}

	return c.JSON(200, JsonMap{"umkm": u})
}

func (a Api) Server() *echo.Echo {
	return a.server
}

func (a Api) StartServer(address string) {
	a.server.Start(address)
}
