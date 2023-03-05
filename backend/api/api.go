package api

import (
	"github.com/ikmv2/backend/pkg/repository"

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
		service: Service{repo: repo},
	}
}

func (a Api) ExposeRoute() {
	a.server.GET("ping", func(c echo.Context) error {
		return c.JSON(200, JsonMap{"reponse": "pong"})
	})
}

func (a Api) StartServer() {
	a.server.Start(":8081")
}
