package datasources

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		mockError    error
		mockResponse *http.Response
		expected     *LTP
		name         string
		pair         string
		expectError  bool
	}{
		{
			name: "Successful response",
			pair: "XBTUSD",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(
					bytes.NewBufferString(`{
						"result": {"XBTUSD": {"c": ["50000.0"]}},
						"error": []
					}`),
				),
			},
			expected:    &LTP{Pair: "XBTUSD", Amount: 50000.0},
			expectError: false,
		},
		{
			name: "Kraken returned error",
			pair: "XBTUSD",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(
					bytes.NewBufferString(`{
						"result": {},
						"error": ["Some error"]
					}`),
				),
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Invalid JSON response",
			pair: "XBTUSD",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{invalid json}`)),
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "HTTP client error",
			pair: "XBTUSD",
			mockResponse: &http.Response{
				StatusCode: http.StatusInternalServerError,
			},
			mockError:   errors.New("HTTP client error"), //nolint: goerr113
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := new(MockClient)
			k := &DataSourceKraken{
				uri:    "https://mockedKraken.example.com/api",
				client: mockClient,
			}

			req, _ := http.NewRequest(http.MethodGet, k.uri, nil)
			q := req.URL.Query()
			q.Add("pair", tt.pair)
			req.URL.RawQuery = q.Encode()

			mockClient.On("Do", req).Return(tt.mockResponse, tt.mockError)

			ctx := context.Background()
			result, err := k.GetData(ctx, tt.pair)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expected, result)
			mockClient.AssertExpectations(t)
		})
	}
}
