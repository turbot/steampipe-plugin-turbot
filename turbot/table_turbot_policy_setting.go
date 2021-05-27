package turbot

import (
	"context"
	"regexp"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableTurbotPolicySetting(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "turbot_policy_setting",
		Description: "Policy settings defined in the Turbot workspace.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("filter"),
			Hydrate:    listPolicySetting,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getPolicySetting,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ID"), Description: "Unique identifier of the policy setting."},
			{Name: "precedence", Type: proto.ColumnType_STRING, Description: "Precedence of the setting: REQUIRED or RECOMMENDED."},
			{Name: "resource_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ResourceID"), Description: "ID of the resource this policy setting is associated with."},
			{Name: "policy_type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type.URI"), Description: "URI of the policy type for this policy setting."},
			{Name: "value", Type: proto.ColumnType_STRING, Description: "Value of the policy setting (for non-calculated policy settings)."},
			{Name: "is_calculated", Type: proto.ColumnType_BOOL, Description: "True if this is a policy setting will be calculated for each value."},
			{Name: "exception", Type: proto.ColumnType_INT, Description: "Number of settings that this setting is an exception to. If zero, the setting is not an exception."},
			{Name: "orphan", Type: proto.ColumnType_INT, Description: "The number of settings that this setting is orphaned by. If zero, the setting is not an orphan."},
			{Name: "note", Type: proto.ColumnType_STRING, Description: "Optional note or comment for the setting."},
			// Other columns
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp"), Description: "When the policy setting was first discovered by Turbot. (It may have been created earlier.)"},
			{Name: "default", Type: proto.ColumnType_BOOL, Description: "True if this policy setting is the default."},
			{Name: "filter", Type: proto.ColumnType_STRING, Hydrate: filterString, Transform: transform.FromValue(), Description: "Filter used for this policy setting list."},
			{Name: "input", Type: proto.ColumnType_STRING, Description: "For calculated policy settings, this is the input GraphQL query."},
			{Name: "policy_type_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.PolicyTypeID"), Description: "ID of the policy type for this policy setting."},
			{Name: "template", Type: proto.ColumnType_STRING, Description: "For a calculated policy setting, this is the nunjucks template string defining a YAML string which is parsed to get the value."},
			{Name: "template_input", Type: proto.ColumnType_STRING, Description: "For calculated policy settings, this GraphQL query is run and used as input to the template."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.Timestamp"), Description: "Timestamp when the policy setting was last modified (created, updated or deleted)."},
			{Name: "update_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.UpdateTimestamp"), Description: "When the policy setting was last updated in Turbot."},
			{Name: "valid_from_timestamp", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the policy setting becomes valid."},
			{Name: "valid_to_timestamp", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the policy setting expires."},
			{Name: "value_source", Type: proto.ColumnType_STRING, Description: "The raw value in YAML format. If the setting was made via YAML template including comments, these will be included here."},
			{Name: "version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.VersionID"), Description: "Unique identifier for this version of the policy setting."},
		},
	}
}

const (
	queryPolicySettingList = `
query policySettingList($filter: [String!], $next_token: String) {
	policySettings(filter: $filter, paging: $next_token) {
		items {
			default
			exception
			input
			isCalculated
			note
			orphan
			precedence
			#secretValue
			#secretValueSource
			template
			templateInput
			type {
				uri
			}
			turbot {
				id
				timestamp
				createTimestamp
				updateTimestamp
				versionId
				policyTypeId
				resourceId
			}
			validFromTimestamp
			validToTimestamp
			value
			valueSource
		}
		paging {
			next
		}
	}
}
`

	queryPolicySettingGet = `
query policySettingGet($id: ID!) {
	policySetting(id: $id) {
		default
		exception
		input
		isCalculated
		note
		orphan
		precedence
		#secretValue
		#secretValueSource
		template
		templateInput
		type {
			uri
		}
		turbot {
			id
			timestamp
			createTimestamp
			updateTimestamp
			versionId
			policyTypeId
			resourceId
		}
		validFromTimestamp
		validToTimestamp
		value
		valueSource
	}
}
`
)

func listPolicySetting(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_policy_setting.listPolicySetting", "connection_error", err)
		return nil, err
	}

	quals := d.KeyColumnQuals
	filter := quals["filter"].GetStringValue()

	// Default to a very large page size. Page sizes earlier in the filter string
	// win, so this is only used as a fallback.
	pageResults := false
	re := regexp.MustCompile(`(^|\s)limit:[0-9]+($|\s)`)
	if !re.MatchString(filter) {
		// The caller did not specify a limit, so set a high limit and page all
		// results.
		pageResults = true
		filter = filter + " limit:5000"
	}

	nextToken := ""
	for {
		result := &PolicySettingsResponse{}
		err = conn.DoRequest(queryPolicySettingList, map[string]interface{}{"filter": filter, "next_token": nextToken}, result)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_policy_setting.listPolicySetting", "query_error", err)
			return nil, err
		}
		for _, r := range result.PolicySettings.Items {
			d.StreamListItem(ctx, r)
		}
		if !pageResults || result.PolicySettings.Paging.Next == "" {
			break
		}
		nextToken = result.PolicySettings.Paging.Next
	}

	return nil, nil
}

func getPolicySetting(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_policy_setting.getPolicySetting", "connection_error", err)
		return nil, err
	}
	quals := d.KeyColumnQuals
	id := quals["id"].GetInt64Value()
	result := &PolicySettingResponse{}
	err = conn.DoRequest(queryPolicySettingGet, map[string]interface{}{"id": id}, result)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_policy_setting.getPolicySetting", "query_error", err)
		return nil, err
	}
	return result.PolicySetting, nil
}
