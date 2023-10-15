package clouding

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	FIREWALL_PATH = "firewalls"
)

type Firewall struct {
	ID             string               `json:"id,omitempty"`
	Name           string               `json:"name"`
	Description    string               `json:"description"`
	NewName        string               `json:"newName,omitempty"`
	NewDescription string               `json:"newDescription,omitempty"`
	Rules          []FirewallRule       `json:"rules,omitempty"`
	Attachments    []FirewallAttachment `json:"attachments,omitempty"`
}

type FirewallRuleID struct {
	FirewallID   string       `json:"firewallId"`
	FirewallRule FirewallRule `json:"firewallRule"`
}

type FirewallRule struct {
	ID           string `json:"id"`
	SourceIP     string `json:"sourceIp"`
	Protocol     string `json:"protocol"`
	Description  string `json:"description"`
	PortRangeMin int64  `json:"portRangeMin"`
	PortRangeMax int64  `json:"portRangeMax"`
	Enabled      bool   `json:"enabled"`
}

type FirewallAttachment struct {
	ServerID   string `json:"serverId"`
	ServerName string `json:"serverName"`
}

// GetFirewallID returns the firewall ID.
func (a *API) GetFirewallID(id string) (Firewall, error) {
	var firewall Firewall

	response, err := a.sendRequest(http.MethodGet, fmt.Sprintf("%s/%s", FIREWALL_PATH, id), nil)
	if err != nil {
		return firewall, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			return firewall, fmt.Errorf("error decoding error response: %s", err)
		}
		return firewall, fmt.Errorf("error getting firewall: %s", errorResponse.Title)
	}

	err = json.NewDecoder(response.Body).Decode(&firewall)
	if err != nil {
		return firewall, fmt.Errorf("error decoding firewall: %s", err)
	}

	return firewall, nil
}

func (a *API) CreateFirewall(firewall *Firewall) error {
	firewallJSON, err := json.Marshal(firewall)
	if err != nil {
		return fmt.Errorf("error marshaling firewall: %s", err)
	}

	response, err := a.sendRequest(http.MethodPost, FIREWALL_PATH, firewallJSON)
	if err != nil {
		return fmt.Errorf("getting error from sendRequest: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		var errorResponse ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			return fmt.Errorf("error decoding error response: %s", err)
		}
		return fmt.Errorf("error creating firewall, status code: %d, title: %s", errorResponse.Status, errorResponse.Title)

	}
	err = json.NewDecoder(response.Body).Decode(&firewall)
	if err != nil {
		return fmt.Errorf("error decoding firewall: %s", err)
	}

	return nil
}

func (a *API) UpdateFirewall(id string, firewall Firewall) error {
	firewallJSON, err := json.Marshal(firewall)
	if err != nil {
		return fmt.Errorf("error marshaling firewall: %s", err)
	}
	response, err := a.sendRequest(http.MethodPatch, fmt.Sprintf("%s/%s", FIREWALL_PATH, id), firewallJSON)
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
		return fmt.Errorf("error updating firewall, status code: %d, title: %s", errorResponse.Status, errorResponse.Title)
	}

	return nil
}

func (a *API) DeleteFirewall(id string) error {
	response, err := a.sendRequest(http.MethodDelete, fmt.Sprintf("%s/%s", FIREWALL_PATH, id), nil)
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
		return fmt.Errorf("error deleting firewall, status code: %d, title: %s", errorResponse.Status, errorResponse.Title)
	}

	return nil
}

func (a *API) GetFirewallRule(id string) (FirewallRuleID, error) {
	var firewallRuleID FirewallRuleID
	response, err := a.sendRequest(http.MethodGet, fmt.Sprintf("%s/%s/%s", FIREWALL_PATH, "rules", id), nil)
	if err != nil {
		return firewallRuleID, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			return firewallRuleID, fmt.Errorf("error decoding error response: %s", err)
		}
		return firewallRuleID, fmt.Errorf("error getting firewall rule, status code: %d, title: %s", errorResponse.Status, errorResponse.Title)
	}
	err = json.NewDecoder(response.Body).Decode(&firewallRuleID)
	if err != nil {
		return firewallRuleID, fmt.Errorf("error decoding firewall rule: %s", err)
	}
	return firewallRuleID, nil
}

func (a *API) CreateFirewallRule(rule *FirewallRuleID) error {
	ruleJSON, err := json.Marshal(rule.FirewallRule)
	if err != nil {
		return fmt.Errorf("error marshaling firewall rule: %s", err)
	}

	response, err := a.sendRequest(http.MethodPost, fmt.Sprintf("%s/%s/%s", FIREWALL_PATH, rule.FirewallID, "rules"), ruleJSON)
	if err != nil {
		return fmt.Errorf("getting error from sendRequest: %s", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		var errorResponse ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			return fmt.Errorf("error decoding error response: %s", err)
		}
		return fmt.Errorf("error creating firewall rule, status code: %d, title: %s", errorResponse.Status, errorResponse.Title)
	}
	err = json.NewDecoder(response.Body).Decode(&rule.FirewallRule)
	if err != nil {
		return fmt.Errorf("error decoding firewall rule: %s", err)
	}

	return nil
}

func (a *API) DeleteFirewallRule(id string) error {
	response, err := a.sendRequest(http.MethodDelete, fmt.Sprintf("%s/%s/%s", FIREWALL_PATH, "rules", id), nil)
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
		return fmt.Errorf("error deleting firewall rule, status code: %d, title: %s", errorResponse.Status, errorResponse.Title)
	}
	return nil
}
