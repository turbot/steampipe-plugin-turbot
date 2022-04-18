package apiClient

import (
	"bytes"
	"fmt"
	"github.com/blang/semver"
	"strings"
)

// NOTE: we do not use a fragment for resource metadata as we just request the full turbot property
// using turbot: get(path:"turbot")
// This is because we saw errors returning null for the turbot property for a non existent resource
// TODO fix this to use a fragment

func turbotGrantMetadataFragment(prefix string) string {
	return applyPrefix(prefix,
		`turbot {
	id
	profileId  
	resourceId 
}`)
}

func turbotActiveGrantMetadataFragment(prefix string) string {
	return applyPrefix(prefix, `turbot {
	id
	grantId
	resourceId
}`)
}

// add the given prefix to each line of the multi-line inputString
func applyPrefix(prefix, inputString string) string {
	return strings.Replace(inputString, "\n", "\n"+prefix, -1)
}

// validation
func validationQuery() (string, ValidationResponse) {
	query := `{
	schema: __schema {
		queryType {
			name
		}
	}
}`

	return query, ValidationResponse{}
}

// policySetting
func createPolicySettingMutation() string {
	return `mutation CreatePolicySetting($input: CreatePolicySettingInput!) {
	policySetting: createPolicySetting(input: $input ) {
		type {
			uri
		}
		value: secretValue
		valueSource: secretValueSource
		template
		precedence
		templateInput
		input
		note
		validFromTimestamp
		validToTimestamp
		turbot {
		  id
		}
	}
}`

}

func readPolicySettingQuery(policySettingId string) string {
	return fmt.Sprintf(`{
policySetting(id:"%s") {
	type {
		uri
	}
	value: secretValue
	valueSource: secretValueSource
	template
	default
	precedence
	templateInput
	input
	note
	validFromTimestamp
	validToTimestamp
	turbot {
		id
		resourceId
	}
}
}`, policySettingId)
}

func updatePolicySettingMutation() string {
	return `mutation UpdatePolicySetting($input: UpdatePolicySettingInput!) {
	policySetting: updatePolicySetting(input: $input ) {
		type {
			uri
		}
		value: secretValue
		valueSource: secretValueSource
		template
		precedence
		templateInput
		input
		note
		validFromTimestamp
		validToTimestamp
		turbot {
			id
		}
	}
}`

}

func deletePolicySettingMutation() string {
	return `mutation DeletePolicySetting($input: DeletePolicySettingInput!) {
	policySetting: deletePolicySetting(input: $input ) {
		value: secretValue
		valueSource: secretValueSource
		template
		precedence
		templateInput
		input
		note
		validFromTimestamp
		validToTimestamp
		turbot {
			id
		}
	}
}`
}

func findPolicySettingQuery(policyTypeUri, resourceAka string) string {
	return fmt.Sprintf(`{
  policySettings: policySettingList(filter: "policyType:%s resource:%s") {
    items {
      	value: secretValue
		valueSource: secretValueSource
		template
		precedence
		templateInput
		input
		note
		validFromTimestamp
		validToTimestamp
		turbot {
			id
		}
    }
  }
}
`, policyTypeUri, resourceAka)
}

// policy value
func readPolicyValueQuery(policyTypeUri string, resourceId string) string {
	return fmt.Sprintf(`{
	policyValue(uri:"%s", resourceId:"%s"){
		value: secretValue
		secretValue
		precedence
		state
		reason
		details
		setting {
			valueSource
			turbot {
				id
			}
		}
		turbot {
			id
		}
	}
}
`, policyTypeUri, resourceId)
}

// smart folder
// filter and description are removed for a workaround, will be removed after a Core change.
func createSmartFolderMutation() string {
	return `mutation CreateSmartFolder($input: CreateSmartFolderInput!) {
		smartFolder: createSmartFolder(input: $input) {
			turbot {
				id
				parentId
				akas
				title
			}
		}
	}`
}

func readSmartFolderQuery(id string) string {
	return fmt.Sprintf(`{
	smartFolder: resource(id:"%s") {
		title: get(path:"turbot.title")
		description: get(path:"description")
		filters: get(path:"filters")
		parent:	get(path:"turbot.parentId")
		turbot: get(path:"turbot")
   		attachedResources{
			items{
				turbot: get(path:"turbot")
			}
		}
	}
}`, id)
}

func updateSmartFolderMutation() string {
	return `mutation UpdateSmartFolder($input: UpdateSmartFolderInput!) {
		smartFolder: updateSmartFolder(input: $input) {
			turbot {
				id
				parentId
				akas
			}
		}
	}`
}

func createSmartFolderAttachmentMutation() string {
	return `mutation AttachSmartFolder($input: AttachSmartFolderInput!) {
		attachSmartFolders(input: $input) {
			turbot {
				id
			}
		}
	}`
}

func detachSmartFolderAttachment() string {
	return `mutation DetachSmartFolder($input: DetachSmartFolderInput!) {
		detachSmartFolder: detachSmartFolders(input: $input) {
    		turbot {
				id
			}
  		}
	}`
}

// mod
func installModMutation() string {
	return `mutation InstallMod($input: InstallModInput!) {
	mod: installMod(input: $input) {
		turbot {
			id
			parentId
			akas
		}
		build
	}
}`
}

func readModQuery(modId string) string {
	return fmt.Sprintf(`{
	mod: resource(id:"%s") {
		uri: get(path: "turbot.akas.0")
		parent: get(path: "turbot.parentId")
		version: get(path: "version")
	}
}`, modId)
}

func uninstallModMutation() string {
	return `mutation UninstallMod($input: UninstallModInput!) {
	uninstallMod(input: $input) {
		success
	}
}`
}

func modVersionsQuery(org, mod string) string {
	return fmt.Sprintf(`{
	versions: modVersionList(orgName: "%s", modName: "%s") {
		items {
			status
			version
		}
	}
}`, org, mod)
}

// resource
func createResourceMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation CreateResource($input: CreateResourceInput!) {
	resource: createResource(input: $input) {
%s
		turbot: get(path:"turbot")
	}
}`, buildResourceProperties(properties))
}

func updateResourceMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation UpdateResource($input: UpdateResourceInput!) {
 	resource: updateResource(input: $input) {
%s
		turbot: get(path:"turbot")
	}
}`, buildResourceProperties(properties))
}

func deleteResourceMutation() string {
	return `mutation DeleteResource($input: DeleteResourceInput!) {
 	resource: deleteResource(input: $input) {
		turbot: get(path:"turbot")
	}
}`
}

// support properties array of Interface
func readResourceQuery(aka string, properties []interface{}) string {
	return fmt.Sprintf(`{
	resource(id:"%s") {
		type {
			uri
		}
%s
		turbot: get(path:"turbot")
  	}
}`, aka, buildResourceProperties(properties))
}

func getResourceTypeIdQuery(aka string) string {
	return fmt.Sprintf(`{
	resource(id:"%s") {
		turbot {
			resourceTypeId
		}
  	}
}`, aka)
}

func readResourceListQuery(filter string, properties map[string]string) string {
	var propertiesString bytes.Buffer
	for alias, propertyPath := range properties {
		propertiesString.WriteString(fmt.Sprintf("\t\t\t%s: get(path: \"%s\")\n", alias, propertyPath))
	}
	return fmt.Sprintf(`{
	resources(filter:"%s") {
		items{
			data
			type { uri }
%s
			turbot: get(path:"turbot")
		}
	}
}`, filter, propertiesString.String())
}

func readFullResourceQuery(aka string) string {
	return fmt.Sprintf(`{
  resource(id:"%s") {
	type {
		uri
	}
    data
    turbot: get(path:"turbot")
  }
}`, aka)
}

// google directory read query
func readGoogleDirectoryQuery(aka string) string {
	return fmt.Sprintf(`{
	directory: resource(id:"%s") {
		title:             	get(path:"title")
		parent:            	get(path:"turbot.parentId")
		description:       	get(path:"description")
		status:            	get(path:"status")
		directoryType:     	get(path:"directoryType")
		profileIdTemplate: 	get(path:"profileIdTemplate")
		clientID:          	get(path:"clientID")
		poolId:            	get(path:"poolId")
		groupIdTemplate:   	get(path:"groupIdTemplate")
		loginNameTemplate: 	get(path:"loginNameTemplate")
		hostedDomain:       get(path:"hostedDomain")
		turbot: 			get(path:"turbot")
	}
}`, aka)
}

// grant
func readGrantQuery(aka string) string {
	return fmt.Sprintf(`{
	grant: grant(id:"%s") {
		permissionTypeId
		permissionLevelId
		%s
	}
  }`, aka, turbotGrantMetadataFragment("\t\t"))
}

func createGrantMutation() string {
	return fmt.Sprintf(`mutation CreateGrant($input: CreateGrantInput!) {
	grants: createGrant(input: $input) {
%s
	}
}`, turbotGrantMetadataFragment("\t\t\t"))
}

func deleteGrantMutation() string {
	return fmt.Sprintf(`mutation DeleteGrant($input: DeleteGrantInput!) {
 	grant: deleteGrant(input: $input) {
%s
	}
}`, turbotGrantMetadataFragment("\t\t\t"))
}

// active grant
func readActiveGrantQuery(aka string) string {
	return fmt.Sprintf(`{
	activeGrant: activeGrant(id:"%s"){
%s
	}
}`, aka, turbotActiveGrantMetadataFragment("\t\t"))
}

func activateGrantMutation() string {
	return fmt.Sprintf(`mutation ActivateGrant($input: ActivateGrantInput!) {
	grantActivate: activateGrant(input: $input) {
%s
	}
}`, turbotActiveGrantMetadataFragment("\t\t\t"))
}

func deactivateGrantMutation() string {
	return fmt.Sprintf(`mutation DeactivateGrant($input: DeactivateGrantInput!) {
	deactivateGrant(input: $input) {
%s
	}
}`, turbotActiveGrantMetadataFragment("\t\t\t"))
}

// turbot directory
func createTurbotDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation createTurbotDirectory($input: CreateTurbotDirectoryInput!) {
 	 	resource: createTurbotDirectory(input: $input){
%s
    	turbot : get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

func updateTurbotDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation updateTurbotDirectory($input: UpdateTurbotDirectoryInput!) {
  		resource: updateTurbotDirectory(input: $input){
%s
		turbot : get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

func buildResourceProperties(resourceProperties []interface{}) string {
	var propertiesString bytes.Buffer
	for _, propertyPath := range resourceProperties {
		property, ok := propertyPath.(map[string]string)
		if ok {
			for alias, property := range property {
				propertiesString.WriteString(fmt.Sprintf("\t\t\t%s: get(path: \"%s\")\n", alias, property))
			}
		} else {
			propertiesString.WriteString(fmt.Sprintf("\t\t\t%s: get(path: \"%s\")\n", propertyPath, propertyPath))
		}

	}
	return propertiesString.String()
}

// local directory
func createLocalDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation createLocalDirectory($input: CreateLocalDirectoryInput!) {
 	 	resource: createLocalDirectory(input: $input){
%s
    	turbot : get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

func updateLocalDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation updateLocalDirectory($input: UpdateLocalDirectoryInput!) {
  		resource: updateLocalDirectory(input: $input){
%s
		turbot : get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

// google directory
func createGoogleDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation createGoogleDirectory($input: CreateGoogleDirectoryInput!) {
  		resource: createGoogleDirectory(input: $input){
%s
    	turbot:get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

func updateGoogleDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation updateGoogleDirectory($input: UpdateGoogleDirectoryInput!) {
  		resource: updateGoogleDirectory(input: $input){
%s
    	turbot:get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

// saml directory
func createSamlDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation createSamlDirectory($input: CreateSamlDirectoryInput!) {
  		resource: createSamlDirectory(input: $input){
%s
    	turbot:get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

func updateSamlDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation updateSamlDirectory($input: UpdateSamlDirectoryInput!) {
  		resource: updateSamlDirectory(input: $input){
%s
    	turbot:get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

//control
func readControlQuery(args string) string {
	return fmt.Sprintf(`{
control(%s){
	type{
		uri
	}
	state
	reason
	details
	turbot {
		id
		resourceId
	}
}
}`, args)
}

// group profile
func createGroupProfileMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation createGroupProfile($input: CreateGroupProfileInput!) {
  		resource: createGroupProfile(input: $input){
%s
    	turbot:get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

func updateGroupProfileMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation updateGroupProfile($input: UpdateGroupProfileInput!) {
  		resource: updateGroupProfile(input: $input){
%s
    	turbot:get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

func deleteGroupProfileMutation() string {
	return `mutation deleteGroupProfile($input: DeleteGroupProfileInput!) {
  		resource: deleteGroupProfile(input: $input){
    	turbot:get(path:"turbot")
  }
}`
}

// ldap directory
func createLdapDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation createLdapDirectory($input: CreateLdapDirectoryInput!) {
  		resource: createLdapDirectory(input: $input){
%s
    	turbot:get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

func updateLdapDirectoryMutation(properties []interface{}) string {
	return fmt.Sprintf(`mutation updateLdapDirectory($input: UpdateLdapDirectoryInput!) {
  		resource: updateLdapDirectory(input: $input){
%s
    	turbot:get(path:"turbot")
  }
}`, buildResourceProperties(properties))
}

func deleteLdapDirectory() string {
	return `mutation DeleteResource($input: DeleteResourceInput!) {
 	resource: deleteResource(input: $input) {
		turbot: get(path:"turbot")
	}
}`
}

// get turbot workspace version
func (client *Client) GetTurbotWorkspaceVersion() (*semver.Version, error) {
	query := readPolicyValueQuery("tmod:@turbot/turbot#/policy/types/workspaceVersion", "tmod:@turbot/turbot#/")
	responseData := &PolicyValueResponse{}
	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error reading policy value: %s", err.Error())
	}
	// convert interface {} to string
	versionValue := fmt.Sprintf("%v", responseData.PolicyValue.Value)
	// convert version value to semver value
	version, err := semver.New(versionValue)
	if err != nil {
		return nil, fmt.Errorf("error reading turbot workspace version value: %s", err.Error())
	}
	return version, nil
}
