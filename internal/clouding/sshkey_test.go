package clouding

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSshKeyID(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
		{
			"id": "Dd8v0nXJ1924rayY",
			"name": "my-ssh-key-without-private-key",
			"fingerprint": "a6:92:c9:0a:c4:ca:7d:2c:fb:98:42:63:79:37:a8:24",
			"publicKey": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDEc2hHw894iG/BULKSwdJceqLG8qGoifKfCVNTkFJgILn7dnJg7ioWtzpztaWCiwJe7R2zVF68ANi0+Ie8UM70/MFk5WKB1mOC7k85HCT8jXRM4F5zhnwYKLmVkGiYQHUchp2RgmKNxdcoJNOfYnvZL0dWZvnrl/M0J3vTmw/xXmqqsbjBI1v5EXvjCYiwyLgZ6ZIcOyDgYFbZj6/YAVpQEaFzmwZSlOFL8AC7YTh2uFYu1hxKOOOVJrAd4/kBXUnd1SqOyj4lPlQ7RMTym9viFkfUisIZOFWQMCIPJrrc4qba5qjHmErpamFGGC7QuFa0bBs6Newq3Qj4mUWBSYC7",
			"hasPrivateKey": false
		}
		`))
	}))

	client, err := NewAPI("token123", WithEndpoint(server.URL))
	if err != nil {
		t.Errorf("getting error creating NewAPI: %s", err)
	}

	sshkey, err := client.GetSshKeyID("Dd8v0nXJ1924rayY")
	if err != nil {
		t.Errorf("getting error calling GetSshKeyID: %s", err)
	}

	assert.Equal(t, "Dd8v0nXJ1924rayY", sshkey.ID)
	assert.Equal(t, "my-ssh-key-without-private-key", sshkey.Name)
	assert.Equal(t, "a6:92:c9:0a:c4:ca:7d:2c:fb:98:42:63:79:37:a8:24", sshkey.Fingerprint)
	assert.Equal(t, "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDEc2hHw894iG/BULKSwdJceqLG8qGoifKfCVNTkFJgILn7dnJg7ioWtzpztaWCiwJe7R2zVF68ANi0+Ie8UM70/MFk5WKB1mOC7k85HCT8jXRM4F5zhnwYKLmVkGiYQHUchp2RgmKNxdcoJNOfYnvZL0dWZvnrl/M0J3vTmw/xXmqqsbjBI1v5EXvjCYiwyLgZ6ZIcOyDgYFbZj6/YAVpQEaFzmwZSlOFL8AC7YTh2uFYu1hxKOOOVJrAd4/kBXUnd1SqOyj4lPlQ7RMTym9viFkfUisIZOFWQMCIPJrrc4qba5qjHmErpamFGGC7QuFa0bBs6Newq3Qj4mUWBSYC7", sshkey.PublicKey)
	assert.Equal(t, false, sshkey.HasPrivateKey)
}

func TestCreateSshKey(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`
		{
			"id": "Dd8v0nXJ1924rayY",
			"name": "my-ssh-key-without-private-key",
			"fingerprint": "a6:92:c9:0a:c4:ca:7d:2c:fb:98:42:63:79:37:a8:24",
			"publicKey": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDEc2hHw894iG/BULKSwdJceqLG8qGoifKfCVNTkFJgILn7dnJg7ioWtzpztaWCiwJe7R2zVF68ANi0+Ie8UM70/MFk5WKB1mOC7k85HCT8jXRM4F5zhnwYKLmVkGiYQHUchp2RgmKNxdcoJNOfYnvZL0dWZvnrl/M0J3vTmw/xXmqqsbjBI1v5EXvjCYiwyLgZ6ZIcOyDgYFbZj6/YAVpQEaFzmwZSlOFL8AC7YTh2uFYu1hxKOOOVJrAd4/kBXUnd1SqOyj4lPlQ7RMTym9viFkfUisIZOFWQMCIPJrrc4qba5qjHmErpamFGGC7QuFa0bBs6Newq3Qj4mUWBSYC7",
			"hasPrivateKey": false
		}
		`))
	}))

	client, err := NewAPI("token123", WithEndpoint(server.URL))
	if err != nil {
		t.Errorf("getting error creating NewAPI: %s", err)
	}

	sshkey := SshKey{
		Name:          "my-ssh-key-without-private-key",
		PublicKey:     "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDEc2hHw894iG/BULKSwdJceqLG8qGoifKfCVNTkFJgILn7dnJg7ioWtzpztaWCiwJe7R2zVF68ANi0+Ie8UM70/MFk5WKB1mOC7k85HCT8jXRM4F5zhnwYKLmVkGiYQHUchp2RgmKNxdcoJNOfYnvZL0dWZvnrl/M0J3vTmw/xXmqqsbjBI1v5EXvjCYiwyLgZ6ZIcOyDgYFbZj6/YAVpQEaFzmwZSlOFL8AC7YTh2uFYu1hxKOOOVJrAd4/kBXUnd1SqOyj4lPlQ7RMTym9viFkfUisIZOFWQMCIPJrrc4qba5qjHmErpamFGGC7QuFa0bBs6Newq3Qj4mUWBSYC7",
		Fingerprint:   "a6:92:c9:0a:c4:ca:7d:2c:fb:98:42:63:79:37:a8:24",
		HasPrivateKey: false,
	}

	err = client.CreateSshKey(&sshkey)
	if err != nil {
		t.Errorf("getting error calling CreateSshKey: %s", err)
	}

	assert.Equal(t, "Dd8v0nXJ1924rayY", sshkey.ID)
	assert.Equal(t, "my-ssh-key-without-private-key", sshkey.Name)
	assert.Equal(t, "a6:92:c9:0a:c4:ca:7d:2c:fb:98:42:63:79:37:a8:24", sshkey.Fingerprint)
	assert.Equal(t, "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDEc2hHw894iG/BULKSwdJceqLG8qGoifKfCVNTkFJgILn7dnJg7ioWtzpztaWCiwJe7R2zVF68ANi0+Ie8UM70/MFk5WKB1mOC7k85HCT8jXRM4F5zhnwYKLmVkGiYQHUchp2RgmKNxdcoJNOfYnvZL0dWZvnrl/M0J3vTmw/xXmqqsbjBI1v5EXvjCYiwyLgZ6ZIcOyDgYFbZj6/YAVpQEaFzmwZSlOFL8AC7YTh2uFYu1hxKOOOVJrAd4/kBXUnd1SqOyj4lPlQ7RMTym9viFkfUisIZOFWQMCIPJrrc4qba5qjHmErpamFGGC7QuFa0bBs6Newq3Qj4mUWBSYC7", sshkey.PublicKey)
}

func TestDeleteSshKey(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}))

	client, err := NewAPI("token123", WithEndpoint(server.URL))
	if err != nil {
		t.Errorf("getting error calling NewAPI:%s", err)
	}

	err = client.DeleteSshKey("jDGPRJXLpGXeV5M1")
	if err != nil {
		t.Errorf("getting error calling DeleteSshKey: %s", err)
	}
}
