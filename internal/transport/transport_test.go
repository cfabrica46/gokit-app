package transport_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/internal/entity/mock"

	"github.com/cfabrica46/gokit-crud/app/service"
	"github.com/stretchr/testify/assert"
)

const (
	//nolint:gosec
	usernamePasswordEmailRequestJSON = `{
		 "username": "username",
		 "password": "password",
		 "email": "email@email.com"
	}`

	//nolint:gosec
	usernamePasswordRequestJSON = `{
		 "username": "username",
		 "password": "password"
	}`
)

func TestDecodeRequestWithoutBody(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		out    any
		in     *http.Request
		name   string
		outErr string
	}{
		{
			name:   mock.NameNoError + "GetAllRequest",
			in:     nil,
			outErr: "",
			out:    service.EmptyRequest{},
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, err := service.DecodeRequestWithoutBody()(context.TODO(), tt.in)

			assert.Empty(t, err)
			assert.Equal(t, tt.out, r)
		})
	}
}

func TestDecodeRequestWithBody(t *testing.T) {
	t.Parallel()

	usernamePasswordEmailReq, err := http.NewRequest(
		http.MethodPost,
		mock.URLTest,
		bytes.NewBuffer([]byte(usernamePasswordEmailRequestJSON)),
	)
	if err != nil {
		assert.Error(t, err)
	}

	usernamePasswordReq, err := http.NewRequest(
		http.MethodPost,
		mock.URLTest,
		bytes.NewBuffer([]byte(usernamePasswordRequestJSON)),
	)
	if err != nil {
		assert.Error(t, err)
	}

	badReq, err := http.NewRequest(http.MethodPost, mock.URLTest, bytes.NewBuffer([]byte{}))
	if err != nil {
		assert.Error(t, err)
	}

	for _, tt := range []struct {
		inType      any
		in          *http.Request
		name        string
		outUsername string
		outPassword string
		outEmail    string
		outErr      string
		outID       int
	}{
		{
			name:        mock.NameNoError + "UsernamePasswordEmailRequest",
			inType:      service.UsernamePasswordEmailRequest{},
			in:          usernamePasswordEmailReq,
			outUsername: mock.UsernameTest,
			outPassword: mock.PasswordTest,
			outEmail:    mock.EmailTest,
			outErr:      "",
		},
		{
			name:        mock.NameNoError + "UsernamePasswordRequest",
			inType:      service.UsernamePasswordRequest{},
			in:          usernamePasswordReq,
			outUsername: mock.UsernameTest,
			outPassword: mock.PasswordTest,
			outErr:      "",
		},
		{
			name:   "BadRequest",
			inType: service.UsernamePasswordEmailRequest{},
			in:     badReq,
			outErr: "EOF",
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string

			var req any

			switch resultType := tt.inType.(type) {
			case service.UsernamePasswordEmailRequest:
				req, err = service.DecodeRequestWithBody(resultType)(context.TODO(), tt.in)
				if err != nil {
					resultErr = err.Error()
				}

				result, ok := req.(service.UsernamePasswordEmailRequest)
				if ok {
					assert.Equal(t, tt.outUsername, result.Username)
					assert.Equal(t, tt.outPassword, result.Password)
					assert.Equal(t, tt.outEmail, result.Email)
					assert.Contains(t, resultErr, tt.outErr)
				} else {
					assert.NotNil(t, err)
				}

			case service.UsernamePasswordRequest:
				req, err = service.DecodeRequestWithBody(resultType)(context.TODO(), tt.in)
				if err != nil {
					resultErr = err.Error()
				}

				result, ok := req.(service.UsernamePasswordRequest)
				if ok {
					assert.Equal(t, tt.outUsername, result.Username)
					assert.Equal(t, tt.outPassword, result.Password)
					assert.Contains(t, resultErr, tt.outErr)
				} else {
					assert.NotNil(t, err)
				}

			default:
				assert.Fail(t, "Error to type inType")
			}
		})
	}
}

func TestDecodeRequestWithHeader(t *testing.T) {
	t.Parallel()

	okReq, err := http.NewRequest(
		http.MethodPost,
		mock.URLTest,
		bytes.NewBuffer([]byte{}),
	)
	if err != nil {
		assert.Error(t, err)
	}

	okReq.Header.Set("Authorization", "token")

	badReq, err := http.NewRequest(http.MethodPost, mock.URLTest, bytes.NewBuffer([]byte{}))
	if err != nil {
		assert.Error(t, err)
	}

	for _, tt := range []struct {
		inType   service.TokenRequest
		in       *http.Request
		name     string
		outErr   string
		outToken string
		outID    int
	}{
		{
			name:     mock.NameNoError,
			inType:   service.TokenRequest{},
			in:       okReq,
			outToken: mock.TokenTest,
			outErr:   "",
		},
		{
			name:   "BadRequest",
			inType: service.TokenRequest{},
			in:     badReq,
			outErr: "failed to get header",
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var resultErr string
			var r any

			r, err = service.DecodeRequestWithHeader(tt.inType)(context.TODO(), tt.in)
			if err != nil {
				resultErr = err.Error()
			}

			result, ok := r.(service.TokenRequest)
			if tt.name == mock.NameNoError {
				if !ok {
					assert.Fail(t, "Error to type inType")
				}
			}

			if tt.name == mock.NameNoError {
				assert.Equal(t, tt.outToken, result.Token)
				assert.Nil(t, err)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
			}
		})
	}
}

func TestEncodeResponse(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name   string
		in     any
		outErr string
	}{
		{
			name:   mock.NameNoError,
			in:     "test",
			outErr: "",
		},
		{
			name:   "ErrorBadEncode",
			in:     func() {},
			outErr: "json: unsupported type: func()",
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var resultErr string

			err := service.EncodeResponse(context.TODO(), httptest.NewRecorder(), tt.in)
			if err != nil {
				resultErr = err.Error()
			}

			if tt.name == mock.NameNoError {
				assert.Empty(t, resultErr)
			} else {
				assert.Contains(t, resultErr, tt.outErr)
			}
		})
	}
}
