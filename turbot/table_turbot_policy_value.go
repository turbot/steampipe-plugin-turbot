package turbot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableTurbotPolicyValue(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "turbot_policy_value",
		Description: "Policy value define the value of controls known to Turbot.",
		List: &plugin.ListConfig{
			Hydrate: listPolicyValue,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "state", Require: plugin.Optional},
				{Name: "policy_type_id", Require: plugin.Optional},
				{Name: "resource_id", Require: plugin.Optional},
				{Name: "resource_type_id", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ID"), Description: "Unique identifier of the policy type."},
			{Name: "policy_value_type_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type.Title"), Description: "Title of the policy value."},
			{Name: "policy_value_type_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type.Trunk.Title"), Description: "Title with full path of the policy value."},
			{Name: "default", Type: proto.ColumnType_BOOL, Description: "Defines the policy value is default or not."},
			{Name: "is_calculated", Type: proto.ColumnType_BOOL, Description: "True if this is a policy setting will be calculated for each value."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "State of the policy value."},
			{Name: "secret_value", Type: proto.ColumnType_STRING, Transform: transform.FromField("Value").Transform(convToString), Description: "Secrect value of the policy value."},
			{Name: "value", Type: proto.ColumnType_STRING, Transform: transform.FromField("Value").Transform(convToString), Description: "Value of the policy value."},
			{Name: "policy_value_type_mod_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type.ModURI"), Description: "URI of the mod that contains the policy value."},
			
			// Other columns
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Filter used for this policy setting list."},
			{Name: "resource_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ResourceId"), Description: "ID of the resource for the policy value."},
			{Name: "policy_type_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.PolicyTypeId"), Description: "ID of the policy type for this policy setting."},
			{Name: "resource_type_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ResourceTypeID"), Description: "ID of the resource type for this policy setting."},
			{Name: "setting_id", Type: proto.ColumnType_INT, Default: 0, Transform: transform.FromField("Turbot.SettingId").Transform(convStringToInt), Description: "Policy setting Id for the policy value."},
			{Name: "dependent_controls", Type: proto.ColumnType_JSON, Description: "The controls that depends upon the policy value."},
			{Name: "dependent_policy_values", Type: proto.ColumnType_JSON, Description: "The policy values that depends upon this policy value."},
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp"), Description: "When the policy value was first set by Turbot. (It may have been created earlier.)"},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.Timestamp"), Description: "Timestamp when the policy value was last modified (created, updated or deleted)."},
			{Name: "update_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.UpdateTimestamp"), Description: "When the policy value was last updated in Turbot."},
			{Name: "workspace", Type: proto.ColumnType_STRING, Hydrate: plugin.HydrateFunc(getTurbotWorkspace).WithCache(), Transform: transform.FromValue(), Description: "Specifies the workspace URL."},
		},
	}
}

const (
	queryPolicyValueList = `
query MyQuery($filter: [String!], $next_token: String) {
	policyValues(filter: $filter, paging: $next_token) {
		items {
			default
			value
			state
			reason
			details
			secretValue
			isCalculated
			type {
				modUri
				title
				trunk {
				  title
				}
			  }
			turbot {
				id
				policyTypeId
				resourceId
				resourceTypeId
				settingId
				createTimestamp
				deleteTimestamp
				timestamp
				updateTimestamp
			}
			dependentControls {
				items {
				turbot {
					controlTypeId
					controlTypePath
					controlCategoryId
					controlCategoryPath
					id
					resourceId
					resourceTypeId
				}
				type {
					modUri
					title
					trunk {
					title
					}
				}
				}
			}
			dependentPolicyValues {
				items {
				type {
					modUri
					uri
					title
					trunk {
					title
					}
					turbot {
					id
					title
					}
				}
				}
			}
			}
			paging {
			next
			}		  
	}
	}
`
)

func listPolicyValue(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_policy_type.listPolicyType", "connection_error", err)
		return nil, err
	}

	filters := []string{}
	quals := d.KeyColumnQuals

	// Additional filters
	if quals["state"] != nil {
		filters = append(filters, fmt.Sprintf("state:%s ", getQualListValues(ctx, quals, "state", "string")))
	}
	
	if quals["policy_type_id"] != nil {
		filters = append(filters, fmt.Sprintf("policyTypeId:%s policyTypeLevel:self", getQualListValues(ctx, quals, "policy_type_id", "string")))
	}

	if quals["resource_id"] != nil {
		filters = append(filters, fmt.Sprintf("resourceId:%s resourceTypeLevel:self", getQualListValues(ctx, quals, "resource_id", "string")))
	}

	if quals["resource_type_id"] != nil {
		filters = append(filters, fmt.Sprintf("resourceTypeId:%s resourceTypeLevel:self", getQualListValues(ctx, quals, "resource_type_id", "string")))
	}

	// Setting a high limit and page all results
	var pageLimit int64 = 5000

	// Adjust page limit, if less than default value
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < pageLimit {
			pageLimit = *limit
		}
	}

	// Setting page limit
	filters = append(filters, fmt.Sprintf("limit:%s", strconv.Itoa(int(pageLimit))))

	nextToken := ""
	for {
		result := &PolicyValuesResponse{}
		err = conn.DoRequest(queryPolicyValueList, map[string]interface{}{"filter": filters, "next_token": nextToken}, result)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_policy_type.listPolicyType", "query_error", err)
			return nil, err
		}
		for _, r := range result.PolicyValues.Items {
			d.StreamListItem(ctx, r)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if result.PolicyValues.Paging.Next == "" {
			break
		}
		nextToken = result.PolicyValues.Paging.Next
	}
	return nil, nil
}
