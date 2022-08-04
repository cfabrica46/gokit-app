package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"app/internal/entity"

	httptransport "github.com/go-kit/kit/transport/http"
)

var errFailedGetHeader = errors.New("failed to get header")

// DecodeRequestWithoutBody ...
func DecodeRequestWithoutBody() httptransport.DecodeRequestFunc {
	return func(_ context.Context, _ *http.Request) (any, error) {
		var request entity.EmptyRequest

		return request, nil
	}
}

// DecodeRequest ...
func DecodeRequestWithBody[req entity.UsernamePasswordEmailRequest |
	entity.UsernamePasswordRequest](request req,
) httptransport.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (any, error) {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return nil, fmt.Errorf("failed to decode request: %w", err)
		}

		return request, nil
	}
}

func DecodeRequestWithHeader(request entity.Token) httptransport.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (any, error) {
		if r.Header.Get("Authorization") == "" {
			return nil, errFailedGetHeader
		}

		request.Token = r.Header.Get("Authorization")

		return request, nil
	}
}

// EncodeResponse ...
func EncodeResponse(_ context.Context, w http.ResponseWriter, response any) (err error) {
	if err = json.NewEncoder(w).Encode(response); err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	return nil
}
