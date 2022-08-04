package service

import (
	"net/http"
)

type myDoFunc func(req *http.Request) (*http.Response, error)

// MockClient ...
type MockClient struct {
	doFunc myDoFunc
}

// NewMockClient ...
func NewMockClient(d myDoFunc) *MockClient {
	return &MockClient{d}
}

// Do ...
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}
