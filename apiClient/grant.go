package apiClient

import (
	"fmt"
)

func (client *Client) CreateGrant(input map[string]interface{}) (*TurbotGrantMetadata, error) {
	query := createGrantMutation()
	responseData := &CreateGrantResponse{}

	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "grant")
	}
	return &responseData.Grants.Turbot, nil
}

func (client *Client) ReadGrant(id string) (*Grant, error) {
	query := readGrantQuery(id)
	responseData := &ReadGrantResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, id, "grant")
	}
	return &responseData.Grant, nil
}

func (client *Client) DeleteGrant(id string) error {
	query := deleteGrantMutation()
	var responseData interface{}
	variables := map[string]interface{}{
		"input": map[string]string{
			"id": id,
		},
	}

	// execute api call
	if err := client.doRequest(query, variables, &responseData); err != nil {
		return fmt.Errorf("error deleting grant: %s", err.Error())
	}
	return nil
}

func (client *Client) GrantExists(id string) (bool, error) {
	grant, err := client.ReadGrant(id)
	if err != nil {
		return false, err
	}
	exists := grant.Turbot.Id != ""
	return exists, nil
}
