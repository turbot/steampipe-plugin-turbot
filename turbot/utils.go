package turbot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-turbot/apiClient"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const (
	filterTimeFormat = "2006-01-02T15:04:05.000Z"
)

func connect(ctx context.Context, d *plugin.QueryData) (*apiClient.Client, error) {

	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "turbot"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*apiClient.Client), nil
	}

	// Start with an empty Turbot config
	config := apiClient.ClientConfig{Credentials: apiClient.ClientCredentials{}}

	// Prefer config options given in Steampipe
	turbotConfig := GetConfig(d.Connection)
	if turbotConfig.Profile != nil {
		config.Profile = *turbotConfig.Profile
	}
	if turbotConfig.Workspace != nil {
		config.Credentials.Workspace = *turbotConfig.Workspace
	}
	if turbotConfig.AccessKey != nil {
		config.Credentials.AccessKey = *turbotConfig.AccessKey
	}
	if turbotConfig.SecretKey != nil {
		config.Credentials.SecretKey = *turbotConfig.SecretKey
	}

	// Create the client
	client, err := apiClient.CreateClient(config)
	if err != nil {
		return nil, fmt.Errorf("Error creating Turbot client: %s", err.Error())
	}
	if err = client.Validate(); err != nil {
		return nil, fmt.Errorf("Error validating Turbot client: %s", err.Error())
	}

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, client)

	// Done
	return client, nil
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
	return v, nil
}

func emptyListIfNil(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	v := d.Value.([]string)
	return v, nil
}

func intToBool(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}
	v := d.Value.(int)
	return v > 0, nil
}

func convToString(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	var v interface{} = fmt.Sprint(d.Value)
	return v, nil
}

func attachedResourceIDs(_ context.Context, d *transform.TransformData) (interface{}, error) {
	objs := d.Value.([]TurbotIDObject)
	ids := []int64{}
	for _, o := range objs {
		id, err := strconv.ParseInt(o.Turbot.ID, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids, nil
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

func escapeQualString(_ context.Context, quals map[string]*proto.QualValue, qualName string) string {
	s := quals[qualName].GetStringValue()
	s = strings.Replace(s, "\\", "\\\\", -1)
	s = strings.Replace(s, "'", "\\'", -1)
	return s
}

func getTurbotWorkspace(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Load workspace name from cache
	cacheKey := "getTurbotWorkspaceInfo"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(string), nil
	}

	// Start with an empty Turbot config
	config := apiClient.ClientConfig{Credentials: apiClient.ClientCredentials{}}

	// Prefer config options given in Steampipe
	turbotConfig := GetConfig(d.Connection)
	if turbotConfig.Profile != nil {
		config.Profile = *turbotConfig.Profile
	}
	if turbotConfig.Workspace != nil {
		config.Credentials.Workspace = *turbotConfig.Workspace
	}
	if turbotConfig.AccessKey != nil {
		config.Credentials.AccessKey = *turbotConfig.AccessKey
	}
	if turbotConfig.SecretKey != nil {
		config.Credentials.SecretKey = *turbotConfig.SecretKey
	}

	credentials, err := apiClient.GetCredentials(config)
	if err != nil {
		return nil, nil
	}
	endpoint := credentials.Workspace // https://pikachu-turbot.cloud.turbot-dev.com/api/latest/graphql
	if endpoint != "" {
		workspaceUrl := strings.Split(endpoint, "/api/")[0]
		return workspaceUrl, nil
	}

	return nil, nil
}

// Get QualValueList as an list of items
func getQualListValues(ctx context.Context, quals map[string]*proto.QualValue, qualName string, qualType string) string {
	switch qualType {
	case "string":
		if quals[qualName].GetStringValue() != "" {
			return fmt.Sprintf("'%s'", escapeQualString(ctx, quals, qualName))
		} else if quals[qualName].GetListValue() != nil {
			values := make([]string, 0)
			for _, value := range quals[qualName].GetListValue().Values {
				str := value.GetStringValue()
				str = strings.Replace(str, "\\", "\\\\", -1)
				str = strings.Replace(str, "'", "\\'", -1)
				values = append(values, fmt.Sprintf("'%s'", str))
			}
			return strings.Join(values, ",")
		}
	case "int64":
		if quals[qualName].GetInt64Value() != 0 {
			return strconv.FormatInt(quals[qualName].GetInt64Value(), 10)
		} else if quals[qualName].GetListValue() != nil {
			values := make([]string, 0)
			for _, value := range quals[qualName].GetListValue().Values {
				values = append(values, strconv.FormatInt(value.GetInt64Value(), 10))
			}
			return strings.Join(values, ",")
		}
	}
	return ""
}
