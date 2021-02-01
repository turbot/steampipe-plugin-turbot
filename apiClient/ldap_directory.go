package apiClient

import (
	"fmt"

	"github.com/turbot/steampipe-plugin-turbot/helpers"
)

var ldapDirectoryProperties = []interface{}{
	map[string]string{"parent": "turbot.parentId"},
	"title",
	"description",
	"profileIdTemplate",
	"status",
	"groupProfileIdTemplate",
	"url",
	"distinguishedName",
	"password",
	"base",
	"userObjectFilter",
	"directoryType",
	"disabledUserFilter",
	"userMatchFilter",
	"userSearchFilter",
	"userSearchAttributes",
	"groupObjectFilter",
	"groupSearchFilter",
	"groupSyncFilter",
	"userCanonicalNameAttribute",
	"userEmailAttribute",
	"userDisplayNameAttribute",
	"userGivenNameAttribute",
	"userFamilyNameAttribute",
	"tlsEnabled",
	"tlsServerCertificate",
	"groupMemberOfAttribute",
	"groupMembershipAttribute",
	"connectivityTestFilter",
	"rejectUnauthorized",
	"disabledGroupFilter",
}

// exclude password from read call, secret id
func getLdapDirectoryReadProperties() []interface{} {
	excludedProperties := []string{"password"}
	return helpers.RemoveProperties(ldapDirectoryProperties, excludedProperties)
}

func (client *Client) CreateLdapDirectory(input map[string]interface{}) (*LdapDirectory, error) {
	query := createLdapDirectoryMutation(ldapDirectoryProperties)
	responseData := &LdapDirectoryResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "ldap directory")
	}
	return &responseData.Resource, nil
}

func (client *Client) ReadLdapDirectory(id string) (*LdapDirectory, error) {
	// create a map of the properties we want the graphql query to return
	query := readResourceQuery(id, getLdapDirectoryReadProperties())
	responseData := &LdapDirectoryResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, id, "ldap directory")
	}
	return &responseData.Resource, nil
}

func (client *Client) UpdateLdapDirectory(input map[string]interface{}) (*LdapDirectory, error) {
	query := updateLdapDirectoryMutation(ldapDirectoryProperties)
	responseData := &LdapDirectoryResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleUpdateError(err, input, "ldap directory")
	}
	return &responseData.Resource, nil
}

func (client *Client) DeleteLdapDirectory(aka string) error {
	query := deleteLdapDirectory()
	// we do not care about the response
	var responseData interface{}

	variables := map[string]interface{}{
		"input": map[string]string{
			"id": aka,
		},
	}

	// execute api call
	if err := client.doRequest(query, variables, &responseData); err != nil {
		return fmt.Errorf("error deleting ldap directory: %s", err.Error())
	}
	return nil
}
