package petition_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"app/internal/entity"
	"app/internal/entity/mock"
	"app/internal/petition"
	"app/internal/service"

	serviceMock "app/internal/service/mock"

	"github.com/stretchr/testify/assert"
)

func TestRequestFunc(t *testing.T) {
	t.Parallel()

	mockOK := serviceMock.NewMockClient(func(_ *http.Request) (*http.Response, error) {
		response := &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(strings.NewReader(`{
				"id": 1
			}`)),
		}

		return response, nil
	})

	mockNotOK := serviceMock.NewMockClient(func(_ *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("%w: error", service.ErrWebServer)
	})

	for _, tt := range []struct {
		Response *entity.IDErrorResponse
		client   petition.HTTPClient
		body     any
		name     string
		inURL    string
		inMethod string
		outErr   string
	}{
		{
			name:   "NoError",
			client: mockOK,
			body: entity.UsernameRequest{
				Username: mock.UsernameTest,
			},
			inURL:    "localhost:8080",
			inMethod: http.MethodPost,
			Response: &entity.IDErrorResponse{},
			outErr:   "",
		},
		{
			name:     "ErrorMarshal",
			client:   mockOK,
			body:     func() {},
			inURL:    "localhost:8080",
			inMethod: http.MethodPost,
			Response: &entity.IDErrorResponse{},
			outErr:   "error to make petition",
		},
		{
			name:   "ErrorURL",
			client: mockOK,
			body: entity.UsernameRequest{
				Username: mock.UsernameTest,
			},
			inURL:    "%%",
			inMethod: http.MethodPost,
			Response: &entity.IDErrorResponse{},
			outErr:   "error to make petition",
		},
		{
			name:   "ErrorService",
			client: mockNotOK,
			body: entity.UsernameRequest{
				Username: mock.UsernameTest,
			},
			inURL:    "localhost:8080",
			inMethod: http.MethodPost,
			Response: &entity.IDErrorResponse{},
			outErr:   "error to make petition",
		},
		{
			name:   "ErrorDecode",
			client: mockOK,
			body: entity.UsernameRequest{
				Username: mock.UsernameTest,
			},
			inURL:    "localhost:8080",
			inMethod: http.MethodPost,
			Response: nil,
			outErr:   "failed to decode request",
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := petition.RequestFunc(
				tt.client,
				tt.body,
				petition.NewHTTPComponents(
					tt.inURL,
					tt.inMethod,
				),
				tt.Response,
			)

			if tt.name == "NoError" {
				assert.Nil(t, err)
				assert.Equal(t, mock.IDTest, tt.Response.ID)
				assert.Zero(t, tt.Response.Err)
			} else {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tt.outErr)
			}
		})
	}
}

func TestRequestFuncWithoutBody(t *testing.T) {
	t.Parallel()

	mockOK := serviceMock.NewMockClient(func(_ *http.Request) (*http.Response, error) {
		response := &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(strings.NewReader(`{
				"users": [
					{
						"id":1,
						"username":"username",
						"password":"password",
						"email":"email@email.com"
					}
				]
			}`)),
		}

		return response, nil
	})

	mockNotOK := serviceMock.NewMockClient(func(_ *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("%w: error", service.ErrWebServer)
	})

	for _, tt := range []struct {
		Response *entity.UsersErrorResponse
		client   petition.HTTPClient
		name     string
		inURL    string
		inMethod string
		outErr   string
	}{
		{
			name:     "NoError",
			client:   mockOK,
			inURL:    "localhost:8080",
			inMethod: http.MethodPost,
			Response: &entity.UsersErrorResponse{},
			outErr:   "",
		},
		{
			name:     "ErrorURL",
			client:   mockOK,
			inURL:    "%%",
			inMethod: http.MethodPost,
			Response: &entity.UsersErrorResponse{},
			outErr:   "error to make petition",
		},
		{
			name:     "ErrorService",
			client:   mockNotOK,
			inURL:    "localhost:8080",
			inMethod: http.MethodPost,
			Response: &entity.UsersErrorResponse{},
			outErr:   "error to make petition",
		},
		{
			name:     "ErrorDecode",
			client:   mockOK,
			inURL:    "localhost:8080",
			inMethod: http.MethodPost,
			Response: nil,
			outErr:   "failed to decode request",
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := petition.RequestFuncWithoutBody(
				tt.client,
				petition.NewHTTPComponents(
					tt.inURL,
					tt.inMethod,
				),
				tt.Response,
			)

			if tt.name == "NoError" {
				assert.Nil(t, err)
				assert.NotNil(t, tt.Response.Users)
				assert.Zero(t, tt.Response.Err)
			} else {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, tt.outErr)
			}
		})
	}
}
