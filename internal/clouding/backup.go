package clouding

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	BACKUP_PATH = "backups"
)

type Backup struct {
	ID           string `json:"id"`
	CreatedAt    string `json:"createdAt"`
	ServerID     string `json:"serverId"`
	ServerName   string `json:"serverName"`
	VolumeSizeGb int    `json:"volumeSizeGb"`
	Image        Image  `json:"image"`
	Status       string `json:"status"`
}

func (a *API) GetBackupID(id string) (Backup, error) {
	var backup Backup

	response, err := a.sendRequest(http.MethodGet, fmt.Sprintf("%s/%s", BACKUP_PATH, id), nil)
	if err != nil {
		return backup, fmt.Errorf("getting error from sendRequest: %s", err)
	}
	if response.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)
		return backup, fmt.Errorf("error getting backup, status code: %d, title: %s", errorResponse.Status, errorResponse.Title)
	}

	err = json.NewDecoder(response.Body).Decode(&backup)
	if err != nil {
		return backup, fmt.Errorf("error decoding backup: %s", err)
	}

	return backup, nil
}
