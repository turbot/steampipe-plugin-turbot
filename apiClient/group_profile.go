package apiClient

import (
	"fmt"
)

var groupProfileProperties = []interface{}{
	map[string]string{"parent": "turbot.parentId"},
	"directory",
	"title",
	"status",
	"groupProfileId",
}

func (client *Client) CreateGroupProfile(input map[string]interface{}) (*GroupProfile, error) {
	query := createGroupProfileMutation(groupProfileProperties)
	responseData := &GroupProfileResponse{}
	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "group profile")
	}
	return &responseData.Resource, nil
}

func (client *Client) ReadGroupProfile(id string) (*GroupProfile, error) {
	// create a map of the properties we want the graphql query to return

	query := readResourceQuery(id, groupProfileProperties)
	responseData := &GroupProfileResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, id, "group profile")
	}
	return &responseData.Resource, nil
}

func (client *Client) UpdateGroupProfile(input map[string]interface{}) (*GroupProfile, error) {
	query := updateGroupProfileMutation(groupProfileProperties)
	responseData := &GroupProfileResponse{}
	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleUpdateError(err, input, "group profile")
	}
	return &responseData.Resource, nil
}

func (client *Client) DeleteGroupProfile(aka string) error {
	query := deleteGroupProfileMutation()
	// we do not care about the response
	var responseData interface{}

	variables := map[string]interface{}{
		"input": map[string]string{
			"id": aka,
		},
	}

	// execute api call
	if err := client.doRequest(query, variables, &responseData); err != nil {
		return fmt.Errorf("error deleting resource: %s", err.Error())
	}
	return nil
}
