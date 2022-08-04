package endpoint_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"app/internal/endpoint"
	"app/internal/entity"
	"app/internal/entity/mock"
	"app/internal/service"

	httpMock "app/internal/service/mock"

	"github.com/stretchr/testify/assert"
)

const (
	nameErrorRequest string = "ErrorRequest"
)

type incorrectRequest struct {
	incorrect bool
}

func TestSignUpEndpoint(t *testing.T) {
	t.Parallel()

	infoServiceTest := service.InfoServices{
		DBHost:    mock.DBHostTest,
		DBPort:    mock.PortTest,
		TokenHost: mock.TokenHostTest,
		TokenPort: mock.PortTest,
		Secret:    mock.SecretTest,
	}

	for _, tt := range []struct {
		name     string
		in       any
		outToken string
		outErr   string
	}{
		{
			name: mock.NameNoError,
			in: entity.UsernamePasswordEmailRequest{
				Username: mock.UsernameTest,
				Password: mock.PasswordTest,
				Email:    mock.EmailTest,
			},
			outToken: mock.TokenTest,
			outErr:   "",
		},
		{
			name: nameErrorRequest,
			in: incorrectRequest{
				incorrect: true,
			},
			outErr: "isn't of type",
		},
		{
			name:     "ErrorWebService",
			in:       entity.UsernamePasswordEmailRequest{},
			outToken: "",
			outErr:   service.ErrWebServer.Error(),
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			testResp := struct {
				Token string `json:"token"`
				Err   string `json:"err"`
				ID    int    `json:"id"`
			}{
				ID:    mock.IDTest,
				Token: tt.outToken,
				Err:   tt.outErr,
			}

			jsonData, err := json.Marshal(testResp)
			if err != nil {
				assert.Error(t, err)
			}

			mockClient := httpMock.NewMockClient(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(jsonData)),
				}, nil
			})

			svc := service.NewService(
				mockClient,
				&infoServiceTest,
			)

			r, err := endpoint.MakeSignUpEndpoint(svc)(context.TODO(), tt.in)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := r.(entity.TokenErrorResponse)
			if !ok {
				if tt.name != nameErrorRequest {
					assert.Error(t, mock.ErrNotTypeIndicated)
				}
			}

			if result.Err != "" {
				resultErr = result.Err
			}

			if tt.name == mock.NameNoError {
				assert.Empty(t, result.Err)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
			}

			assert.Equal(t, tt.outToken, result.Token)
		})
	}
}

func TestSignInEndpoint(t *testing.T) {
	t.Parallel()

	infoServiceTest := service.InfoServices{
		DBHost:    mock.DBHostTest,
		DBPort:    mock.PortTest,
		TokenHost: mock.TokenHostTest,
		TokenPort: mock.PortTest,
		Secret:    mock.SecretTest,
	}

	for _, tt := range []struct {
		name     string
		in       any
		outToken string
		outErr   string
	}{
		{
			name: mock.NameNoError,
			in: entity.UsernamePasswordRequest{
				Username: mock.UsernameTest,
				Password: mock.PasswordTest,
			},
			outToken: mock.TokenTest,
			outErr:   "",
		},
		{
			name: nameErrorRequest,
			in: incorrectRequest{
				incorrect: true,
			},
			outErr: "isn't of type",
		},
		{
			name:     "ErrorWebService",
			in:       entity.UsernamePasswordRequest{},
			outToken: "",
			outErr:   mock.ErrWebServer.Error(),
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			testResp := struct {
				Token string `json:"token"`
				Err   string `json:"err"`
				User  entity.User
			}{
				User: entity.User{
					ID:       mock.IDTest,
					Username: mock.UsernameTest,
					Password: mock.PasswordTest,
					Email:    mock.EmailTest,
				},
				Token: tt.outToken,
				Err:   tt.outErr,
			}

			jsonData, err := json.Marshal(testResp)
			if err != nil {
				assert.Error(t, err)
			}

			mockClient := httpMock.NewMockClient(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(jsonData)),
				}, nil
			})

			svc := service.NewService(
				mockClient,
				&infoServiceTest,
			)

			r, err := endpoint.MakeSignInEndpoint(svc)(context.TODO(), tt.in)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := r.(entity.TokenErrorResponse)
			if !ok {
				if tt.name != nameErrorRequest {
					assert.Error(t, mock.ErrNotTypeIndicated)
				}
			}

			if result.Err != "" {
				resultErr = result.Err
			}

			if tt.name == mock.NameNoError {
				assert.Empty(t, result.Err)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
			}

			assert.Equal(t, tt.outToken, result.Token)
		})
	}
}

func TestLogOutEndpoint(t *testing.T) {
	t.Parallel()

	infoServiceTest := service.InfoServices{
		DBHost:    mock.DBHostTest,
		DBPort:    mock.PortTest,
		TokenHost: mock.TokenHostTest,
		TokenPort: mock.PortTest,
		Secret:    mock.SecretTest,
	}

	for _, tt := range []struct {
		name   string
		in     any
		outErr string
	}{
		{
			name: mock.NameNoError,
			in: entity.Token{
				Token: mock.TokenTest,
			},
			outErr: "",
		},
		{
			name: nameErrorRequest,
			in: incorrectRequest{
				incorrect: true,
			},
			outErr: "isn't of type",
		},
		{
			name:   "ErrorWebService",
			in:     entity.Token{},
			outErr: mock.ErrWebServer.Error(),
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			testResp := struct {
				Err   string `json:"err"`
				Check bool   `json:"check"`
			}{
				Check: true,
				Err:   tt.outErr,
			}

			jsonData, err := json.Marshal(testResp)
			if err != nil {
				assert.Error(t, err)
			}

			mockClient := httpMock.NewMockClient(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(jsonData)),
				}, nil
			})

			svc := service.NewService(
				mockClient,
				&infoServiceTest,
			)

			r, err := endpoint.MakeLogOutEndpoint(svc)(context.TODO(), tt.in)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := r.(entity.ErrorResponse)
			if !ok {
				if tt.name != nameErrorRequest {
					assert.Error(t, mock.ErrNotTypeIndicated)
				}
			}

			if result.Err != "" {
				resultErr = result.Err
			}

			if tt.name == mock.NameNoError {
				assert.Empty(t, result.Err)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
			}
		})
	}
}

func TestGetAllUsersEndpoint(t *testing.T) {
	t.Parallel()

	infoServiceTest := service.InfoServices{
		DBHost:    mock.DBHostTest,
		DBPort:    mock.PortTest,
		TokenHost: mock.TokenHostTest,
		TokenPort: mock.PortTest,
		Secret:    mock.SecretTest,
	}

	for _, tt := range []struct {
		name     string
		outErr   string
		outUsers []entity.User
	}{
		{
			name: mock.NameNoError,
			outUsers: []entity.User{
				{
					ID:       mock.IDTest,
					Username: mock.UsernameTest,
					Password: mock.PasswordTest,
					Email:    mock.EmailTest,
				},
			},
			outErr: "",
		},
		{
			name:     "ErrorWebService",
			outUsers: nil,
			outErr:   mock.ErrWebServer.Error(),
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testResp := struct {
				Err   string        `json:"err"`
				Users []entity.User `json:"users"`
			}{
				Users: tt.outUsers,
				Err:   tt.outErr,
			}

			jsonData, err := json.Marshal(testResp)
			if err != nil {
				assert.Error(t, err)
			}

			mockClient := httpMock.NewMockClient(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(jsonData)),
				}, nil
			})

			svc := service.NewService(
				mockClient,
				&infoServiceTest,
			)

			r, err := endpoint.MakeGetAllUsersEndpoint(svc)(context.TODO(), nil)
			if err != nil {
				assert.Error(t, err)
			}

			result, ok := r.(entity.UsersErrorResponse)
			if !ok {
				assert.Error(t, mock.ErrNotTypeIndicated)
			}

			if tt.name == mock.NameNoError {
				assert.Empty(t, result.Err)
			} else {
				assert.Contains(t, result.Err, tt.outErr)
			}

			assert.Equal(t, tt.outUsers, result.Users)
		})
	}
}

func TestProfileEndpoint(t *testing.T) {
	t.Parallel()

	infoServiceTest := service.InfoServices{
		DBHost:    mock.DBHostTest,
		DBPort:    mock.PortTest,
		TokenHost: mock.TokenHostTest,
		TokenPort: mock.PortTest,
		Secret:    mock.SecretTest,
	}

	for _, tt := range []struct {
		in      any
		name    string
		outErr  string
		outUser entity.User
	}{
		{
			name: mock.NameNoError,
			in: entity.Token{
				Token: mock.TokenTest,
			},
			outUser: entity.User{
				ID:       mock.IDTest,
				Username: mock.UsernameTest,
				Password: mock.PasswordTest,
				Email:    mock.EmailTest,
			},
			outErr: "",
		},
		{
			name: nameErrorRequest,
			in: incorrectRequest{
				incorrect: true,
			},
			outErr: "isn't of type",
		},
		{
			name:    "ErrorWebService",
			in:      entity.Token{},
			outUser: entity.User{},
			outErr:  mock.ErrWebServer.Error(),
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			testResp := struct {
				Username string      `json:"username"`
				Email    string      `json:"email"`
				Err      string      `json:"err"`
				User     entity.User `json:"user"`
				ID       int         `json:"id"`
				Check    bool        `json:"check"`
			}{
				User:     tt.outUser,
				ID:       tt.outUser.ID,
				Username: tt.outUser.Username,
				Email:    tt.outUser.Email,
				Check:    true,
				Err:      tt.outErr,
			}

			jsonData, err := json.Marshal(testResp)
			if err != nil {
				assert.Error(t, err)
			}

			mockClient := httpMock.NewMockClient(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(jsonData)),
				}, nil
			})

			svc := service.NewService(
				mockClient,
				&infoServiceTest,
			)

			r, err := endpoint.MakeProfileEndpoint(svc)(context.TODO(), tt.in)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := r.(entity.UserErrorResponse)
			if !ok {
				if tt.name != nameErrorRequest {
					assert.Error(t, mock.ErrNotTypeIndicated)
				}
			}

			if result.Err != "" {
				resultErr = result.Err
			}

			if tt.name == mock.NameNoError {
				assert.Empty(t, result.Err)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
			}

			assert.Equal(t, tt.outUser, result.User)
		})
	}
}

func TestDeleteAccountEndpoint(t *testing.T) {
	t.Parallel()

	infoServiceTest := service.InfoServices{
		DBHost:    mock.DBHostTest,
		DBPort:    mock.PortTest,
		TokenHost: mock.TokenHostTest,
		TokenPort: mock.PortTest,
		Secret:    mock.SecretTest,
	}

	for _, tt := range []struct {
		name   string
		in     any
		outErr string
	}{
		{
			name: mock.NameNoError,
			in: entity.Token{
				Token: mock.TokenTest,
			},
			outErr: "",
		},
		{
			name: nameErrorRequest,
			in: incorrectRequest{
				incorrect: true,
			},
			outErr: "isn't of type",
		},
		{
			name:   "ErrorWebService",
			in:     entity.Token{},
			outErr: mock.ErrWebServer.Error(),
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			testResp := struct {
				Username string `json:"username"`
				Email    string `json:"email"`
				Err      string `json:"err"`
				ID       int    `json:"id"`
				Check    bool   `json:"check"`
			}{
				ID:       mock.IDTest,
				Username: mock.UsernameTest,
				Email:    mock.EmailTest,
				Check:    true,
				Err:      tt.outErr,
			}

			jsonData, err := json.Marshal(testResp)
			if err != nil {
				assert.Error(t, err)
			}

			mockClient := httpMock.NewMockClient(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(jsonData)),
				}, nil
			})

			svc := service.NewService(
				mockClient,
				&infoServiceTest,
			)

			r, err := endpoint.MakeDeleteAccountEndpoint(svc)(context.TODO(), tt.in)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := r.(entity.ErrorResponse)
			if !ok {
				if tt.name != nameErrorRequest {
					assert.Error(t, mock.ErrNotTypeIndicated)
				}
			}

			if result.Err != "" {
				resultErr = result.Err
			}

			if tt.name == mock.NameNoError {
				assert.Empty(t, result.Err)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
			}
		})
	}
}
