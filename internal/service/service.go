package service

import (
	"errors"
	"fmt"
	"net/http"

	"app/internal/entity"
	"app/internal/petition"
)

type InfoServices struct {
	DBHost    string
	DBPort    string
	TokenHost string
	TokenPort string
	Secret    string
}

type Service interface {
	SignUp(string, string, string) (string, error)
	SignIn(string, string) (string, error)
	LogOut(string) error
	GetAllUsers() ([]entity.User, error)
	Profile(string) (entity.User, error)
	DeleteAccount(string) error
}

// service ...
type service struct {
	client                    petition.HTTPClient
	dbHost, tokenHost, secret string
}

var (
	ErrResponse      = errors.New("error to response")
	ErrTokenNotValid = errors.New("token not validate")
	ErrWebServer     = errors.New("error from web server")
)

// NewService ...
func NewService(client petition.HTTPClient, is *InfoServices) *service {
	return &service{
		client,
		"http://" + is.DBHost + ":" + is.DBPort, "http://" + is.TokenHost + ":" + is.TokenPort, is.Secret,
	}
}

// SignUp ...
func (s *service) SignUp(username, password, email string) (token string, err error) {
	var (
		errorDBResponse    entity.ErrorResponse
		idResponse         entity.IDErrorResponse
		tokenResponse      entity.Token
		errorTokenResponse entity.ErrorResponse
	)

	if err = petition.RequestFunc(
		s.client,
		entity.UsernamePasswordEmailRequest{
			Username: username,
			Password: password,
			Email:    email,
		},
		petition.NewHTTPComponents(
			s.dbHost+"/user",
			http.MethodPost,
		),
		&errorDBResponse,
	); err != nil {
		return "", fmt.Errorf("%w:%s", ErrWebServer, err.Error())
	}

	if errorDBResponse.Err != "" {
		return "", fmt.Errorf("%w:%s", ErrWebServer, errorDBResponse.Err)
	}

	if err = petition.RequestFunc(
		s.client,
		entity.UsernameRequest{
			Username: username,
		},
		petition.NewHTTPComponents(
			s.dbHost+"/id/username",
			http.MethodGet,
		),
		&idResponse,
	); err != nil {
		return "", fmt.Errorf("%w:%s", ErrWebServer, err.Error())
	}

	if idResponse.Err != "" {
		return "", fmt.Errorf("%w:%s", ErrWebServer, idResponse.Err)
	}

	if err = petition.RequestFunc(
		s.client,
		entity.IDUsernameEmailSecretRequest{
			ID:       idResponse.ID,
			Username: username,
			Email:    email,
			Secret:   s.secret,
		},
		petition.NewHTTPComponents(
			s.tokenHost+"/generate",
			http.MethodPost,
		),
		&tokenResponse,
	); err != nil {
		return "", fmt.Errorf("%w:%s", ErrWebServer, err.Error())
	}

	if err = petition.RequestFunc(
		s.client,
		entity.Token{
			Token: tokenResponse.Token,
		},
		petition.NewHTTPComponents(
			s.tokenHost+"/token",
			http.MethodPost,
		),
		&errorTokenResponse,
	); err != nil {
		return "", fmt.Errorf("%w:%s", ErrWebServer, err.Error())
	}

	if errorTokenResponse.Err != "" {
		return "", fmt.Errorf("%w:%s", ErrWebServer, errorTokenResponse.Err)
	}

	return tokenResponse.Token, nil
}

// SignIn ...
func (s *service) SignIn(username, password string) (token string, err error) {
	var (
		userErrorResponse entity.UserErrorResponse
		tokenResponse     entity.Token
		errorResponse     entity.ErrorResponse
	)

	if err = petition.RequestFunc(
		s.client,
		entity.UsernamePasswordRequest{
			Username: username,
			Password: password,
		},
		petition.NewHTTPComponents(
			s.dbHost+"/user/username_password",
			http.MethodGet,
		),
		&userErrorResponse,
	); err != nil {
		return "", fmt.Errorf("%w:%s", ErrWebServer, err.Error())
	}

	if userErrorResponse.Err != "" {
		return "", fmt.Errorf("%w:%s", ErrWebServer, userErrorResponse.Err)
	}

	if err = petition.RequestFunc(
		s.client,
		entity.IDUsernameEmailSecretRequest{
			ID:       userErrorResponse.User.ID,
			Username: userErrorResponse.User.Username,
			Email:    userErrorResponse.User.Email,
			Secret:   s.secret,
		},
		petition.NewHTTPComponents(
			s.tokenHost+"/generate",
			http.MethodPost,
		),
		&tokenResponse,
	); err != nil {
		return "", fmt.Errorf("%w:%s", ErrWebServer, err.Error())
	}

	if err = petition.RequestFunc(
		s.client,
		entity.Token{
			Token: tokenResponse.Token,
		},
		petition.NewHTTPComponents(
			s.tokenHost+"/token",
			http.MethodPost,
		),
		&errorResponse,
	); err != nil {
		return "", fmt.Errorf("%w:%s", ErrWebServer, err.Error())
	}

	if errorResponse.Err != "" {
		return "", fmt.Errorf("%w:%s", ErrWebServer, errorResponse.Err)
	}

	return tokenResponse.Token, nil
}

// LogOut ...
func (s *service) LogOut(token string) (err error) {
	var (
		checkErrorResponse entity.CheckErrResponse
		errorResponse      entity.ErrorResponse
	)

	if err = petition.RequestFunc(
		s.client,
		entity.Token{
			Token: token,
		},
		petition.NewHTTPComponents(
			s.tokenHost+"/check",
			http.MethodPost,
		),
		&checkErrorResponse,
	); err != nil {
		return fmt.Errorf("%w:%s", ErrWebServer, err.Error())
	}

	if checkErrorResponse.Err != "" {
		return fmt.Errorf("%w:%s", ErrWebServer, checkErrorResponse.Err)
	}

	if !checkErrorResponse.Check {
		err = ErrTokenNotValid

		return err
	}

	if err = petition.RequestFunc(
		s.client,
		entity.Token{
			Token: token,
		},
		petition.NewHTTPComponents(
			s.tokenHost+"/token",
			http.MethodDelete,
		),
		&errorResponse,
	); err != nil {
		return fmt.Errorf("%w:%s", ErrWebServer, err.Error())
	}

	if errorResponse.Err != "" {
		return fmt.Errorf("%w:%s", ErrWebServer, errorResponse.Err)
	}

	return nil
}

// GetAllUsers  ...
func (s *service) GetAllUsers() (users []entity.User, err error) {
	var usersErrorResponse entity.UsersErrorResponse

	if err = petition.RequestFuncWithoutBody(
		s.client,
		petition.NewHTTPComponents(
			s.dbHost+"/users",
			http.MethodGet,
		),
		&usersErrorResponse,
	); err != nil {
		return nil, fmt.Errorf("%w:%s", ErrWebServer, err.Error())
	}

	if usersErrorResponse.Err != "" {
		return nil, fmt.Errorf("%w:%s", ErrWebServer, usersErrorResponse.Err)
	}

	return usersErrorResponse.Users, nil
}

// Profile  ...
func (s *service) Profile(token string) (user entity.User, err error) {
	var (
		checkErrorResponse         entity.CheckErrResponse
		idUsernameEmailErrResponse entity.IDUsernameEmailErrResponse
		userErrorResponse          entity.UserErrorResponse
	)

	if err = petition.RequestFunc(
		s.client,
		entity.Token{
			Token: token,
		},
		petition.NewHTTPComponents(
			s.tokenHost+"/check",
			http.MethodPost,
		),
		&checkErrorResponse,
	); err != nil {
		return entity.User{}, fmt.Errorf("%w:%s", ErrWebServer, err.Error())
	}

	if checkErrorResponse.Err != "" {
		return entity.User{}, fmt.Errorf("%w:%s", ErrWebServer, checkErrorResponse.Err)
	}

	if !checkErrorResponse.Check {
		err = ErrTokenNotValid

		return entity.User{}, err
	}

	if err = petition.RequestFunc(
		s.client,
		entity.TokenSecretRequest{
			Token:  token,
			Secret: s.secret,
		},
		petition.NewHTTPComponents(
			s.tokenHost+"/extract",
			http.MethodPost,
		),
		&idUsernameEmailErrResponse,
	); err != nil {
		return entity.User{}, fmt.Errorf("%w:%s", ErrWebServer, err.Error())
	}

	if idUsernameEmailErrResponse.Err != "" {
		return entity.User{}, fmt.Errorf("%w:%s", ErrWebServer, idUsernameEmailErrResponse.Err)
	}

	if err = petition.RequestFunc(
		s.client,
		entity.IDRequest{
			ID: idUsernameEmailErrResponse.ID,
		},
		petition.NewHTTPComponents(
			s.dbHost+"/user/id",
			http.MethodGet,
		),
		&userErrorResponse,
	); err != nil {
		return entity.User{}, fmt.Errorf("%w:%s", ErrWebServer, err.Error())
	}

	if userErrorResponse.Err != "" {
		return entity.User{}, fmt.Errorf("%w:%s", ErrWebServer, userErrorResponse.Err)
	}

	return userErrorResponse.User, nil
}

// DeleteAccount  ...
func (s *service) DeleteAccount(token string) (err error) {
	var (
		checkErrorResponse         entity.CheckErrResponse
		idUsernameEmailErrResponse entity.IDUsernameEmailErrResponse
		errorResponse              entity.ErrorResponse
	)

	if err = petition.RequestFunc(
		s.client,
		entity.Token{
			Token: token,
		},
		petition.NewHTTPComponents(
			s.tokenHost+"/check",
			http.MethodPost,
		),
		&checkErrorResponse,
	); err != nil {
		return fmt.Errorf("%w:%s", ErrWebServer, err.Error())
	}

	if checkErrorResponse.Err != "" {
		return fmt.Errorf("%w:%s", ErrWebServer, checkErrorResponse.Err)
	}

	if !checkErrorResponse.Check {
		err = ErrTokenNotValid

		return err
	}

	if err = petition.RequestFunc(
		s.client,
		entity.TokenSecretRequest{
			Token:  token,
			Secret: s.secret,
		},
		petition.NewHTTPComponents(
			s.tokenHost+"/extract",
			http.MethodPost,
		),
		&idUsernameEmailErrResponse,
	); err != nil {
		return fmt.Errorf("%w:%s", ErrWebServer, err.Error())
	}

	if idUsernameEmailErrResponse.Err != "" {
		return fmt.Errorf("%w:%s", ErrWebServer, idUsernameEmailErrResponse.Err)
	}

	if err = petition.RequestFunc(
		s.client,
		entity.IDRequest{
			ID: idUsernameEmailErrResponse.ID,
		},
		petition.NewHTTPComponents(
			s.dbHost+"/user",
			http.MethodDelete,
		),
		&errorResponse,
	); err != nil {
		return fmt.Errorf("%w:%s", ErrWebServer, err.Error())
	}

	return nil
}
