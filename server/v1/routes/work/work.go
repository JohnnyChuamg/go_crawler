package work

import (
	"net/http"

	"github.com/gin-gonic/gin"

	middlewares "api_server/internal/middlewares/auth"
	srv "api_server/internal/services/work"
)

type Work struct {
	srv            *srv.Work
	authMiddleware *middlewares.Auth
}

func New(srv *srv.Work, authMiddleware *middlewares.Auth) *Work {
	return &Work{
		srv:            srv,
		authMiddleware: authMiddleware,
	}
}

func (c *Work) SetRoute(engine *gin.RouterGroup) {
	//g := engine.Group("work", c.authMiddleware.Authentication)
	g := engine.Group("work")
	g.GET("pic", c.Action)
}

func (c *Work) Action(ctx *gin.Context) {
	//const target = "https://tw.yahoo.com"
	//const target = "https://www.taiwan.net.tw/m1.aspx?sNo=0012076"
	url := ctx.Query("url")
	if url == "" {
		ctx.JSON(http.StatusBadRequest, "url cant be null or empty")
		return
	}
	data, err := c.srv.CrawlerImage(url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	_, err = ctx.Writer.Write(data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}
}
