package clouding

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	IMAGE_PATH = "images"
)

type Image struct {
	ID                  string            `json:"id"`
	Name                string            `json:"name"`
	MinimumSizeGB       int64             `json:"minimumSizeGb"`
	AccessMethods       ImageAccessMethod `json:"accessMethods"`
	PricePerHour        float64           `json:"pricePerHour"`
	PricePerMonthApprox float64           `json:"pricePerMonthApprox"`
	BillingUnit         string            `json:"billingUnit"`
}

type ImageAccessMethod struct {
	SshKey   string `json:"sshKey"`
	Password string `json:"password"`
}

func (a *API) GetImageID(id string) (Image, error) {
	var image Image

	response, err := a.sendRequest(http.MethodGet, fmt.Sprintf("%s/%s", IMAGE_PATH, id), nil)
	if err != nil {
		return image, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		err = json.NewDecoder(response.Body).Decode(&errorResponse)
		if err != nil {
			return image, fmt.Errorf("error decoding error response: %s", err)
		}
		return image, fmt.Errorf("error getting image: %s", errorResponse.Detail)
	}

	err = json.NewDecoder(response.Body).Decode(&image)
	if err != nil {
		return image, fmt.Errorf("error decoding image: %s", err)
	}
	return image, nil
}
