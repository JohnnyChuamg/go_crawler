package auth

import (
	"api_server/internal/repositories"
	"api_server/internal/services/auth/contracts"
	"crypto/md5"
	"errors"
	"fmt"
	"time"
)

type Auth struct {
	repo *repositories.Auth
}

func New(auth *repositories.Auth) *Auth {
	return &Auth{
		repo: auth,
	}
}

func (srv *Auth) SignIn(model contracts.SignIn) (string, error) {
	if model.Username == "" || model.Password == "" {
		return "", errors.New("username or password cant be null or empty")
	}
	if r, err := srv.repo.SignIn(model.Username, model.Password); err != nil {
		return "", err
	} else if !r {
		return "", nil
	}
	return encrypt(model.Username, model.Password), nil
}

func encrypt(username string, password string) string {
	t := fmt.Sprint(time.Now().Format("20060102"))
	result := fmt.Sprintf("%s%s%s", username, password, t)
	return fmt.Sprintf("%x", md5.Sum([]byte(result)))
}
