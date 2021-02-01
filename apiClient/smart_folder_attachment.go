package apiClient

import (
	"fmt"
)

func (client *Client) CreateSmartFolderAttachment(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := createSmartFolderAttachmentMutation()
	responseData := &CreateSmartFolderAttachResponse{}

	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "smart folder attachment")
	}
	return &responseData.SmartFolderAttach.Turbot, nil
}

func (client *Client) DeleteSmartFolderAttachment(input map[string]interface{}) error {
	query := detachSmartFolderAttachment()
	var responseData interface{}

	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return fmt.Errorf("error deleting smart folder attachment: %s", err.Error())
	}
	return nil
}
