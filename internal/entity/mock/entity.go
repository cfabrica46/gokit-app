package mock

import "errors"

const (
	IDTest       int    = 1
	UsernameTest string = "username"
	PasswordTest string = "password"
	EmailTest    string = "email@email.com"
	SecretTest   string = "secret"

	URLTest       string = "localhost:8080"
	DBHostTest    string = "db"
	TokenHostTest string = "token"
	PortTest      string = "8080"
	TokenTest     string = "token"

	NameNoError string = "NoError"
)

var (
	ErrWebServer        = errors.New("error from web server")
	ErrNotTypeIndicated = errors.New("response is not of the type indicated")
)
