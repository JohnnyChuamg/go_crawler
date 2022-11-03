package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"api_server/internal/repositories"
	authSrv "api_server/internal/services/auth"
	authHandler "api_server/server/routes/auth"
)

type IRoute interface {
	SetRoute(engine *gin.RouterGroup)
}

func InitRoutes(engine *gin.Engine) error {
	group := engine.Group("")

	group.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "["+c.ClientIP()+"] pong!!!")
	})

	group.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!!!")
	})

	routes, err := createRoutes()

	if err != nil {
		return err
	}

	for _, route := range routes {
		route.SetRoute(group)
	}

	return nil
}

func createRoutes() ([]IRoute, error) {
	routes := []IRoute{
		authHandler.New(authSrv.New(repositories.NewAuth())),
	}
	return routes, nil
}
