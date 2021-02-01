package apiClient

import (
	"fmt"
)

func (client *Client) ReadControl(args string) (*Control, error) {
	query := readControlQuery(args)
	var responseData = &ReadControlResponse{}

	// execute api call
	err := client.doRequest(query, nil, responseData)
	if err != nil {
		return nil, fmt.Errorf("error reading control: %s", err.Error())
	}
	control := responseData.Control

	return &control, nil
}
