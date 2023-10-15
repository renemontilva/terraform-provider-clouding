package clouding

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	SERVER_PATH = "servers"
)

type Server struct {
	ID                            string               `json:"id,omitempty"`
	Name                          string               `json:"name,omitempty"`
	NewServerName                 string               `json:"newServerName,omitempty"`
	Hostname                      string               `json:"hostname,omitempty"`
	VCores                        float64              `json:"vCores,omitempty"`
	RamGb                         int                  `json:"ramGb,omitempty"`
	FlavorID                      string               `json:"flavorId,omitempty"`
	Flavor                        string               `json:"flavor,omitempty"`
	FirewallID                    string               `json:"firewallId,omitempty"`
	AccessConfiguration           *AccessConfiguration `json:"accessConfiguration,omitempty"`
	RequestedAccessConfiguration  *AccessConfiguration `json:"requestedAccessConfiguration,omitempty"`
	VolumeSizeGb                  int64                `json:"volumeSizeGb,omitempty"`
	Volume                        *Volume              `json:"volume,omitempty"`
	EnablePrivateNetwork          bool                 `json:"enablePrivateNetwork,omitempty"`
	EnableStrictAntiDDoSFiltering bool                 `json:"enableStrictAntiDDoSFiltering,omitempty"`
	UserData                      string               `json:"userData,omitempty"`
	BackupPreference              *BackupPreference    `json:"backupPreferences,omitempty"`
	Image                         Image                `json:"image,omitempty"`
	Status                        string               `json:"status,omitempty"`
	PowerState                    string               `json:"powerState,omitempty"`
	Features                      []string             `json:"features,omitempty"`
	PendingFeatures               []string             `json:"pendingFeatures,omitempty"`
	PendingFirewalls              []string             `json:"pendingFirewalls,omitempty"`
	CreatedAt                     string               `json:"createdAt,omitempty"`
	DnsAddress                    string               `json:"dnsAddress,omitempty"`
	PublicIP                      string               `json:"publicIp,omitempty"`
	PrivateIP                     string               `json:"privateIp,omitempty"`
	SshKeyID                      string               `json:"sshKeyId,omitempty"`
	Firewalls                     []Firewall           `json:"firewalls,omitempty"`
	Snapshots                     []Snapshot           `json:"snapshots,omitempty"`
	Backups                       []Backup             `json:"backups,omitempty"`
	Cost                          ServerCost           `json:"cost,omitempty"`
	Action                        Action               `json:"action,omitempty"`
}

type AccessConfiguration struct {
	SshKey       string `json:"ssh_key,omitempty"`
	SshKeyID     string `json:"sshKeyId,omitempty"`
	Password     string `json:"password,omitempty"`
	HasPassword  bool   `json:"hasPassword,omitempty"`
	SavePassword bool   `json:"savePassword,omitempty"`
}

type Volume struct {
	ID             string `json:"id,omitempty"`
	Source         string `json:"source,omitempty"`
	SsdGb          int64  `json:"ssdGb,omitempty"`
	ShutDownSource bool   `json:"shutDownSource,omitempty"`
}

type BackupPreference struct {
	Slots     int64  `json:"slots,omitempty"`
	Frequency string `json:"frequency,omitempty"`
}

type ServerCost struct {
	PricePerHour        float64 `json:"pricePerHour,omitempty"`
	PricePerMonthApprox float64 `json:"pricePerMonthApprox,omitempty"`
}

func (a *API) GetServerID(server *Server) error {

	response, err := a.sendRequest(http.MethodGet, fmt.Sprintf("%s/%s", SERVER_PATH, server.ID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			return fmt.Errorf("error decoding error response: %s", err)
		}
		return fmt.Errorf("error getting server: %s", errorResponse.Detail)
	}

	err = json.NewDecoder(response.Body).Decode(&server)
	if err != nil {
		return fmt.Errorf("error decoding server: %s", err)
	}
	server.FlavorID = server.Flavor
	server.FirewallID = server.Firewalls[0].ID
	if server.Volume != nil {
		server.Volume.SsdGb = server.VolumeSizeGb
		// FIXME: This is a workaround to avoid volume source value inconsistency"
		server.Volume.ID = server.Image.ID
	}

	return nil
}

func (a *API) CreateServer(server *Server) error {
	serverJSON, err := json.Marshal(server)
	var serverResponse Server
	if err != nil {
		return fmt.Errorf("error marshaling server: %s", err)
	}

	response, err := a.sendRequest(http.MethodPost, SERVER_PATH, serverJSON)
	if err != nil {
		return fmt.Errorf("getting error from sendRequest: %s", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted {
		var errorResponse ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			return fmt.Errorf("error decoding error response: %s", err)
		}
		return fmt.Errorf("error creating server, status code: %d, title: %s, detail: %s", errorResponse.Status, errorResponse.Title, errorResponse.Detail)

	}

	err = json.NewDecoder(response.Body).Decode(&serverResponse)
	if err != nil {
		return fmt.Errorf("error decoding server: %s", err)
	}

	server.ID = serverResponse.ID
	server.AccessConfiguration.SshKeyID = serverResponse.RequestedAccessConfiguration.SshKeyID
	server.AccessConfiguration.SavePassword = serverResponse.RequestedAccessConfiguration.SavePassword
	server.Status = serverResponse.Status
	server.Action = serverResponse.Action

	return nil
}

func (a *API) DeleteServer(id string) (Action, error) {
	var action Action
	response, err := a.sendRequest(http.MethodDelete, fmt.Sprintf("%s/%s", SERVER_PATH, id), nil)
	if err != nil {
		return action, fmt.Errorf("getting error from sendRequest: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		var errorResponse ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			return action, fmt.Errorf("error decoding error response: %s", err)
		}
		return action, fmt.Errorf("error deleting server, status code: %d, title: %s", errorResponse.Status, errorResponse.Title)
	}

	err = json.NewDecoder(response.Body).Decode(&action)
	if err != nil {
		return action, fmt.Errorf("error decoding action: %s", err)
	}

	return action, nil
}

func (a *API) UpdateServerName(id, name string) error {
	server := Server{
		NewServerName: name,
	}
	serverJSON, err := json.Marshal(server)
	if err != nil {
		return fmt.Errorf("error marshaling server: %s", err)
	}

	response, err := a.sendRequest(http.MethodPatch, fmt.Sprintf("%s/%s/rename", SERVER_PATH, id), serverJSON)
	if err != nil {
		return fmt.Errorf("getting error from sendRequest: %s", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusNoContent {
		var errorResponse ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			return fmt.Errorf("error decoding error response: %s", err)
		}
		return fmt.Errorf("error updating server name, status code: %d, title: %s", errorResponse.Status, errorResponse.Title)
	}

	return nil
}
