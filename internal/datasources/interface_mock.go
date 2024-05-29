package datasources

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// MockingClient Mock for the HTTPClient.
type MockClient struct {
	mock.Mock
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)

	return args.Get(0).(*http.Response), args.Error(1) //nolint: forcetypeassert, wrapcheck
}

// Type check.
var _ HTTPClient = new(MockClient)
