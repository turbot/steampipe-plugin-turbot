package apiClient

var folderProperties = []interface{}{
	//explicit mapping
	map[string]string{
		"parent": "turbot.parentId",
	},
	// implicit mapping
	"title",
	"description",
}

func (client *Client) CreateFolder(input map[string]interface{}) (*Folder, error) {
	query := createResourceMutation(folderProperties)
	responseData := &FolderResponse{}
	// set type in input data
	input["type"] = "tmod:@turbot/turbot#/resource/types/folder"
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "folder")
	}
	return &responseData.Resource, nil
}

func (client *Client) ReadFolder(id string) (*Folder, error) {
	// create a map of the properties we want the graphql query to return

	query := readResourceQuery(id, folderProperties)
	responseData := &FolderResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, id, "folder")
	}
	return &responseData.Resource, nil
}

func (client *Client) UpdateFolder(input map[string]interface{}) (*Folder, error) {
	query := updateResourceMutation(folderProperties)
	responseData := &FolderResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleUpdateError(err, input, "folder")
	}
	return &responseData.Resource, nil
}
