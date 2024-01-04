package httput

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestHTTPClient demonstrate how to use httptest.NewServer() to reply a mocked response to a client.
func TestHTTPClient(t *testing.T) {
	// Using httptest NewServer we get a new Server instance from which we can mock responses.
	expected := "dummy data"
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	}))
	// Server must be closed when the test ends.
	t.Cleanup(func() { svr.Close() })

	// Client call would be the piece of logic replaced by your service call.
	out, err := httpClientCall(svr.URL)
	require.NoError(t, err)

	require.Equal(t, "dummy data", string(out))
}

// httpClientCall simulates the business logic of an existing service calling another system HTTP endpoint.
func httpClientCall(url string) (string, error) {
	res, err := http.Get(url + "/upper?word=anything")
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	out, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
