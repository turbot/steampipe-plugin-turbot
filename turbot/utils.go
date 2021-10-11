package turbot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-turbot/apiClient"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
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
	if &turbotConfig != nil {
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

// QualOrFieldValue:: To get the values of fields, which are also used as optional qual
func QualOrFieldValue(_ context.Context, d *transform.TransformData) (interface{}, error) {
	// {Name: "resource_id", Type: proto.ColumnType_INT, Transform: transform.FromQual("resource_id").TransformP(QualOrFieldValue, "Turbot.ResourceID"), Description: "ID of the resource for this notification."},
	if d.Value != nil {
		return d.Value, nil
	}

	var item = d.HydrateItem
	var fieldNames []string

	switch p := d.Param.(type) {
	case []string:
		fieldNames = p
	case string:
		fieldNames = []string{p}
	default:
		return nil, fmt.Errorf("'FieldValue' requires one or more string parameters containing property path but received %v", d.Param)
	}

	for _, propertyPath := range fieldNames {
		fieldValue, ok := helpers.GetNestedFieldValueFromInterface(item, propertyPath)
		if ok {
			return fieldValue, nil

		}

	}
	return nil, nil
}
