package turbot

import (
	"context"
	"strconv"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableTurbotControl(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "turbot_control",
		Description: "TODO",
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
			{Name: "control_type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type.URI"), Description: "URI of the control type for this control."},

			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.Timestamp"), Description: "Timestamp when the control was last modified (created, updated or deleted)."},
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp"), Description: "When the control was first discovered by Turbot. (It may have been created earlier.)"},
			{Name: "update_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.UpdateTimestamp"), Description: "When the control was last updated in Turbot."},
			{Name: "version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.VersionID"), Description: "Unique identifier for this version of the control."},
			{Name: "control_type_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ControlTypeID"), Description: "ID of the control type for this control."},
			{Name: "resource_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ResourceID"), Description: "ID of the resource this control is associated with."},
			//{Name: "delete_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.DeleteTimestamp"), Description: "When the control was deleted from Turbot."},
			{Name: "filter", Type: proto.ColumnType_STRING, Hydrate: filterString, Transform: transform.FromValue(), Description: "Filter used for this control list."},
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

type ControlsResponse struct {
	Controls struct {
		Items  []Control
		Paging struct {
			Next string
		}
	}
}

type ControlResponse struct {
	Control Control
}

type Control struct {
	State   string
	Reason  string
	Details interface{}
	Type    struct {
		URI string
	}
	Turbot TurbotControlMetadata
}

type TurbotControlMetadata struct {
	ID              string
	VersionID       string
	Timestamp       string
	CreateTimestamp string
	DeleteTimestamp string
	UpdateTimestamp string
	ControlTypeID   string
	ResourceID      string
}

func listControl(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_control.listControl", "connection_error", err)
		return nil, err
	}
	quals := d.KeyColumnQuals
	filter := quals["filter"].GetStringValue()
	plugin.Logger(ctx).Warn("listControl", "filter", filter, "d", d)

	nextToken := ""

	for {
		result := &ControlsResponse{}
		err = conn.DoRequest(queryControlList, map[string]interface{}{"filter": filter, "next_token": nextToken}, result)
		plugin.Logger(ctx).Warn("listControl", "result", result, "next", result.Controls.Paging.Next, "err", err)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_control.listControl", "query_error", err)
			return nil, err
		}
		for _, r := range result.Controls.Items {
			d.StreamListItem(ctx, r)
		}
		if result.Controls.Paging.Next == "" {
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
	idStr := strconv.FormatInt(id, 10)
	plugin.Logger(ctx).Warn("getControl", "id", id, "id.str", idStr, "d", d)

	result := &ControlResponse{}

	start := time.Now()
	err = conn.DoRequest(queryControlGet, map[string]interface{}{"id": id}, result)
	plugin.Logger(ctx).Warn("getControl", "time", time.Since(start).Milliseconds(), "result", result, "err", err)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_control.getControl", "query_error", err)
		return nil, err
	}
	return result.Control, nil
}
