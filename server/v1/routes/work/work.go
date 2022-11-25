package work

import (
	"api_server/internal/infrastures/configuration"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"

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
	g.GET("pic", c.Pic)
	g.GET("print", c.Print)
}

func (c *Work) Pic(ctx *gin.Context) {
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
	imgContentChannel := make(chan string, len(imgs))
	var wg, wg2 sync.WaitGroup
	for _, img := range imgs {
		if img == nil || len(img) <= 0 {
			continue
		}
		wg.Add(1)
		go getDataImageBase64Async(img, imgContentChannel, &wg)
	}
	wg2.Add(1)
	//主序等channel close會卡死，因此另開一序作寫入buffer
	go func(channel <-chan string, buffer *bytes.Buffer, wg *sync.WaitGroup) {
		defer wg.Done()

		for val := range channel {
			buffer.WriteString(fmt.Sprintf("<img src=\"%s\" />", val))
		}
	}(imgContentChannel, &buffer, &wg2)

	//等編碼的序
	wg.Wait()
	//等編碼的序完成後，關閉channel( 即不會再有任務寫入至channel中)
	close(imgContentChannel)
	//等所有資料都轉移到buffer中。
	wg2.Wait()
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

func getDataImageBase64Async(source []byte, channel chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	var base64Encoding string

	mimeType := http.DetectContentType(source)

	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}
	base64Encoding += base64.StdEncoding.EncodeToString(source)
	channel <- base64Encoding
}
