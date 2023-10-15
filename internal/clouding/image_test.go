package clouding

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetImageID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`
		{
  			"id": "d3mKbx4zd3XEQaqP",
  			"name": "Quantum Nova OS (English 64Bit)",
  			"minimumSizeGb": 25,
  			"accessMethods": {
  			  "sshKey": "required-with-private-key",
  			  "password": "not-supported"
  			},
  			"pricePerHour": 0.00684,
  			"pricePerMonthApprox": 4.9932,
  			"billingUnit": "Core"
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
	image, err := client.GetImageID("d3mKbx4zd3XEQaqP")
	if err != nil {
		t.Errorf("getting error calling GetImageID: %s", err)
	}

	assert.Equal(t, "d3mKbx4zd3XEQaqP", image.ID)
	assert.Equal(t, "Quantum Nova OS (English 64Bit)", image.Name)
	assert.Equal(t, int64(25), image.MinimumSizeGB)
	assert.Equal(t, "required-with-private-key", image.AccessMethods.SshKey)
	assert.Equal(t, "not-supported", image.AccessMethods.Password)
	assert.Equal(t, 0.00684, image.PricePerHour)
	assert.Equal(t, 4.9932, image.PricePerMonthApprox)
	assert.Equal(t, "Core", image.BillingUnit)
}
