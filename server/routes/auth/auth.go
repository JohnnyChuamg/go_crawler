package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"api_server/internal/infrastures/configuration"
	"api_server/internal/services/auth"
	"api_server/internal/services/auth/contracts"
)

type Auth struct {
	srv *auth.Auth
}

func (c *Auth) SetRoute(engine *gin.RouterGroup) {
	g := engine.Group("auth")
	g.POST("signIn", c.SignIn)
}

func New(srv *auth.Auth) *Auth {
	return &Auth{
		srv: srv,
	}
}

func (c *Auth) SignIn(ctx *gin.Context) {
	var model contracts.SignIn
	if err := ctx.BindJSON(model); err != nil {
		ctx.JSON(http.StatusInternalServerError, "parse request body to json failed: "+err.Error())
		return
	}

	if model.Username == "" || model.Password == "" {
		ctx.JSON(http.StatusBadRequest, "username or password invalid")
		return
	}
	token, err := c.srv.SignIn(model)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "unknown error: "+err.Error())
		return
	} else if token == "" {
		ctx.JSON(http.StatusBadRequest, "username or password invalid")
		return
	}
	ctx.SetCookie(configuration.GlobalConfig.Auth.Key, token, configuration.GlobalConfig.Auth.TTL, "/", "localhost", false, true)
	ctx.JSON(http.StatusNoContent, "")
}
