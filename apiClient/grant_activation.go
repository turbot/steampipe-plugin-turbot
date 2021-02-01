package apiClient

import (
	"fmt"
)

func (client *Client) CreateGrantActivation(input map[string]interface{}) (*TurbotActiveGrantMetadata, error) {
	query := activateGrantMutation()
	responseData := &ActivateGrantResponse{}

	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "grant activation")
	}
	return &responseData.GrantActivate.Turbot, nil
}

func (client *Client) ReadGrantActivation(id string) (*ActiveGrant, error) {
	query := readActiveGrantQuery(id)
	responseData := &ReadActiveGrantResponse{}
	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, id, "grant activation")
	}
	return &responseData.ActiveGrant, nil
}

func (client *Client) DeleteGrantActivation(id string) error {
	query := deactivateGrantMutation()
	var responseData interface{}

	variables := map[string]interface{}{
		"input": map[string]string{
			"activation": id,
		},
	}

	// execute api call
	if err := client.doRequest(query, variables, &responseData); err != nil {
		return fmt.Errorf("error deleting grant activation: %s", err.Error())
	}
	return nil
}

func (client *Client) GrantActivationExists(id string) (bool, error) {
	grantActivate, err := client.ReadGrantActivation(id)
	if err != nil {
		return false, err
	}
	exists := grantActivate.Turbot.Id != ""
	return exists, nil
}
