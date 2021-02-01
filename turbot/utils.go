package turbot

import (
	"context"
	"fmt"
	"log"

	"github.com/turbot/steampipe-plugin-turbot/apiClient"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func connect(ctx context.Context) (*apiClient.Client, error) {
	/*
		config := apiClient.ClientConfig{
			Credentials: apiClient.ClientCredentials{
				AccessKey: d.Get("access_key").(string),
				SecretKey: d.Get("secret_key").(string),
				Workspace: d.Get("workspace").(string),
			},
			Profile:         d.Get("profile").(string),
			CredentialsPath: d.Get("credentials_file").(string),
		}
	*/
	config := apiClient.ClientConfig{}
	client, err := apiClient.CreateClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %s", err.Error())
	}
	log.Println("[INFO] Turbot API client initialized, now validating...", client)
	if err = client.Validate(); err != nil {
		return nil, err
	}
	return client, nil
}

func filterString(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	s := quals["filter"].GetStringValue()
	return s, nil
}
