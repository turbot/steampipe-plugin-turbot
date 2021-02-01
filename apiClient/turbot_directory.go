package apiClient

var turbotDirectoryProperties = []interface{}{
	map[string]string{"parent": "turbot.parentId"},
	"title",
	"description",
	"status",
	"directoryType",
	"profileIdTemplate",
	"server",
}

func (client *Client) CreateTurbotDirectory(input map[string]interface{}) (*TurbotDirectory, error) {
	query := createTurbotDirectoryMutation(turbotDirectoryProperties)
	responseData := &TurbotDirectoryResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "turbot directory")
	}
	return &responseData.Resource, nil
}

func (client *Client) ReadTurbotDirectory(id string) (*TurbotDirectory, error) {
	// create a map of the properties we want the graphql query to return
	query := readResourceQuery(id, turbotDirectoryProperties)
	responseData := &TurbotDirectoryResponse{}
	// execute api call
	if err := client.doRequest(query, nil, &responseData); err != nil {
		return nil, client.handleReadError(err, id, "turbot directory")
	}
	return &responseData.Resource, nil
}

func (client *Client) UpdateTurbotDirectory(input map[string]interface{}) (*TurbotDirectory, error) {
	query := updateTurbotDirectoryMutation(turbotDirectoryProperties)
	responseData := &TurbotDirectoryResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, &responseData); err != nil {
		return nil, client.handleUpdateError(err, input, "turbot directory")
	}
	return &responseData.Resource, nil
}
