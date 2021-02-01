package apiClient

var googleDirectoryProperties = []interface{}{
	// implicit mappings
	"title", "poolId", "profileIdTemplate", "groupIdTemplate", "loginNameTemplate", "clientSecret", "hostedDomain", "description", "clientId"}

func (client *Client) ReadGoogleDirectory(id string) (*GoogleDirectory, error) {
	/*
		GoogleDirectory read response has clientSecret attribute,
		which is fetched from getSecret(path:"clientSecret") and
		not from get() resolver.
		That's why we used separate query and not readResourceQuery()
	*/
	query := readGoogleDirectoryQuery(id)
	responseData := &ReadGoogleDirectoryResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, id, "google")
	}
	return &responseData.Directory, nil
}

func (client *Client) CreateGoogleDirectory(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := createGoogleDirectoryMutation(googleDirectoryProperties)
	responseData := &CreateResourceResponse{}
	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "google")
	}
	return &responseData.Resource.Turbot, nil
}

func (client *Client) UpdateGoogleDirectory(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := updateGoogleDirectoryMutation(googleDirectoryProperties)
	responseData := &UpdateResourceResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleUpdateError(err, input, "google")
	}
	return &responseData.Resource.Turbot, nil
}
