package endpoint

import (
	"context"
	"errors"
	"fmt"

	"app/internal/entity"
	"app/internal/service"

	"github.com/go-kit/kit/endpoint"
)

var ErrRequest = errors.New("error to request")

// MakeSignUpEndpoint ...
func MakeSignUpEndpoint(svc service.Service) endpoint.Endpoint {
	return func(_ context.Context, request any) (any, error) {
		var errMessage string

		req, ok := request.(entity.UsernamePasswordEmailRequest)
		if !ok {
			return nil, fmt.Errorf("%w: isn't of type GenerateTokenRequest", ErrRequest)
		}

		token, err := svc.SignUp(req.Username, req.Password, req.Email)
		if err != nil {
			errMessage = err.Error()
		}

		return entity.TokenErrorResponse{Token: token, Err: errMessage}, nil
	}
}

// MakeSignInEndpoint ...
func MakeSignInEndpoint(svc service.Service) endpoint.Endpoint {
	return func(_ context.Context, request any) (any, error) {
		var errMessage string

		req, ok := request.(entity.UsernamePasswordRequest)
		if !ok {
			return nil, fmt.Errorf("%w: isn't of type GenerateTokenRequest", ErrRequest)
		}

		token, err := svc.SignIn(req.Username, req.Password)
		if err != nil {
			errMessage = err.Error()
		}

		return entity.TokenErrorResponse{Token: token, Err: errMessage}, nil
	}
}

// MakeLogOutEndpoint ...
func MakeLogOutEndpoint(svc service.Service) endpoint.Endpoint {
	return func(_ context.Context, request any) (any, error) {
		var errMessage string

		req, ok := request.(entity.Token)
		if !ok {
			return nil, fmt.Errorf("%w: isn't of type GenerateTokenRequest", ErrRequest)
		}

		err := svc.LogOut(req.Token)
		if err != nil {
			errMessage = err.Error()
		}

		return entity.ErrorResponse{Err: errMessage}, nil
	}
}

// MakeGetAllUsersEndpoint ...
func MakeGetAllUsersEndpoint(svc service.Service) endpoint.Endpoint {
	return func(_ context.Context, _ any) (any, error) {
		var errMessage string

		users, err := svc.GetAllUsers()
		if err != nil {
			errMessage = err.Error()
		}

		return entity.UsersErrorResponse{Users: users, Err: errMessage}, nil
	}
}

// MakeProfileEndpoint ...
func MakeProfileEndpoint(svc service.Service) endpoint.Endpoint {
	return func(_ context.Context, request any) (any, error) {
		var errMessage string

		req, ok := request.(entity.Token)
		if !ok {
			return nil, fmt.Errorf("%w: isn't of type GenerateTokenRequest", ErrRequest)
		}

		user, err := svc.Profile(req.Token)
		if err != nil {
			errMessage = err.Error()
		}

		return entity.UserErrorResponse{User: user, Err: errMessage}, nil
	}
}

// MakeDeleteAccountEndpoint ...
func MakeDeleteAccountEndpoint(svc service.Service) endpoint.Endpoint {
	return func(_ context.Context, request any) (any, error) {
		var errMessage string

		req, ok := request.(entity.Token)
		if !ok {
			return nil, fmt.Errorf("%w: isn't of type GenerateTokenRequest", ErrRequest)
		}

		err := svc.DeleteAccount(req.Token)
		if err != nil {
			errMessage = err.Error()
		}

		return entity.ErrorResponse{Err: errMessage}, nil
	}
}
