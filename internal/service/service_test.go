package service_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"app/internal/entity"
	"app/internal/entity/mock"
	"app/internal/service"

	httpMock "app/internal/service/mock"

	"github.com/stretchr/testify/assert"
)

type errorHTTPComponents struct {
	errorURL, errorMethod string
}

type signUpTestStruct struct {
	name                            string
	inUsername, inPassword, inEmail string
	url                             string
	method                          string
	isError                         bool
	isErrorInsideRequest            bool
}

type profileTestStruct struct {
	name                 string
	inToken              string
	url                  string
	method               string
	outUser              entity.User
	outCheck             bool
	isError              bool
	isErrorInsideRequest bool
}

func newErrorHTTPComponets(url, method string) errorHTTPComponents {
	return errorHTTPComponents{
		errorURL:    url,
		errorMethod: method,
	}
}

func getSignUpTestEntity() []signUpTestStruct {
	return []signUpTestStruct{
		{
			name:       "NoError",
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			inEmail:    mock.EmailTest,
			isError:    false,
			url:        "http://db:8080/user",
			method:     http.MethodPost,
		},
		{
			name:       "ErrorInsertUser",
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			inEmail:    mock.EmailTest,
			isError:    true,
			url:        "http://db:8080/user",
			method:     http.MethodPost,
		},
		{
			name:                 "ErrorInsideInsertUser",
			inUsername:           mock.UsernameTest,
			inPassword:           mock.PasswordTest,
			inEmail:              mock.EmailTest,
			isError:              true,
			isErrorInsideRequest: true,
			url:                  "http://db:8080/user",
			method:               http.MethodPost,
		},
		{
			name:       "ErrorGetID",
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			inEmail:    mock.EmailTest,
			isError:    true,
			url:        "http://db:8080/id/username",
			method:     http.MethodGet,
		},
		{
			name:                 "ErrorInsideGetID",
			inUsername:           mock.UsernameTest,
			inPassword:           mock.PasswordTest,
			inEmail:              mock.EmailTest,
			isError:              true,
			isErrorInsideRequest: true,
			url:                  "http://db:8080/id/username",
			method:               http.MethodGet,
		},
		{
			name:       "ErrorGenerate",
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			inEmail:    mock.EmailTest,
			isError:    true,
			url:        "http://token:8080/generate",
			method:     http.MethodPost,
		},
		{
			name:       "ErrorSetToken",
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			inEmail:    mock.EmailTest,
			isError:    true,
			url:        "http://token:8080/token",
			method:     http.MethodPost,
		},
		{
			name:                 "ErrorInsideSetToken",
			inUsername:           mock.UsernameTest,
			inPassword:           mock.PasswordTest,
			inEmail:              mock.EmailTest,
			isError:              true,
			isErrorInsideRequest: true,
			url:                  "http://token:8080/token",
			method:               http.MethodPost,
		},
	}
}

func getProfileTestEntity() []profileTestStruct {
	return []profileTestStruct{
		{
			name:    "NoError",
			inToken: mock.TokenTest,
			outUser: entity.User{
				ID:       mock.IDTest,
				Username: mock.UsernameTest,
				Password: mock.PasswordTest,
				Email:    mock.EmailTest,
			},
			outCheck: true,
			isError:  false,
			url:      "http://db:8080/user/id",
			method:   http.MethodGet,
		},
		{
			name:     "ErrorCheckToken",
			inToken:  mock.TokenTest,
			outUser:  entity.User{},
			outCheck: true,
			isError:  true,
			url:      "http://token:8080/check",
			method:   http.MethodPost,
		},
		{
			name:                 "ErrorInsideCheckToken",
			inToken:              mock.TokenTest,
			outUser:              entity.User{},
			outCheck:             true,
			isError:              true,
			isErrorInsideRequest: true,
			url:                  "http://token:8080/check",
			method:               http.MethodPost,
		},
		{
			name:     "FalseCheckToken",
			inToken:  mock.TokenTest,
			outUser:  entity.User{},
			outCheck: false,
			isError:  false,
			url:      "http://token:8080/check",
			method:   http.MethodPost,
		},
		{
			name:     "ErrorExtractToken",
			inToken:  mock.TokenTest,
			outUser:  entity.User{},
			outCheck: true,
			isError:  true,
			url:      "http://token:8080/extract",
			method:   http.MethodPost,
		},
		{
			name:                 "ErrorInsideExtractToken",
			inToken:              mock.TokenTest,
			outUser:              entity.User{},
			outCheck:             true,
			isError:              true,
			isErrorInsideRequest: true,
			url:                  "http://token:8080/extract",
			method:               http.MethodPost,
		},
		{
			name:     "ErrorGetID",
			inToken:  mock.TokenTest,
			outUser:  entity.User{},
			outCheck: true,
			isError:  true,
			url:      "http://db:8080/user/id",
			method:   http.MethodGet,
		},
		{
			name:                 "ErrorInsideGetID",
			inToken:              mock.TokenTest,
			outUser:              entity.User{},
			outCheck:             true,
			isError:              true,
			isErrorInsideRequest: true,
			url:                  "http://db:8080/user/id",
			method:               http.MethodGet,
		},
	}
}

func TestSignUp(t *testing.T) {
	t.Parallel()

	infoServiceTest := service.InfoServices{
		DBHost:    mock.DBHostTest,
		DBPort:    mock.PortTest,
		TokenHost: mock.TokenHostTest,
		TokenPort: mock.PortTest,
		Secret:    mock.SecretTest,
	}

	for _, tt := range getSignUpTestEntity() {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultToken string
			var resultErr error
			var tokenResponse, errorResponse string

			var respFunc func(*http.Request) (*http.Response, error)

			if tt.isError {
				errorResponse = mock.ErrWebServer.Error()
			} else {
				tokenResponse = mock.TokenTest
			}

			responseJSON := `{
						"token": "token",
						"id": 1
			}`

			if tt.isError {
				respFunc = getIsErrorMock(
					tt.isErrorInsideRequest,
					newErrorHTTPComponets(tt.url, tt.method),
					responseJSON,
				)
			} else {
				respFunc = getMock(
					responseJSON,
				)
			}

			mockHTTP := httpMock.NewMockClient(respFunc)

			svc := service.NewService(
				mockHTTP,
				&infoServiceTest,
			)

			resultToken, resultErr = svc.SignUp(tt.inUsername, tt.inPassword, tt.inEmail)

			if !tt.isError {
				assert.Nil(t, resultErr)
			} else {
				assert.ErrorContains(t, resultErr, errorResponse)
			}
			assert.Equal(t, tokenResponse, resultToken)
		})
	}
}

func TestSignIn(t *testing.T) {
	t.Parallel()

	infoServiceTest := service.InfoServices{
		DBHost:    mock.DBHostTest,
		DBPort:    mock.PortTest,
		TokenHost: mock.TokenHostTest,
		TokenPort: mock.PortTest,
		Secret:    mock.SecretTest,
	}

	for _, tt := range []struct {
		name                   string
		inUsername, inPassword string
		url                    string
		method                 string
		isError                bool
		isErrorInsideRequest   bool
	}{
		{
			name:       "NoError",
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			isError:    false,
			url:        "http://token:8080/generate",
			method:     http.MethodPost,
		},
		{
			name:       "ErrorGetUser",
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			isError:    true,
			url:        "http://db:8080/user/username_password",
			method:     http.MethodGet,
		},
		{
			name:                 "ErrorInsideGetToken",
			inUsername:           mock.UsernameTest,
			inPassword:           mock.PasswordTest,
			isError:              true,
			isErrorInsideRequest: true,
			url:                  "http://db:8080/user/username_password",
			method:               http.MethodGet,
		},
		{
			name:       "ErrorGenerateToken",
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			isError:    true,
			url:        "http://token:8080/generate",
			method:     http.MethodPost,
		},
		{
			name:       "ErrorSetToken",
			inUsername: mock.UsernameTest,
			inPassword: mock.PasswordTest,
			isError:    true,
			url:        "http://token:8080/token",
			method:     http.MethodPost,
		},
		{
			name:                 "ErrorInsideSetToken",
			inUsername:           mock.UsernameTest,
			inPassword:           mock.PasswordTest,
			isError:              true,
			isErrorInsideRequest: true,
			url:                  "http://token:8080/token",
			method:               http.MethodPost,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultToken string
			var resultErr error
			var tokenResponse, errorResponse string
			var mockHTTP *httpMock.MockClient

			if tt.isError {
				errorResponse = mock.ErrWebServer.Error()
			} else {
				tokenResponse = mock.TokenTest
			}

			responseJSON := `{
					"token":"token",
					"user":{
						"id":       1,
						"username": "username",
						"password": "password",
						"email":    "email@email.com"
					}
			}`

			if tt.isError {
				mockHTTP = httpMock.NewMockClient(getIsErrorMock(
					tt.isErrorInsideRequest,
					newErrorHTTPComponets(tt.url, tt.method),
					responseJSON,
				))
			} else {
				mockHTTP = httpMock.NewMockClient(getMock(
					responseJSON,
				))
			}

			svc := service.NewService(
				mockHTTP,
				&infoServiceTest,
			)

			resultToken, resultErr = svc.SignIn(tt.inUsername, tt.inPassword)

			if !tt.isError {
				assert.Nil(t, resultErr)
			} else {
				assert.ErrorContains(t, resultErr, errorResponse)
			}
			assert.Equal(t, tokenResponse, resultToken)
		})
	}
}

func TestLogOut(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name                 string
		inToken              string
		url                  string
		method               string
		outCheck             bool
		isError              bool
		isErrorInsideRequest bool
	}{
		{
			name:     "NoError",
			inToken:  mock.TokenTest,
			outCheck: true,
			isError:  false,
			url:      "http://token:8080/check",
			method:   http.MethodPost,
		},
		{
			name:     "ErrorCheckToken",
			inToken:  mock.TokenTest,
			outCheck: true,
			isError:  true,
			url:      "http://token:8080/check",
			method:   http.MethodPost,
		},
		{
			name:                 "ErrorInsideCheckToken",
			inToken:              mock.TokenTest,
			outCheck:             true,
			isError:              true,
			isErrorInsideRequest: true,
			url:                  "http://token:8080/check",
			method:               http.MethodPost,
		},
		{
			name:     "FalseCheckToken",
			inToken:  mock.TokenTest,
			outCheck: false,
			isError:  false,
			url:      "http://token:8080/check",
			method:   http.MethodPost,
		},
		{
			name:     "ErrorDeleteToken",
			inToken:  mock.TokenTest,
			outCheck: true,
			isError:  true,
			url:      "http://token:8080/token",
			method:   http.MethodDelete,
		},
		{
			name:                 "ErrorDeleteToken",
			inToken:              mock.TokenTest,
			outCheck:             true,
			isError:              true,
			isErrorInsideRequest: true,
			url:                  "http://token:8080/token",
			method:               http.MethodDelete,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr error
			var errorResponse string
			var mockHTTP *httpMock.MockClient

			infoServiceTest := service.InfoServices{
				DBHost:    mock.DBHostTest,
				DBPort:    mock.PortTest,
				TokenHost: mock.TokenHostTest,
				TokenPort: mock.PortTest,
				Secret:    mock.SecretTest,
			}

			if tt.isError {
				errorResponse = mock.ErrWebServer.Error()
			}

			responseJSON := fmt.Sprintf(`{
					"check": %t
				}`, tt.outCheck)

			if tt.isError {
				mockHTTP = httpMock.NewMockClient(getIsErrorMock(
					tt.isErrorInsideRequest,
					newErrorHTTPComponets(tt.url, tt.method),
					responseJSON,
				))
			} else {
				mockHTTP = httpMock.NewMockClient(getMock(
					responseJSON,
				))
			}

			svc := service.NewService(
				mockHTTP,
				&infoServiceTest,
			)

			resultErr = svc.LogOut(tt.inToken)

			if !tt.isError {
				if tt.outCheck {
					assert.Nil(t, resultErr)
				} else {
					assert.ErrorContains(t, resultErr, errorResponse)
				}
			} else {
				assert.ErrorContains(t, resultErr, errorResponse)
			}
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name                 string
		url                  string
		method               string
		outUsers             []entity.User
		isError              bool
		isErrorInsideRequest bool
	}{
		{
			name: "NoError",
			outUsers: []entity.User{
				{
					ID:       mock.IDTest,
					Username: mock.UsernameTest,
					Password: mock.PasswordTest,
					Email:    mock.EmailTest,
				},
			},
			isError: false,
			url:     "http://db:8080/users",
			method:  http.MethodGet,
		},
		{
			name:     "ErrorGetAllUsers",
			outUsers: nil,
			isError:  true,
			url:      "http://db:8080/users",
			method:   http.MethodGet,
		},
		{
			name:                 "ErrorInsideGetAllUsers",
			outUsers:             nil,
			isError:              true,
			isErrorInsideRequest: true,
			url:                  "http://db:8080/users",
			method:               http.MethodGet,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infoServiceTest := service.InfoServices{
				DBHost:    mock.DBHostTest,
				DBPort:    mock.PortTest,
				TokenHost: mock.TokenHostTest,
				TokenPort: mock.PortTest,
				Secret:    mock.SecretTest,
			}

			var resultUsers []entity.User
			var resultErr error
			var errorResponse string
			var mockHTTP *httpMock.MockClient

			if tt.isError {
				errorResponse = mock.ErrWebServer.Error()
			}

			responseJSON := `{
				"users":[
					{
						"username":"username",
						"password":"password",
						"email":"email@email.com",
						"id":1
					}
				]
			}`

			if tt.isError {
				mockHTTP = httpMock.NewMockClient(getIsErrorMock(
					tt.isErrorInsideRequest,
					newErrorHTTPComponets(tt.url, tt.method),
					responseJSON,
				))
			} else {
				mockHTTP = httpMock.NewMockClient(getMock(
					responseJSON,
				))
			}

			svc := service.NewService(
				mockHTTP,
				&infoServiceTest,
			)

			resultUsers, resultErr = svc.GetAllUsers()

			if !tt.isError {
				assert.Nil(t, resultErr)
			} else {
				assert.ErrorContains(t, resultErr, errorResponse)
			}
			assert.Equal(t, tt.outUsers, resultUsers)
		})
	}
}

func TestProfile(t *testing.T) {
	t.Parallel()

	for _, tt := range getProfileTestEntity() {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infoServiceTest := service.InfoServices{
				DBHost:    mock.DBHostTest,
				DBPort:    mock.PortTest,
				TokenHost: mock.TokenHostTest,
				TokenPort: mock.PortTest,
				Secret:    mock.SecretTest,
			}

			var resultUser entity.User
			var resultErr error
			var errorResponse string
			var mockHTTP *httpMock.MockClient

			if tt.isError {
				errorResponse = mock.ErrWebServer.Error()
			}

			responseJSON := fmt.Sprintf(`{
					"user":{
						"username":"username",
						"password":"password",
						"email":"email@email.com",
						"id":1
					},
					"id":1,
					"username":"usename",
					"email":"email@email.com",
					"check":%t
				}`, tt.outCheck)

			if tt.isError {
				mockHTTP = httpMock.NewMockClient(getIsErrorMock(
					tt.isErrorInsideRequest,
					newErrorHTTPComponets(tt.url, tt.method),
					responseJSON,
				))
			} else {
				mockHTTP = httpMock.NewMockClient(getMock(
					responseJSON,
				))
			}

			svc := service.NewService(
				mockHTTP,
				&infoServiceTest,
			)

			resultUser, resultErr = svc.Profile(tt.inToken)

			if !tt.isError {
				if tt.outCheck {
					assert.Nil(t, resultErr)
				} else {
					assert.ErrorContains(t, resultErr, errorResponse)
				}
			} else {
				assert.ErrorContains(t, resultErr, errorResponse)
			}

			assert.Equal(t, tt.outUser, resultUser)
		})
	}
}

func TestDeleteAccount(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name                 string
		inToken              string
		url                  string
		method               string
		outCheck             bool
		isError              bool
		isErrorInsideRequest bool
	}{
		{
			name:     "NoError",
			inToken:  mock.TokenTest,
			outCheck: true,
			isError:  false,
			url:      "http://token:8080/check",
			method:   http.MethodPost,
		},
		{
			name:     "ErrorCheckToken",
			inToken:  mock.TokenTest,
			outCheck: true,
			isError:  true,
			url:      "http://token:8080/check",
			method:   http.MethodPost,
		},
		{
			name:                 "ErrorInsideCheckToken",
			inToken:              mock.TokenTest,
			outCheck:             true,
			isError:              true,
			isErrorInsideRequest: true,
			url:                  "http://token:8080/check",
			method:               http.MethodPost,
		},
		{
			name:     "FalseCheckToken",
			inToken:  mock.TokenTest,
			outCheck: false,
			isError:  false,
			url:      "http://token:8080/check",
			method:   http.MethodPost,
		},
		{
			name:     "ErrorExtractToken",
			inToken:  mock.TokenTest,
			outCheck: true,
			isError:  true,
			url:      "http://token:8080/extract",
			method:   http.MethodPost,
		},
		{
			name:                 "ErrorInsideExtractToken",
			inToken:              mock.TokenTest,
			outCheck:             true,
			isError:              true,
			isErrorInsideRequest: true,
			url:                  "http://token:8080/extract",
			method:               http.MethodPost,
		},
		{
			name:     "ErrorDeleteToken",
			inToken:  mock.TokenTest,
			outCheck: true,
			isError:  true,
			url:      "http://db:8080/user",
			method:   http.MethodDelete,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infoServiceTest := service.InfoServices{
				DBHost:    mock.DBHostTest,
				DBPort:    mock.PortTest,
				TokenHost: mock.TokenHostTest,
				TokenPort: mock.PortTest,
				Secret:    mock.SecretTest,
			}

			var resultErr error
			var errorResponse string
			var mockHTTP *httpMock.MockClient

			if tt.isError {
				errorResponse = mock.ErrWebServer.Error()
			}

			responseJSON := fmt.Sprintf(`{
					"user":{
						"username":"username",
						"password":"password",
						"email":"email@email.com",
						"id":1
					},
					"id":1,
					"username":"usename",
					"email":"email@email.com",
					"check":%t
				}`, tt.outCheck)

			if tt.isError {
				mockHTTP = httpMock.NewMockClient(getIsErrorMock(
					tt.isErrorInsideRequest,
					newErrorHTTPComponets(tt.url, tt.method),
					responseJSON,
				))
			} else {
				mockHTTP = httpMock.NewMockClient(getMock(
					responseJSON,
				))
			}

			svc := service.NewService(
				mockHTTP,
				&infoServiceTest,
			)

			resultErr = svc.DeleteAccount(tt.inToken)

			if !tt.isError {
				if tt.outCheck {
					assert.Nil(t, resultErr)
				} else {
					assert.ErrorContains(t, resultErr, errorResponse)
				}
			} else {
				assert.ErrorContains(t, resultErr, errorResponse)
			}
		})
	}
}

func getMock(jsonResponse string) func(*http.Request) (*http.Response, error) {
	return func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			Body: io.NopCloser(strings.NewReader(jsonResponse)),
		}, nil
	}
}

//nolint:revive
func getIsErrorMock(
	isErrorInsideRequest bool,
	errorHTTPComponents errorHTTPComponents,
	jsonResponse string,
) func(*http.Request) (*http.Response, error) {
	return func(r *http.Request) (*http.Response, error) {
		if r.URL.String() == errorHTTPComponents.errorURL && r.Method == errorHTTPComponents.errorMethod {
			if isErrorInsideRequest {
				return &http.Response{
					Body: io.NopCloser(bytes.NewReader([]byte(`{
										"err":"error"
									}`),
					)),
				}, nil
			}

			return nil, mock.ErrWebServer
		}

		return &http.Response{
			Body: io.NopCloser(strings.NewReader(jsonResponse)),
		}, nil
	}
}
