package datasources

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"

	"github.com/rmarku/ltp_api/internal/entities"
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

// MockDataSource Mock for the datasources.
type MockDataSource struct {
	mock.Mock
}

func (m *MockDataSource) GetData(ctx context.Context, pair string) (*entities.LTP, error) {
	args := m.Called(ctx, pair)

	return args.Get(0).(*entities.LTP), args.Error(1) //nolint: forcetypeassert, wrapcheck
}

var _ DataSource = new(MockDataSource)
