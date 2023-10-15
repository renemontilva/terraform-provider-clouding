package clouding

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	SSHKEY_PATH = "keypairs"
)

type SshKey struct {
	ID            string `json:"id"`
	Name          string `json:"name,omitempty"`
	Fingerprint   string `json:"fingerprint,omitempty"`
	PublicKey     string `json:"publicKey,omitempty"`
	PrivateKey    string `json:"privateKey,omitempty"`
	HasPrivateKey bool   `json:"hasPrivateKey,omitempty"`
}

func (a *API) GetSshKeyID(id string) (SshKey, error) {
	var sshKey SshKey

	response, err := a.sendRequest(http.MethodGet, fmt.Sprintf("%s/%s", SSHKEY_PATH, id), nil)
	if err != nil {
		return sshKey, fmt.Errorf("getting error from sendRequest: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)
		return sshKey, fmt.Errorf("error getting sshkey, status code: %d, title: %s", errorResponse.Status, errorResponse.Title)
	}

	err = json.NewDecoder(response.Body).Decode(&sshKey)
	if err != nil {
		return sshKey, fmt.Errorf("error decoding sshkey: %s", err)
	}

	return sshKey, nil
}

func (a *API) CreateSshKey(sshKey *SshKey) error {
	sshKeyJSON, err := json.Marshal(sshKey)
	if err != nil {
		return fmt.Errorf("error marshaling sshkey: %s", err)
	}

	response, err := a.sendRequest(http.MethodPost, SSHKEY_PATH, sshKeyJSON)
	if err != nil {
		return fmt.Errorf("getting error from sendRequest: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)
		return fmt.Errorf("error creating sshkey, status code: %d, title: %s", errorResponse.Status, errorResponse.Title)
	}

	json.NewDecoder(response.Body).Decode(sshKey)
	return nil
}

func (a *API) DeleteSshKey(id string) error {
	response, err := a.sendRequest(http.MethodDelete, fmt.Sprintf("%s/%s", SSHKEY_PATH, id), nil)
	if err != nil {
		return fmt.Errorf("getting error from sendRequest: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)
		return fmt.Errorf("error deleting sshkey, status code: %d, title: %s", errorResponse.Status, errorResponse.Title)
	}

	return nil
}
