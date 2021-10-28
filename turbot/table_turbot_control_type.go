package turbot

import (
	"context"
	"fmt"
	"strconv"

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
			KeyColumns: []*plugin.KeyColumn{
				{Name: "category_uri", Require: plugin.Optional},
				{Name: "uri", Require: plugin.Optional},
			},
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
			{Name: "workspace", Type: proto.ColumnType_STRING, Hydrate: plugin.HydrateFunc(getTurbotWorkspace).WithCache(), Transform: transform.FromValue(), Description: "Specifies the workspace URL."},
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

	filters := []string{}
	quals := d.KeyColumnQuals

	// Additional filters
	if quals["uri"] != nil {
		filters = append(filters, fmt.Sprintf("controlTypeId:%s controlTypeLevel:self", getQualListValues(ctx, quals, "uri", "string")))
	}

	if quals["category_uri"] != nil {
		filters = append(filters, fmt.Sprintf("controlCategory:%s", getQualListValues(ctx, quals, "category_uri", "string")))
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

	plugin.Logger(ctx).Trace("turbot_control_type.listControlType", "quals", quals)
	plugin.Logger(ctx).Trace("turbot_control_type.listControlType", "filters", filters)

	nextToken := ""
	for {
		result := &ControlTypesResponse{}
		err = conn.DoRequest(queryControlTypeList, map[string]interface{}{"filter": filters, "next_token": nextToken}, result)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_control_type.listControlType", "query_error", err)
			return nil, err
		}
		for _, r := range result.ControlTypes.Items {
			d.StreamListItem(ctx, r)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
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
