package turbot

import (
	"context"
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
			KeyColumns: plugin.SingleColumn("filter"),
			Hydrate:    listControl,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getControl,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ID"), Description: "Unique identifier of the control."},
			{Name: "state", Type: proto.ColumnType_STRING, Description: "State of the control."},
			{Name: "reason", Type: proto.ColumnType_STRING, Description: "Reason for this control state."},
			{Name: "details", Type: proto.ColumnType_JSON, Description: "Details associated with this control state."},
			{Name: "resource_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ResourceID"), Description: "ID of the resource this control is associated with."},
			{Name: "control_type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type.URI"), Description: "URI of the control type for this control."},
			// Other columns
			{Name: "control_type_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ControlTypeID"), Description: "ID of the control type for this control."},
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp"), Description: "When the control was first discovered by Turbot. (It may have been created earlier.)"},
			{Name: "filter", Type: proto.ColumnType_STRING, Hydrate: filterString, Transform: transform.FromValue(), Description: "Filter used for this control list."},
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
			type {
				uri
			}
			turbot {
				id
				timestamp
				createTimestamp
				updateTimestamp
				versionId
				controlTypeId
				resourceId
			}
		}
		paging {
			next
		}
	}
}
`

	queryControlGet = `
query controlGet($id: ID!) {
	control(id: $id) {
		state
		reason
		details
		type {
			uri
		}
		turbot {
			id
			timestamp
			createTimestamp
			updateTimestamp
			versionId
			controlTypeId
			resourceId
		}
	}
}
`
)

func listControl(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_control.listControl", "connection_error", err)
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
		result := &ControlsResponse{}
		err = conn.DoRequest(queryControlList, map[string]interface{}{"filter": filter, "next_token": nextToken}, result)
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

func getControl(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_control.getControl", "connection_error", err)
		return nil, err
	}
	quals := d.KeyColumnQuals
	id := quals["id"].GetInt64Value()
	result := &ControlResponse{}
	err = conn.DoRequest(queryControlGet, map[string]interface{}{"id": id}, result)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_control.getControl", "query_error", err)
		return nil, err
	}
	return result.Control, nil
}
