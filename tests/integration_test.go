//go:build integration

package nginx

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type response struct {
	Ltp []ltp `json:"ltp"`
}

type ltp struct {
	Pair   string  `json:"pair"`
	Amount float64 `json:"amount"`
}

type appContainer struct {
	testcontainers.Container
	URI string
}

func startContainer(ctx context.Context) (*appContainer, error) {
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "..",
			Dockerfile: "Dockerfile",
		},
		ExposedPorts: []string{"8081/tcp"},
		WaitingFor:   wait.ForHTTP("/_healthz").WithStartupTimeout(10 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return nil, err //nolint: wrapcheck
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err //nolint: wrapcheck
	}

	mappedPort, err := container.MappedPort(ctx, "8081")
	if err != nil {
		return nil, err //nolint: wrapcheck
	}

	uri := "http://" + ip + ":" + mappedPort.Port()

	return &appContainer{Container: container, URI: uri}, nil
}

func TestIntegration(t *testing.T) { //nolint: funlen, paralleltest
	app, err := startContainer(context.Background())
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}
	defer app.Terminate(context.Background()) //nolint: errcheck

	// Test the application
	t.Run("Health check", func(t *testing.T) { //nolint: paralleltest
		resp, err := http.Get(app.URI + "/_healthz") //nolint: noctx
		if err != nil {
			t.Fatalf("failed to get health from container: %s", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Test correct query", func(t *testing.T) { //nolint: paralleltest
		resp, err := http.Get(app.URI + "/api/v1/ltp") //nolint: noctx
		if err != nil {
			t.Fatalf("failed to get response from container: %s", err)
		}

		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		var r response

		err = json.Unmarshal(body, &r)
		require.NoError(t, err)

		assert.Len(t, r.Ltp, 3)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Test correct query for one pair", func(t *testing.T) { //nolint: paralleltest
		resp, err := http.Get(app.URI + "/api/v1/ltp?pairs=BTC/USD") //nolint: noctx
		if err != nil {
			t.Fatalf("failed to get response from container: %s", err)
		}

		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		var r response

		err = json.Unmarshal(body, &r)
		require.NoError(t, err)

		assert.Len(t, r.Ltp, 1)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
	t.Run("Test correct query for one pair", func(t *testing.T) { //nolint: paralleltest
		resp, err := http.Get(app.URI + "/api/v1/ltp?pairs=BTC/USD,BTC/CHF") //nolint: noctx
		if err != nil {
			t.Fatalf("failed to get response from container: %s", err)
		}

		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		var r response

		err = json.Unmarshal(body, &r)
		require.NoError(t, err)

		assert.Len(t, r.Ltp, 2)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
	t.Run("Test bad pair", func(t *testing.T) { //nolint: paralleltest
		resp, err := http.Get(app.URI + "/api/v1/ltp?pairs=BADPAIR") //nolint: noctx
		if err != nil {
			t.Fatalf("failed to get response from container: %s", err)
		}

		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
