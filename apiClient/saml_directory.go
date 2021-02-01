package apiClient

// create a map of the properties we want the graphql query to return
var samlDirectoryProperties = []interface{}{
	map[string]string{"parent": "turbot.parentId"},
	"title",
	"description",
	"status",
	"directoryType",
	"profileIdTemplate",
	"entryPoint",
	"certificate",
	"issuer",
	"nameIdFormat",
	"signRequests",
	"signaturePrivateKey",
	"signatureAlgorithm",
	"poolId",
	"allowIdpInitiatedSso",
	"allowGroupSyncing",
	"profileGroupsAttribute",
	"groupFilter",
}

func (client *Client) ReadSamlDirectory(id string) (*SamlDirectory, error) {

	query := readResourceQuery(id, samlDirectoryProperties)
	responseData := &SamlDirectoryResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, id, "saml directory")
	}
	return &responseData.Resource, nil
}

func (client *Client) CreateSamlDirectory(input map[string]interface{}) (*SamlDirectory, error) {
	query := createSamlDirectoryMutation(samlDirectoryProperties)
	responseData := &SamlDirectoryResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "saml directory")
	}
	return &responseData.Resource, nil
}

func (client *Client) UpdateSamlDirectory(input map[string]interface{}) (*SamlDirectory, error) {
	query := updateSamlDirectoryMutation(samlDirectoryProperties)
	responseData := &SamlDirectoryResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleUpdateError(err, input, "saml directory")
	}
	return &responseData.Resource, nil
}
