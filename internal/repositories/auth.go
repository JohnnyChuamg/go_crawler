package repositories

const (
	_username = "johnny"
	_password = "123456"
)

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
}

func (*Auth) SignIn(username string, password string) (bool, error) {
	return username == _username && password == _password, nil
}
