package work

import (
	"api_server/internal/infrastures/configuration"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

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
	//_, err = ctx.Writer.Write(data)
	//if err != nil {
	//	ctx.JSON(http.StatusInternalServerError, err.Error())
	//	return
	//}
	base64Encoding := getDataImageBase64(data)

	html := fmt.Sprintf("<img src=\"%s\" />", base64Encoding)
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func (c *Work) Print(ctx *gin.Context) {
	url := ctx.Query("url")
	if url == "" {
		ctx.JSON(http.StatusBadRequest, "url cant be null or empty")
		return
	}
	imgs, err := c.srv.CrawlerImagesAsync(url)
	//imgs, err := c.srv.CrawlerImages(url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var buffer bytes.Buffer
	for _, img := range imgs {
		if len(img) <= 0 {
			continue
		}
		//if idx >= 5 {
		//	break
		//}
		base64Encoding := getDataImageBase64(img)
		buffer.WriteString(fmt.Sprintf("<img src=\"%s\" />", base64Encoding))
	}
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", buffer.Bytes())
}

func getDataImageBase64(source []byte) string {
	var base64Encoding string

	mimeType := http.DetectContentType(source)

	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}
	base64Encoding += base64.StdEncoding.EncodeToString(source)
	return base64Encoding
}
