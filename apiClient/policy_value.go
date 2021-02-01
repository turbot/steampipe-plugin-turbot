package apiClient

func (client *Client) ReadPolicyValue(policyTypeUri, resourceAka string) (*PolicyValue, error) {
	query := readPolicyValueQuery(policyTypeUri, resourceAka)
	responseData := &PolicyValueResponse{}
	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, policyTypeUri, "policy setting")
	}

	return &responseData.PolicyValue, nil
}
