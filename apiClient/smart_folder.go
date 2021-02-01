package apiClient

func (client *Client) CreateSmartFolder(input map[string]interface{}) (*SmartFolder, error) {
	query := createSmartFolderMutation()
	responseData := &SmartFolderResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "smart folder")
	}
	return &responseData.SmartFolder, nil
}

func (client *Client) ReadSmartFolder(id string) (*SmartFolder, error) {
	query := readSmartFolderQuery(id)
	responseData := &SmartFolderResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, id, "smart folder")
	}
	return &responseData.SmartFolder, nil
}

func (client *Client) UpdateSmartFolder(input map[string]interface{}) (*SmartFolder, error) {
	query := updateSmartFolderMutation()
	responseData := &SmartFolderResponse{}
	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleUpdateError(err, input, "smart folder")
	}
	return &responseData.SmartFolder, nil
}
