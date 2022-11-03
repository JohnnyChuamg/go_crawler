package work

import (
	"api_server/internal/infrastures/configuration"
	"bytes"
	"fmt"
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
	var g *gin.RouterGroup
	if configuration.GlobalConfig.Auth.Enable {
		g = engine.Group("work", c.authMiddleware.Authentication)
	} else {
		g = engine.Group("work")
	}
	//g := engine.Group("work")
	g.GET("pic", c.Action)
	g.GET("print", c.Print)
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
		return
	}
	_, err = ctx.Writer.Write(data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
}

func (c *Work) Print(ctx *gin.Context) {
	url := ctx.Query("url")
	if url == "" {
		ctx.JSON(http.StatusBadRequest, "url cant be null or empty")
		return
	}
	//url := "https://www.taiwan.net.tw/m1.aspx?sNo=0012076"
	data, err := c.srv.CrawlerImagesAndPrintAll(url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	var buffer bytes.Buffer

	for _, img := range data {
		buffer.WriteString(fmt.Sprintf("<img src=\"/%s\" />", img))
	}
	//
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", buffer.Bytes())
}
