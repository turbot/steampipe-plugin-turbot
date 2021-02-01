package apiClient

var localDirectoryProperties = []interface{}{
	map[string]string{"parent": "turbot.parentId"},
	"title",
	"description",
	"status",
	"directoryType",
	"profileIdTemplate",
}

func (client *Client) ReadLocalDirectory(id string) (*LocalDirectory, error) {
	// create a map of the properties we want the graphql query to return
	query := readResourceQuery(id, localDirectoryProperties)
	responseData := &LocalDirectoryResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, id, "local directory")
	}
	return &responseData.Resource, nil
}

func (client *Client) CreateLocalDirectory(input map[string]interface{}) (*LocalDirectory, error) {
	query := createLocalDirectoryMutation(localDirectoryProperties)
	responseData := &LocalDirectoryResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "local directory")
	}
	return &responseData.Resource, nil
}

func (client *Client) UpdateLocalDirectory(input map[string]interface{}) (*LocalDirectory, error) {
	query := updateLocalDirectoryMutation(localDirectoryProperties)
	responseData := &LocalDirectoryResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleUpdateError(err, input, "local directory")
	}
	return &responseData.Resource, nil
}
