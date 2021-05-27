package turbot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-turbot/apiClient"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

/*

TODO - Need to support steampipe config, but default to Turbot terraform config

func connect(_ context.Context, d *plugin.QueryData) (*search.Client, error) {

	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "turbot"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*search.Client), nil
	}

	// Default to using env vars
	workspace := os.Getenv("TURBOT_WORKSPACE")
	accessKey := os.Getenv("TURBOT_ACCESS_KEY")
	secretKey := os.Getenv("TURBOT_SECRET_KEY")

	config := apiClient.ClientConfig{}

	// But prefer the config
	turbotConfig := GetConfig(d.Connection)
	if &turbotConfig != nil {
		if turbotConfig.Workspace != nil {
			config.
			workspace = *turbotConfig.Workspace
		}
		if turbotConfig.AccessKey != nil {
			accessKey = *turbotConfig.AccessKey
		}
		if turbotConfig.SecretKey != nil {
			secretKey = *turbotConfig.SecretKey
		}
	}

	if workspace == "" || accessKey == "" || secretKey == "" {
		// Credentials not set
		return nil, errors.New("workspace, access_key and secret_key must be configured")
	}

	conn := search.NewClient(appID, apiKey)

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, conn)

	return conn, nil
}

*/

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

func getMapValue(_ context.Context, d *transform.TransformData) (interface{}, error) {
	param := d.Param.(string)
	inputMap := d.Value.(map[string]interface{})
	if inputMap[param] != nil {
		return inputMap[param], nil
	}
	return "", nil
}

func emptyMapIfNil(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	v := d.Value.(map[string]interface{})
	plugin.Logger(ctx).Warn("emptyMapIfNil", "v", v)
	return v, nil
}

func emptyListIfNil(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	v := d.Value.([]string)
	plugin.Logger(ctx).Warn("emptyListIfNil", "v", v)
	return v, nil
}

func pathToArray(_ context.Context, d *transform.TransformData) (interface{}, error) {
	pathStr := types.SafeString(d.Value)
	pathStrs := strings.Split(pathStr, ".")
	pathInts := []int64{}
	for _, s := range pathStrs {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		pathInts = append(pathInts, i)
	}
	return pathInts, nil
}
