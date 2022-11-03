package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	configManager "api_server/internal/infrastures/configuration"
	baseRoute "api_server/server/routes"
	v1 "api_server/server/v1"
)

func NewServer() error {

	ginEngine := newGinEngine()

	if err := baseRoute.InitRoutes(ginEngine); err != nil {
		return err
	}

	err := v1.InitRoutes(ginEngine)

	// start http server
	httpServer := &http.Server{
		Addr:    configManager.GlobalConfig.HTTPBind,
		Handler: ginEngine,
	}
	go func() {
		// service connection
		log.Info().Msgf("main: Listening and serving HTTP on %s", httpServer.Addr)
		err = httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Panic().Msgf("main: http server listen failed: %v", err)
		}
	}()
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan
	log.Printf("main: shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Panic().Msgf("main: http server shutdown error: %v", err)
	} else {
		log.Info().Msgf("main: gracefully stopped")
	}
	return nil
}

func newGinEngine() *gin.Engine {
	defaultRouter := gin.Default()

	corsConfig := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:    []string{"*", "Authorization", "Content-Type", "Origin", "Content-Length"},
		// firefox 和 safari 不支援 *, 所以需要一個一個打，但更好的是要抓 Access-Control-Request-Headers
		// https://stackoverflow.com/questions/54666673/cors-check-fails-for-firefox-but-passes-for-chrome
	}

	defaultRouter.Use(cors.New(corsConfig))

	defaultRouter.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "StatusNotFound")
	})

	return defaultRouter
}
