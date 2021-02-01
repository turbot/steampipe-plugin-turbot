package apiClient

import (
	"fmt"
	"strings"
)

func (client *Client) InstallMod(input map[string]interface{}) (*InstallModData, error) {
	query := installModMutation()
	responseData := &InstallModResponse{}

	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, fmt.Errorf("error installing mod: %s", err.Error())
	}
	return &responseData.Mod, nil
}

func (client *Client) ReadMod(id string) (*Mod, error) {
	query := readModQuery(id)
	responseData := &ReadModResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, id, "mod")
	}

	// convert uri into org and mod
	org, mod := ParseModUri(responseData.Mod.Uri)
	responseData.Mod.Org = org
	responseData.Mod.Mod = mod
	return &responseData.Mod, nil
}

func ParseModUri(uri string) (org, mod string) {
	if uri == "" {
		org = ""
		mod = ""
		return
	}
	// uri will be of form "tmod:@<org>/<mod>"
	segments := strings.Split(strings.TrimPrefix(uri, "tmod:@"), "/")
	org = segments[0]
	mod = segments[1]
	return
}

func (client *Client) UninstallMod(modId string) error {
	query := uninstallModMutation()
	responseData := &UninstallModResponse{}

	variables := map[string]interface{}{
		"input": map[string]string{
			"id": modId,
		},
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return fmt.Errorf("error uninstalling mod: %s", err.Error())
	}
	if !responseData.UninstallMod.Success {
		return fmt.Errorf(" uninstallMod mutation ran with no errors but failed to uninstall the mod")
	}

	return nil
}

func (client *Client) GetModVersions(org, mod string) ([]ModRegistryVersion, error) {
	query := modVersionsQuery(org, mod)
	responseData := &ModVersionResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error fetching mod versions mod: %s", err.Error())
	}

	return responseData.Versions.Items, nil
}
