package mockutils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	server := httptest.NewServer(Handler(t, []Request{
		{
			Method: "GET", Path: "/",
			Status: 200,
			JSON: struct {
				Data string `json:"data"`
			}{
				Data: "Hello",
			},
		},
		{
			Method: "GET", Path: "/",
			Status:  400,
			JSONRaw: `{"error": "failed"}`,
		},
	}))
	defer server.Close()

	// Request 1
	resp, err := http.Get(server.URL)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, `{"data":"Hello"}`, readBody(t, resp))

	// Request 2
	resp, err = http.Get(server.URL)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, `{"error": "failed"}`, readBody(t, resp))
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	return strings.TrimSuffix(string(body), "\n")
}