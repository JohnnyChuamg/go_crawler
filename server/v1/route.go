package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	authMiddleware "api_server/internal/middlewares/auth"
	"api_server/internal/repositories"
	workSrv "api_server/internal/services/work"
	workHandler "api_server/server/v1/routes/work"
)

type IRoute interface {
	SetRoute(engine *gin.RouterGroup)
}

func InitRoutes(engine *gin.Engine) error {
	v1 := engine.Group("v1")

	v1.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "["+c.ClientIP()+"] pong!!!")
	})

	routes, err := createRoutes()

	if err != nil {
		return err
	}

	for _, route := range routes {
		route.SetRoute(v1)
	}

	return nil
}

func createRoutes() ([]IRoute, error) {
	auth := authMiddleware.NewAuth()
	routes := []IRoute{
		workHandler.New(workSrv.New(repositories.NewWork()), auth),
	}
	return routes, nil
}
