package endpoints_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AbdoAnss/go-fantasy-pl/client"
	"github.com/AbdoAnss/go-fantasy-pl/endpoints"
	"github.com/AbdoAnss/go-fantasy-pl/internal/cache"
	"github.com/stretchr/testify/require"
)

const liveTestEnv = "FPL_LIVE_TEST"

func skipUnlessLive(t *testing.T) {
	t.Helper()
	if os.Getenv(liveTestEnv) == "" {
		t.Skipf("set %s=1 to run live API tests", liveTestEnv)
	}
}

func newEndpointTestClient(t *testing.T) (*client.Client, *httptest.Server) {
	t.Helper()

	endpoints.SetSharedCache(cache.NewMemoryCache())

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.URL.Path == "/bootstrap-static/":
			writeTestdata(t, w, "bootstrap-static.json")
		case r.URL.Path == "/fixtures/":
			writeTestdata(t, w, "fixtures.json")
		case strings.HasPrefix(r.URL.Path, "/element-summary/"):
			id := strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/"), "/element-summary/")
			writeTestdata(t, w, fmt.Sprintf("element-summary-%s.json", id))
		default:
			http.NotFound(w, r)
		}
	}))

	c, err := client.NewClient(
		client.WithBaseURL(server.URL),
		client.WithMemoryCache(),
	)
	require.NoError(t, err)

	return c, server
}

func writeTestdata(t *testing.T, w http.ResponseWriter, name string) {
	t.Helper()

	path := filepath.Join("testdata", name)
	body, err := os.ReadFile(path)
	require.NoError(t, err)

	_, err = w.Write(body)
	require.NoError(t, err)
}
