package clouding

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetActionID(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "N3V2ryXQjWa6pvok",
			"status": "completed",
			"type": "delete",
			"startedAt": "2023-01-03T11:55:00.0000000Z",
			"completedAt": "2023-01-03T11:55:05.0000000Z",
			"resourceId": "m1LrZ3W8exDzN60o",
			"resourceType": "server"
		}`))
	}))

	client, err := NewAPI("token123", WithEndpoint(server.URL))
	if err != nil {
		t.Errorf("getting error calling NewAPI:%s", err)
	}

	action, err := client.GetAction("N3V2ryXQjWa6pvok")
	if err != nil {
		t.Errorf("getting error calling GetAction: %s", err)
	}

	assert.Equal(t, "N3V2ryXQjWa6pvok", action.ID)
	assert.Equal(t, "completed", action.Status)
	assert.Equal(t, "delete", action.Type)
	assert.Equal(t, "2023-01-03T11:55:00.0000000Z", action.StartedAt)
	assert.Equal(t, "2023-01-03T11:55:05.0000000Z", action.CompletedAt)
	assert.Equal(t, "m1LrZ3W8exDzN60o", action.ResourceID)
	assert.Equal(t, "server", action.ResourceType)
}
