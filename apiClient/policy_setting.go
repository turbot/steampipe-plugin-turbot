package apiClient

import (
	"fmt"
)

func (client *Client) CreatePolicySetting(input map[string]interface{}) (*PolicySetting, error) {
	query := createPolicySettingMutation()
	responseData := &PolicySettingResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "policy setting")
	}
	return &responseData.PolicySetting, nil
}

func (client *Client) ReadPolicySetting(id string) (*PolicySetting, error) {
	query := readPolicySettingQuery(id)
	responseData := &PolicySettingResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, id, "policy setting")
	}
	return &responseData.PolicySetting, nil
}

func (client *Client) UpdatePolicySetting(input map[string]interface{}) (*PolicySetting, error) {
	query := updatePolicySettingMutation()
	responseData := &PolicySettingResponse{}

	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleUpdateError(err, input, "policy setting")
	}
	return &responseData.PolicySetting, nil
}

func (client *Client) DeletePolicySetting(id string) error {
	query := deletePolicySettingMutation()
	responseData := &PolicySettingResponse{}
	variables := map[string]interface{}{
		"input": map[string]string{
			"id": id,
		},
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return fmt.Errorf("error deleting policy: %s", err.Error())
	}
	return nil
}

func (client *Client) FindPolicySetting(policyTypeUri, resourceAka string) (PolicySetting, error) {
	responseData := &FindPolicySettingResponse{}

	query := findPolicySettingQuery(policyTypeUri, resourceAka)

	// execute api call
	if err := client.doRequest(query, nil, &responseData); err != nil {
		return PolicySetting{}, client.handleReadError(err, policyTypeUri, "policy setting")
	}

	for _, setting := range responseData.PolicySettings.Items {
		if setting.Default {
			return setting, nil
		}
	}
	return PolicySetting{}, nil
}
