package clouding

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFirewallID(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
			  "id": "LywOkvx5LWAp28NP",
			  "name": "Allow all private traffic",
			  "description": "Allow traffic from all private subnets",
			  "rules": [
			    {
			      "id": "Dd8v0nXJ1924rayY",
			      "description": "Allow TCP for private subnet 10.0.0.0/8",
			      "protocol": "tcp",
			      "portRangeMin": 1,
			      "portRangeMax": 65535,
			      "sourceIp": "10.0.0.0/8",
			      "enabled": true
			    },
			    {
			      "id": "N3V2ryXQjWa6pvok",
			      "description": "Allow TCP for private subnet 192.168.0.0/16",
			      "protocol": "tcp",
			      "portRangeMin": 1,
			      "portRangeMax": 65535,
			      "sourceIp": "192.168.0.0/16",
			      "enabled": true
			    },
			    {
			      "id": "2OM84qx6aWdz7JGr",
			      "description": "Allow TCP for private subnet 172.16.0.0/12",
			      "protocol": "tcp",
			      "portRangeMin": 1,
			      "portRangeMax": 65535,
			      "sourceIp": "172.16.0.0/12",
			      "enabled": true
			    }
			  ],
			  "attachments": [
			    {
			      "serverId": "Q7y1OZWlVn9mk6l3",
			      "serverName": "internal-server"
			    }
			  ]
			}`))
		if err != nil {
			t.Errorf("error writing response: %s", err)
		}

	}))
	client, err := NewAPI("token123", WithEndpoint(server.URL))
	if err != nil {
		t.Errorf("getting error calling NewAPI:%s", err)
	}
	firewall, err := client.GetFirewallID("LywOkvx5LWAp28NP")
	if err != nil {
		t.Errorf("getting error calling GetFirewallID: %s", err)
	}
	assert.Equal(t, "Allow all private traffic", firewall.Name)
	assert.Equal(t, "2OM84qx6aWdz7JGr", firewall.Rules[2].ID)
	assert.Equal(t, true, firewall.Rules[1].Enabled)
	assert.Equal(t, "Q7y1OZWlVn9mk6l3", firewall.Attachments[0].ServerID)
}

func TestCreateFirewall(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(`
			{
  				"id": "ZPlL0kxDYQ9Q3Yb5",
  				"name": "My firewall",
  				"description": "A firewall that restricts network accesses to my server",
  				"rules": [],
  				"attachments": []
			}
		`))
		if err != nil {
			t.Errorf("error writing response: %s", err)
		}
	}))

	client, err := NewAPI("token123", WithEndpoint(server.URL))
	if err != nil {
		t.Errorf("getting error calling NewAPI:%s", err)
	}

	firewall := Firewall{
		Name:        "My firewall",
		Description: "A firewall that restricts network accesses to my server",
	}
	err = client.CreateFirewall(&firewall)
	if err != nil {
		t.Errorf("getting error calling CreateFirewall: %s", err)
	}

	assert.Equal(t, "ZPlL0kxDYQ9Q3Yb5", firewall.ID)
}

func TestUpdateFirewall(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var firewall Firewall
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		err := json.NewDecoder(r.Body).Decode(&firewall)
		if err != nil {
			t.Errorf("error decoding firewall: %s", err)
		}
		assert.Equal(t, "the-new-name", firewall.NewName)
		assert.Equal(t, "The new description of the firewall", firewall.NewDescription)
	}))

	client, err := NewAPI("token123", WithEndpoint(server.URL))
	if err != nil {
		t.Errorf("getting error calling NewAPI:%s", err)
	}

	err = client.UpdateFirewall("ZPlL0kxDYQ9Q3Yb5", Firewall{
		NewName:        "the-new-name",
		NewDescription: "The new description of the firewall",
	})
	if err != nil {
		t.Errorf("getting error calling UpdateFirewall: %s", err)
	}
}

func TestUpdateFirewallWithError(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`
			{
  				"title": "One or more validation errors occurred.",
  				"status": 400,
  				"detail": "Please refer to the errors property for additional details.",
  				"instance": "/v1/firewalls/{id}",
  				"traceId": "00000000-0000-0000-0000-000000000000",
  				"errors": [
  				  {
  				    "propertyName": [
  				      "Validation error1.",
  				      "Validation error2."
  				    ]
  				  }
  				]
			}
		`))
		if err != nil {
			t.Errorf("error writing response: %s", err)
		}
	}))

	client, err := NewAPI("token123", WithEndpoint(server.URL))
	if err != nil {
		t.Errorf("getting error calling NewAPI:%s", err)
	}

	firewall := Firewall{
		Name:        "My firewall",
		Description: "A firewall that restricts network accesses to my server",
	}
	err = client.UpdateFirewall("ZPlL0kxDYQ9Q3Yb5", firewall)
	if err != nil {
		assert.Equal(t, "error updating firewall, status code: 400, title: One or more validation errors occurred.", err.Error())
	}
}

func TestDeleteFirewall(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}))
	client, err := NewAPI("token123", WithEndpoint(server.URL))
	if err != nil {
		t.Errorf("getting error calling NewAPI:%s", err)
	}
	err = client.DeleteFirewall("mYaRvlx1OmXApk6N")
	if err != nil {
		t.Errorf("getting error calling DeleteFirewall: %s", err)
	}
}

func TestDeleteFirewallWithError(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`
		{
  			"title": "One or more validation errors occurred.",
  			"status": 400,
  			"detail": "Please refer to the errors property for additional details.",
  			"instance": "/v1/firewalls/{id}",
  			"traceId": "00000000-0000-0000-0000-000000000000",
  			"errors": [
  			  {
  			    "propertyName": [
  			      "Validation error1.",
  			      "Validation error2."
  			    ]
  			  }
  			]
		}`))
		if err != nil {
			t.Errorf("error writing response: %s", err)
		}
	}))

	client, err := NewAPI("token123", WithEndpoint(server.URL))
	if err != nil {
		t.Errorf("getting error calling NewAPI:%s", err)
	}
	err = client.DeleteFirewall("mYaRvlx1OmXApk6N")
	if err != nil {
		assert.Equal(t, "error deleting firewall, status code: 400, title: One or more validation errors occurred.", err.Error())
	}
}

func TestGetFirewallRule(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`
		{
  			"firewallRule": {
  			  "id": "2OM84qx6aWdz7JGr",
  			  "description": "Allow TCP for private subnet 172.16.0.0/12",
  			  "protocol": "tcp",
  			  "portRangeMin": 1,
  			  "portRangeMax": 65535,
  			  "sourceIp": "172.16.0.0/12",
  			  "enabled": true
  			},
  			"firewallId": "m1LrZ3W8exDzN60o"
		}	
		`))
		if err != nil {
			t.Errorf("error writing response: %s", err)
		}
	}))

	client, err := NewAPI("token123", WithEndpoint(server.URL))
	if err != nil {
		t.Errorf("getting error calling NewAPI:%s", err)
	}
	firewallRule, err := client.GetFirewallRule("2OM84qx6aWdz7JGr")
	if err != nil {
		t.Errorf("getting error calling GetFirewallRule: %s", err)
	}

	assert.Equal(t, "2OM84qx6aWdz7JGr", firewallRule.FirewallRule.ID)
	assert.Equal(t, "Allow TCP for private subnet 172.16.0.0/12", firewallRule.FirewallRule.Description)
	assert.Equal(t, true, firewallRule.FirewallRule.Enabled)

}

func TestCreateFirewallRule(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var firewallRule FirewallRule
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err := json.NewDecoder(r.Body).Decode(&firewallRule)
		if err != nil {
			t.Errorf("error decoding firewallRule: %s", err)
		}
		assert.Equal(t, "10.0.0.0/8", firewallRule.SourceIP)
		_, err = w.Write([]byte(`
		{
			"id": "eAMVoaXqP9BLJwR6"
		}
		`))
		if err != nil {
			t.Errorf("error writing response: %s", err)
		}
	}))

	client, err := NewAPI("token123", WithEndpoint(server.URL))
	if err != nil {
		t.Errorf("getting error calling NewAPI:%s", err)
	}

	firewallRuleID := FirewallRuleID{
		FirewallID: "ZPlL0kxDYQ9Q3Yb5",
		FirewallRule: FirewallRule{
			SourceIP:     "10.0.0.0/8",
			Protocol:     "tcp",
			Description:  "Allow TCP for private subnet 10.0.0.0/8",
			PortRangeMin: 1,
			PortRangeMax: 65535,
		},
	}
	err = client.CreateFirewallRule(&firewallRuleID)
	if err != nil {
		t.Errorf("getting error calling CreateFirewallRule: %s", err)
	}
	assert.Equal(t, "eAMVoaXqP9BLJwR6", firewallRuleID.FirewallRule.ID)
}
