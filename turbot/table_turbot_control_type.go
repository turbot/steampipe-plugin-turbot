package turbot

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableTurbotControlType(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "turbot_control_type",
		Description: "Control types define the types of controls known to Turbot.",
		List: &plugin.ListConfig{
			Hydrate: listControlType,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getControlType,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ID"), Description: "Unique identifier of the control type."},
			{Name: "uri", Type: proto.ColumnType_STRING, Description: "URI of the control type."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the control type."},
			{Name: "trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Trunk.Title"), Description: "Title with full path of the control type."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the control type."},
			{Name: "targets", Type: proto.ColumnType_JSON, Description: "URIs of the resource types targeted by this control type."},
			// Other columns
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("Turbot.Akas"), Description: "AKA (also known as) identifiers for the control type."},
			{Name: "category_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Category.Turbot.ID"), Description: "ID of the control category for the control type."},
			{Name: "category_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Category.URI"), Description: "URI of the control category for the control type."},
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp"), Description: "When the control type was first discovered by Turbot. (It may have been created earlier.)"},
			{Name: "icon", Type: proto.ColumnType_STRING, Description: "Icon of the control type."},
			{Name: "mod_uri", Type: proto.ColumnType_STRING, Description: "URI of the mod that contains the control type."},
			{Name: "parent_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Turbot.ParentID"), Description: "ID for the parent of this control type."},
			{Name: "path", Type: proto.ColumnType_JSON, Transform: transform.FromField("Turbot.Path").Transform(pathToArray), Description: "Hierarchy path with all identifiers of ancestors of the control type."},
			// TODO - does not work {Name: "resource_target_ids", Type: proto.ColumnType_JSON, Description: "IDs of the resource types targeted by this control type."},
			{Name: "update_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.UpdateTimestamp"), Description: "When the control type was last updated in Turbot."},
			{Name: "version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.VersionID"), Description: "Unique identifier for this version of the control type."},
		},
	}
}

const (
	queryControlTypeList = `
query controlTypeList($filter: [String!], $next_token: String) {
	controlTypes(filter: $filter, paging: $next_token) {
		items {
			category {
				turbot {
					id
				}
				uri
			}
			description
			icon
			modUri
			targets
			title
			trunk {
				title
			}
			turbot {
				akas
				createTimestamp
				id
				parentId
				path
				#resourceTargetIds
				title
				updateTimestamp
				versionId
			}
			uri
		}
		paging {
			next
		}
	}
}
`

	queryControlTypeGet = `
query controlTypeGet($id: ID!) {
	controlType(id: $id) {
		category {
			turbot {
				id
			}
			uri
		}
		description
		icon
		modUri
		targets
		title
		trunk {
			title
		}
		turbot {
			akas
			createTimestamp
			id
			parentId
			path
			#resourceTargetIds
			title
			updateTimestamp
			versionId
		}
		uri
	}
}
`
)

func listControlType(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_control_type.listControlType", "connection_error", err)
		return nil, err
	}
	filter := "limit:5000"
	nextToken := ""
	for {
		result := &ControlTypesResponse{}
		err = conn.DoRequest(queryControlTypeList, map[string]interface{}{"filter": filter, "next_token": nextToken}, result)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_control_type.listControlType", "query_error", err)
			return nil, err
		}
		for _, r := range result.ControlTypes.Items {
			d.StreamListItem(ctx, r)
		}
		if result.ControlTypes.Paging.Next == "" {
			break
		}
		nextToken = result.ControlTypes.Paging.Next
	}

	return nil, nil
}

func getControlType(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_control_type.getControlType", "connection_error", err)
		return nil, err
	}
	quals := d.KeyColumnQuals
	id := quals["id"].GetInt64Value()
	result := &ControlTypeResponse{}
	err = conn.DoRequest(queryControlTypeGet, map[string]interface{}{"id": id}, result)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_control_type.getControlType", "query_error", err)
		return nil, err
	}
	return result.ControlType, nil
}