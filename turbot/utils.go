package turbot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"github.com/turbot/steampipe-plugin-turbot/apiClient"
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

func handleNilString(_ context.Context, d *transform.TransformData) (interface{}, error) {
	value := types.SafeString(d.Value)
	if value == "" {
		return "false", nil
	}
	return value, nil
}

func listResources(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData, resourceType string) (interface{}, error) {
	conn, err := connect(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_resource.listResource", "connection_error", err)
		return nil, err
	}

	filter := resourceType + " limit:500 resourceTypeLevel:self"
	plugin.Logger(ctx).Warn("listResources", "filter", filter, "d", d)

	nextToken := ""

	for {
		result := &ResourcesResponse{}
		err = conn.DoRequest(queryResourceList, map[string]interface{}{"filter": filter, "next_token": nextToken}, result)
		plugin.Logger(ctx).Warn("listResource", "result", result, "next", result.Resources.Paging.Next, "err", err)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_resource.listResource", "query_error", err)
			return nil, err
		}
		for _, r := range result.Resources.Items {
			r.Data["__turbot"] = r.Turbot
			r.Data["__metadata"] = r.Metadata

			d.StreamListItem(ctx, r.Data)
		}
		if result.Resources.Paging.Next == "" {
			break
		}
		nextToken = result.Resources.Paging.Next
	}

	return nil, nil
}

func getResourceById(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData, id int64) (interface{}, error) {
	conn, err := connect(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("getResourceById", "connection_error", err)
		return nil, err
	}

	idStr := strconv.FormatInt(id, 10)
	plugin.Logger(ctx).Warn("getResourceById", "id", id, "id.str", idStr, "d", d)

	result := &ResourceResponse{}

	start := time.Now()
	err = conn.DoRequest(queryResourceGet, map[string]interface{}{"id": id}, result)
	plugin.Logger(ctx).Warn("getResourceById", "time", time.Since(start).Milliseconds(), "result", result, "err", err)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_resource.getResource", "query_error", err)
		return nil, err
	}
	return result.Resource, nil //.(map[string]interface{}), nil
}
