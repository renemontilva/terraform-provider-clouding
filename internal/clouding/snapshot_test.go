package clouding

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSnapShotID(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`
		{
			"id": "jDGPRJXLpGXeV5M1",
			"sizeGb": 15,
			"name": "snapshot-with-mysql",
			"description": "A snapshot of the server after mysql installation",
			"createdAt": "2023-01-03T12:00:00.0000000Z",
			"sourceServerName": "db-server",
			"image": {
			  "id": "DGPRJXLAOGWeV5M1",
			  "name": "Galaxion 15 (CelestiaOS 2.04 64Bit)",
			  "accessMethods": {
				"sshKey": "optional",
				"password": "optional"
			  }
			},
			"cost": {
			  "pricePerHour": 0.0021,
			  "pricePerMonthApprox": 1.533
			}
		  }
		`))
		if err != nil {
			t.Errorf("error writing response: %s", err)
		}
	}))

	client, err := NewAPI("token123", WithEndpoint(server.URL))
	if err != nil {
		t.Errorf("getting error creating NewAPI: %s", err)
	}

	snapshot, err := client.GetSnapshotID("jDGPRJXLpGXeV5M1")
	if err != nil {
		t.Errorf("getting error calling GetSnapshotID: %s", err)
	}

	assert.Equal(t, "jDGPRJXLpGXeV5M1", snapshot.ID)
	assert.Equal(t, "snapshot-with-mysql", snapshot.Name)
	assert.Equal(t, int64(15), snapshot.SizeGb)
	assert.Equal(t, "A snapshot of the server after mysql installation", snapshot.Description)
	assert.Equal(t, "2023-01-03T12:00:00.0000000Z", snapshot.CreatedAt)
	assert.Equal(t, "db-server", snapshot.SourceServeName)
	assert.Equal(t, "DGPRJXLAOGWeV5M1", snapshot.Image.ID)
	assert.Equal(t, "Galaxion 15 (CelestiaOS 2.04 64Bit)", snapshot.Image.Name)
	assert.Equal(t, "optional", snapshot.Image.AccessMethods.SshKey)
	assert.Equal(t, "optional", snapshot.Image.AccessMethods.Password)
	assert.Equal(t, 0.0021, snapshot.Cost.PricePerHour)
	assert.Equal(t, 1.533, snapshot.Cost.PricePerMonthApprox)
}
