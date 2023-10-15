package clouding

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Action struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	Type         string `json:"type"`
	StartedAt    string `json:"startedAt"`
	CompletedAt  string `json:"completedAt"`
	ResourceID   string `json:"resourceId"`
	ResourceType string `json:"resourceType"`
}

const (
	// Action
	ACTION_PATH = "actions"
	PENDING     = "pending"
	PROCESSING  = "inProgress"
	COMPLETED   = "completed"
	ERROR       = "errored"
)

func (a *API) GetAction(id string) (Action, error) {
	var action Action

	response, err := a.sendRequest(http.MethodGet, fmt.Sprintf("%s/%s", ACTION_PATH, id), nil)
	if err != nil {
		return action, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)
		return action, fmt.Errorf("error getting action: %s", errorResponse.Detail)
	}

	json.NewDecoder(response.Body).Decode(&action)
	return action, nil
}

func (a *API) WaitForAction(ctx context.Context, action *Action, waitTime time.Duration) error {
	actionResponse, err := a.GetAction(action.ID)
	if err != nil {
		return err
	}

	for actionResponse.Status == PENDING || actionResponse.Status == PROCESSING {
		actionResponse, err = a.GetAction(action.ID)
		if err != nil {
			return err
		}
		time.Sleep(waitTime)
	}

	if actionResponse.Status == ERROR {
		return fmt.Errorf("action id: %s failed", action.ID)
	}

	// Update the action with the latest status
	action.Status = actionResponse.Status
	action.CompletedAt = actionResponse.CompletedAt

	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}
