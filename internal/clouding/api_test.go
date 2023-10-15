package clouding

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendRequest(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc          string
		tokenExpected string
	}{
		{
			desc:          "Test sendRequest token",
			tokenExpected: "token123",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			responseHeader := http.StatusOK
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("X-API-KEY") != tC.tokenExpected {
					responseHeader = http.StatusBadRequest
				}
				w.WriteHeader(responseHeader)
			}))
			defer server.Close()
			client, err := NewAPI("token123", WithEndpoint(server.URL))
			if err != nil {
				t.Errorf("getting error creating NewAPI: %s", err)
			}
			response, err := client.sendRequest(http.MethodGet, "firewall", []byte{})
			if err != nil {
				t.Errorf("getting error sending request to %s: %s", server.URL, err)
			}
			assert.Equal(t, http.StatusOK, response.StatusCode)
		})
	}
}
