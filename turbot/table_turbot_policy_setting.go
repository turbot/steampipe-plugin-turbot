package turbot

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableTurbotPolicySetting(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "turbot_policy_setting",
		Description: "Policy settings defined in the Turbot workspace.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "id", Require: plugin.Optional},
				{Name: "resource_id", Require: plugin.Optional},
				{Name: "policy_type_id", Require: plugin.Optional},
				{Name: "policy_type_uri", Require: plugin.Optional},
				{Name: "filter", Require: plugin.Optional},
			},
			Hydrate: listPolicySetting,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ID"), Description: "Unique identifier of the policy setting."},
			{Name: "precedence", Type: proto.ColumnType_STRING, Description: "Precedence of the setting: REQUIRED or RECOMMENDED."},
			{Name: "resource_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ResourceID"), Description: "ID of the resource this policy setting is associated with."},
			{Name: "resource_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Resource.Trunk.Title"), Description: "Full title (including ancestor trunk) of the resource."},
			{Name: "policy_type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type.URI"), Description: "URI of the policy type for this policy setting."},
			{Name: "policy_type_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type.Trunk.Title"), Description: "Full title (including ancestor trunk) of the policy type."},
			{Name: "value", Type: proto.ColumnType_STRING, Description: "Value of the policy setting (for non-calculated policy settings)."},
			{Name: "is_calculated", Type: proto.ColumnType_BOOL, Description: "True if this is a policy setting will be calculated for each value."},
			{Name: "exception", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Exception").Transform(intToBool), Description: "True if this setting is an exception to a higher level setting."},
			{Name: "orphan", Type: proto.ColumnType_BOOL, Transform: transform.FromField("Orphan").Transform(intToBool), Description: "True if this setting is orphaned by a higher level setting."},
			{Name: "note", Type: proto.ColumnType_STRING, Description: "Optional note or comment for the setting."},
			// Other columns
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp"), Description: "When the policy setting was first discovered by Turbot. (It may have been created earlier.)"},
			{Name: "default", Type: proto.ColumnType_BOOL, Description: "True if this policy setting is the default."},
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Filter used for this policy setting list."},
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
			{Name: "workspace", Type: proto.ColumnType_STRING, Hydrate: plugin.HydrateFunc(getTurbotWorkspace).WithCache(), Transform: transform.FromValue(), Description: "Specifies the workspace URL."},
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
			resource {
				trunk {
					title
				}
			}
			#secretValue
			#secretValueSource
			template
			templateInput
			type {
				uri
				trunk {
					title
				}
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
)

func listPolicySetting(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_policy_setting.listPolicySetting", "connection_error", err)
		return nil, err
	}

	filters := []string{}
	quals := d.KeyColumnQuals

	filter := ""
	if quals["filter"] != nil {
		filter = quals["filter"].GetStringValue()
		filters = append(filters, filter)
	}

	// Additional filters
	if quals["id"] != nil {
		filters = append(filters, fmt.Sprintf("id:%s", getQualListValues(ctx, quals, "id", "int64")))
	}

	if quals["policy_type_id"] != nil {
		filters = append(filters, fmt.Sprintf("policyTypeId:%s policyTypeLevel:self", getQualListValues(ctx, quals, "policy_type_id", "int64")))
	}

	if quals["policy_type_uri"] != nil {
		filters = append(filters, fmt.Sprintf("policyTypeId:%s policyTypeLevel:self", getQualListValues(ctx, quals, "policy_type_uri", "string")))
	}

	if quals["resource_id"] != nil {
		filters = append(filters, fmt.Sprintf("resourceId:%s resourceTypeLevel:self", getQualListValues(ctx, quals, "resource_id", "int64")))
	}

	// Default to a very large page size. Page sizes earlier in the filter string
	// win, so this is only used as a fallback.
	pageResults := false
	// Add a limit if they haven't given one in the filter field
	re := regexp.MustCompile(`(^|\s)limit:[0-9]+($|\s)`)
	if !re.MatchString(filter) {
		// The caller did not specify a limit, so set a high limit and page all
		// results.
		pageResults = true
		var pageLimit int64 = 5000

		// Adjust page limit, if less than default value
		limit := d.QueryContext.Limit
		if d.QueryContext.Limit != nil {
			if *limit < pageLimit {
				pageLimit = *limit
			}
		}
		filters = append(filters, fmt.Sprintf("limit:%s", strconv.Itoa(int(pageLimit))))
	}

	plugin.Logger(ctx).Trace("turbot_policy_setting.listPolicySetting", "quals", quals)
	plugin.Logger(ctx).Trace("turbot_policy_setting.listPolicySetting", "filters", filters)

	nextToken := ""
	for {
		result := &PolicySettingsResponse{}
		err = conn.DoRequest(queryPolicySettingList, map[string]interface{}{"filter": filters, "next_token": nextToken}, result)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_policy_setting.listPolicySetting", "query_error", err)
			return nil, err
		}
		for _, r := range result.PolicySettings.Items {
			d.StreamListItem(ctx, r)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if !pageResults || result.PolicySettings.Paging.Next == "" {
			break
		}
		nextToken = result.PolicySettings.Paging.Next
	}

	return nil, nil
}
