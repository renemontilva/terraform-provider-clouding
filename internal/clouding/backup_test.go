package clouding

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBackupID(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
		{
		  "id": "86EAL1xB769Z4q2w",
		  "createdAt": "2023-01-01T12:00:00.0000000Z",
		  "serverId": "mawqYZWOojWQyOV0",
		  "serverName": "my-test-server",
		  "volumeSizeGb": 25,
		  "image": {
		    "id": "lo1qJ9oZb1xGMEgD",
		    "name": "CelestiaOS 2.04 (64 Bit)",
		    "accessMethods": {
		      "sshKey": "optional",
		      "password": "optional"
		    }
		  }
		}
		`))
	}))

	client, err := NewAPI("token123", WithEndpoint(server.URL))
	if err != nil {
		t.Errorf("getting error creating NewAPI: %s", err)
	}

	backup, err := client.GetBackupID("86EAL1xB769Z4q2w")
	if err != nil {
		t.Errorf("getting error calling GetBackupID: %s", err)
	}

	assert.Equal(t, "86EAL1xB769Z4q2w", backup.ID)
	assert.Equal(t, "2023-01-01T12:00:00.0000000Z", backup.CreatedAt)
	assert.Equal(t, "mawqYZWOojWQyOV0", backup.ServerID)
	assert.Equal(t, "my-test-server", backup.ServerName)
	assert.Equal(t, 25, backup.VolumeSizeGb)
	assert.Equal(t, "lo1qJ9oZb1xGMEgD", backup.Image.ID)
	assert.Equal(t, "CelestiaOS 2.04 (64 Bit)", backup.Image.Name)
	assert.Equal(t, "optional", backup.Image.AccessMethods.SshKey)
	assert.Equal(t, "optional", backup.Image.AccessMethods.Password)

}
