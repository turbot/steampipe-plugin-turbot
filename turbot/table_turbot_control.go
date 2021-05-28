package turbot

import (
	"context"
	"fmt"
	"regexp"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableTurbotControl(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "turbot_control",
		Description: "Controls show the current state of checks in the Turbot workspace.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AnyColumn([]string{"id", "control_type_id", "control_type_uri", "resource_type_id", "resource_type_uri", "state", "filter"}),
			Hydrate:    listControl,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ID"), Description: "Unique identifier of the control."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "State of the control."},
			{Name: "reason", Type: proto.ColumnType_STRING, Description: "Reason for this control state."},
			{Name: "details", Type: proto.ColumnType_JSON, Description: "Details associated with this control state."},
			{Name: "resource_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ResourceID"), Description: "ID of the resource this control is associated with."},
			{Name: "resource_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Resource.Trunk.Title"), Description: "Full title (including ancestor trunk) of the resource."},
			{Name: "control_type_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type.Trunk.Title"), Description: "Full title (including ancestor trunk) of the control type."},
			// Other columns
			{Name: "control_type_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ControlTypeID"), Description: "ID of the control type for this control."},
			{Name: "control_type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type.URI"), Description: "URI of the control type for this control."},
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp"), Description: "When the control was first discovered by Turbot. (It may have been created earlier.)"},
			{Name: "filter", Type: proto.ColumnType_STRING, Hydrate: filterString, Transform: transform.FromValue(), Description: "Filter used for this control list."},
			{Name: "resource_type_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ResourceTypeID"), Description: "ID of the resource type for this control."},
			{Name: "resource_type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Resource.Type.URI"), Description: "URI of the resource type for this control."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.Timestamp"), Description: "Timestamp when the control was last modified (created, updated or deleted)."},
			{Name: "update_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.UpdateTimestamp"), Description: "When the control was last updated in Turbot."},
			{Name: "version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.VersionID"), Description: "Unique identifier for this version of the control."},
		},
	}
}

const (
	queryControlList = `
query controlList($filter: [String!], $next_token: String) {
	controls(filter: $filter, paging: $next_token) {
		items {
			state
			reason
			details
			resource {
				type {
					uri
				}
				trunk {
					title
				}
			}
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
				controlTypeId
				resourceId
				resourceTypeId
			}
		}
		paging {
			next
		}
	}
}
`
)

func listControl(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_control.listControl", "connection_error", err)
		return nil, err
	}

	filters := []string{}
	quals := d.KeyColumnQuals
	filter := ""
	if quals["filter"] != nil {
		filter = quals["filter"].GetStringValue()
		filters = append(filters, filter)
	}
	if quals["id"] != nil {
		filters = append(filters, fmt.Sprintf("id:%d", quals["id"].GetInt64Value()))
	}
	if quals["control_type_id"] != nil {
		filters = append(filters, fmt.Sprintf("controlTypeId:%d controlTypeLevel:self", quals["control_type_id"].GetInt64Value()))
	}
	if quals["control_type_uri"] != nil {
		filters = append(filters, fmt.Sprintf("controlTypeId:'%s' controlTypeLevel:self", escapeQualString(ctx, quals, "control_type_uri")))
	}
	if quals["resource_type_id"] != nil {
		filters = append(filters, fmt.Sprintf("resourceTypeId:%d resrouceTypeLevel:self", quals["resource_type_id"].GetInt64Value()))
	}
	if quals["resource_type_uri"] != nil {
		filters = append(filters, fmt.Sprintf("resourceTypeId:'%s' resourceTypeLevel:self", escapeQualString(ctx, quals, "resource_type_uri")))
	}
	if quals["state"] != nil {
		filters = append(filters, fmt.Sprintf("state:'%s'", escapeQualString(ctx, quals, "state")))
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
		filters = append(filters, "limit:5000")
	}

	plugin.Logger(ctx).Trace("turbot_control.listControl", "quals", quals)
	plugin.Logger(ctx).Trace("turbot_control.listControl", "filters", filters)

	nextToken := ""
	for {
		result := &ControlsResponse{}
		err = conn.DoRequest(queryControlList, map[string]interface{}{"filter": filters, "next_token": nextToken}, result)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_control.listControl", "query_error", err)
			return nil, err
		}
		for _, r := range result.Controls.Items {
			d.StreamListItem(ctx, r)
		}
		if !pageResults || result.Controls.Paging.Next == "" {
			break
		}
		nextToken = result.Controls.Paging.Next
	}

	return nil, nil
}
