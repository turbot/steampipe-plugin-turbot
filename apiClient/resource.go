package apiClient

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/turbot/steampipe-plugin-turbot/errors"
	"github.com/turbot/steampipe-plugin-turbot/helpers"
	"log"
)

func (client *Client) CreateResource(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := createResourceMutation(nil)
	responseData := &CreateResourceResponse{}
	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "resource")
	}
	return &responseData.Resource.Turbot, nil
}

// properties is a map of terraform property name to turbot property path - it is used to add 'get' resolvers to the query
func (client *Client) ReadResource(resourceAka string, properties map[string]string) (*Resource, error) {
	var propertiesArray = []interface{}{properties}
	query := readResourceQuery(resourceAka, propertiesArray)
	var responseData = &ReadResourceResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, resourceAka, "resource")
	}

	resource, err := client.AssignResourceResults(responseData.Resource, properties)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (client *Client) ReadFullResource(resourceAka string) (*Resource, error) {
	query := readFullResourceQuery(resourceAka)
	var responseData = &ReadResourceResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, client.handleReadError(err, resourceAka, "resource")
	}

	resource, err := client.AssignResourceResults(responseData.Resource, nil)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// read a resource including all properties, then convert into a 'serializable' resource, consisting of simple types and string maps
func (client *Client) ReadSerializableResource(resourceAka string) (*SerializableResource, error) {
	// read the resource, passing an empty string as the property path in the properties map to force a full read
	properties := []interface{}{
		map[string]string{
			"data": "",
			"akas": "turbot.akas",
			"tags": "turbot.tags",
		},
	}

	query := readResourceQuery(resourceAka, properties)
	var responseData = &ReadSerializableResourceResponse{}

	// execute api call
	err := client.doRequest(query, nil, responseData)
	if err != nil {
		return nil, client.handleReadError(err, resourceAka, "resource")
	}
	resource := responseData.Resource

	// convert the data to JSON
	// (NOTE: remove the 'turbot' properties as this has been read separately)
	delete(resource.Data, "turbot")
	dataJson, err := helpers.MapToJsonString(resource.Data)
	if err != nil {
		return nil, err
	}
	// create a copy of the turbot object with all complex properties converted to JSON (as terraform schema cannot handle complex nested maps :/)

	// now convert to a map[string]string
	turbotStringMap, err := helpers.ConvertToStringMap(resource.Turbot)
	if err != nil {
		return nil, err
	}

	result := SerializableResource{
		Data:     dataJson,
		Turbot:   turbotStringMap,
		Tags:     resource.Tags,
		Akas:     resource.Akas,
		Metadata: turbotStringMap["custom"],
	}

	return &result, nil
}

func (client *Client) ReadResourceList(filter string, properties map[string]string) ([]Resource, error) {
	query := readResourceListQuery(filter, properties)
	var responseData = &ReadResourceListResponse{}

	// execute api call
	if err := client.doRequest(query, nil, responseData); err != nil {
		return nil, fmt.Errorf("error fetching resource list: %s", err.Error())
	}

	return responseData.Resources.Items, nil
}

func (client *Client) UpdateResource(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := updateResourceMutation(nil)
	responseData := &UpdateResourceResponse{}
	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleUpdateError(err, input, "resource")
	}
	return &responseData.Resource.Turbot, nil
}

func (client *Client) DeleteResource(aka string) error {
	query := deleteResourceMutation()
	// we do not care about the response
	var responseData interface{}

	variables := map[string]interface{}{
		"input": map[string]string{
			"id": aka,
		},
	}

	// execute api call
	if err := client.doRequest(query, variables, &responseData); err != nil {
		return fmt.Errorf("error deleting resource: %s", err.Error())
	}
	return nil
}

func (client *Client) ResourceExists(id string) (bool, error) {
	resource, err := client.ReadResource(id, nil)

	if err != nil {
		if errors.NotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	exists := resource.Turbot.Id != ""
	return exists, nil
}

func (client *Client) GetResourceAkas(resourceAka string) ([]string, error) {
	resource, err := client.ReadResource(resourceAka, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to load target resource; %s", err)
		return nil, err
	}
	resourceAkas := resource.Turbot.Akas
	// if this resource has no akas, just use the one passed in
	if resourceAkas == nil {
		resourceAkas = []string{resourceAka}
	}
	return resourceAkas, nil
}

// assign the ReadResource results into a Resource object, based on the 'properties' map
func (client *Client) AssignResourceResults(responseData interface{}, properties map[string]string) (*Resource, error) {
	var resource Resource
	// initialise map
	resource.Data = make(map[string]interface{})
	// convert turbot property to structure
	if err := mapstructure.Decode(responseData.(map[string]interface{})["turbot"], &resource.Turbot); err != nil {
		return nil, err
	}
	// convert type property
	if err := mapstructure.Decode(responseData.(map[string]interface{})["type"], &resource.Type); err != nil {
		return nil, err
	}
	// convert object property to structure
	if err := mapstructure.Decode(responseData.(map[string]interface{})["data"], &resource.Data); err != nil {
		return nil, err
	}
	// write properties into a map
	if properties != nil {
		for p := range properties {
			resource.Data[p] = responseData.(map[string]interface{})[p]
		}
	}

	return &resource, nil

}
