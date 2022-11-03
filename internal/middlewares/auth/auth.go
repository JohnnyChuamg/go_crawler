package auth

import (
	"net/http"

	"api_server/internal/infrastures/configuration"
	"github.com/gin-gonic/gin"
)

type Auth struct {
}

func New() *Auth {
	return &Auth{}
}

// Authentication Determines whether users are who they claim to be 確定用戶是否是他們聲稱的人
func (m *Auth) Authentication(ctx *gin.Context) {
	token, err := ctx.Cookie(configuration.GlobalConfig.Auth.Key)
	if err != nil || token == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.Next()
}

// Authorization Determines what users can and cannot access 確定用戶可以訪問和不能訪問的內容
func (m *Auth) Authorization() {

}
