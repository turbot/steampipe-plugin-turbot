package apiClient

var profileProperties = []interface{}{
	map[string]string{"parent": "turbot.parentId"},
	"title",
	"status",
	"displayName",
	"email",
	"givenName",
	"familyName",
	"directoryPoolId",
	"profileId",
	"middleName",
	"picture",
	"externalId",
	"lastLoginTimestamp",
}

func (client *Client) CreateProfile(input map[string]interface{}) (*Profile, error) {
	query := createResourceMutation(profileProperties)
	responseData := &ProfileResponse{}
	// set type in input data
	input["type"] = "tmod:@turbot/turbot-iam#/resource/types/profile"
	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "profile")
	}
	return &responseData.Resource, nil
}

func (client *Client) ReadProfile(id string) (*Profile, error) {
	// create a map of the properties we want the graphql query to return

	query := readResourceQuery(id, profileProperties)
	responseData := &ProfileResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, id, "profile")
	}
	return &responseData.Resource, nil
}

func (client *Client) UpdateProfile(input map[string]interface{}) (*Profile, error) {
	query := updateResourceMutation(profileProperties)
	responseData := &ProfileResponse{}
	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleUpdateError(err, input, "profile")
	}
	return &responseData.Resource, nil
}
