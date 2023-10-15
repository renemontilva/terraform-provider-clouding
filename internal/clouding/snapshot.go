package clouding

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	SNAPSHOT_PATH = "snapshots"
)

type Snapshot struct {
	ID              string       `json:"id"`
	Name            string       `json:"name,omitempty"`
	SizeGb          int64        `json:"sizeGb,omitempty"`
	Description     string       `json:"description,omitempty"`
	CreatedAt       string       `json:"createdAt,omitempty"`
	SourceServeName string       `json:"sourceServerName,omitempty"`
	ShutDownServer  bool         `json:"shutDownServer,omitempty"`
	Image           Image        `json:"image,omitempty"`
	Cost            SnapshotCost `json:"cost,omitempty"`
}

type SnapshotCost struct {
	PricePerHour        float64 `json:"pricePerHour"`
	PricePerMonthApprox float64 `json:"pricePerMonthApprox"`
}

func (a *API) GetSnapshotID(id string) (Snapshot, error) {
	var snapshot Snapshot

	response, err := a.sendRequest(http.MethodGet, fmt.Sprintf("%s/%s", SNAPSHOT_PATH, id), nil)
	if err != nil {
		return snapshot, fmt.Errorf("getting error from sendRequest: %s", err)
	}
	if response.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			return snapshot, fmt.Errorf("error decoding error response: %s", err)
		}
		return snapshot, fmt.Errorf("error getting snapshot, status code: %d, title: %s", errorResponse.Status, errorResponse.Title)
	}

	err = json.NewDecoder(response.Body).Decode(&snapshot)
	if err != nil {
		return snapshot, fmt.Errorf("error decoding snapshot: %s", err)
	}

	return snapshot, nil
}
