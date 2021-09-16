package turbot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableTurbotResourceType(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "turbot_resource_type",
		Description: "Resource types define the types of resources known to Turbot.",
		List: &plugin.ListConfig{
			Hydrate: listResourceType,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "category_uri",
					Require: plugin.Optional,
				},
				{
					Name:    "uri",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getResourceType,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ID"), Description: "Unique identifier of the resource type."},
			{Name: "uri", Type: proto.ColumnType_STRING, Description: "URI of the resource type."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the resource type."},
			{Name: "trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Trunk.Title"), Description: "Title with full path of the resource type."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the resource type."},
			// Other columns
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("Turbot.Akas"), Description: "AKA (also known as) identifiers for the resource type."},
			{Name: "category_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Category.Turbot.ID"), Description: "ID of the resource category for the resource type."},
			{Name: "category_uri", Type: proto.ColumnType_STRING, Description: "URI of the resource category for the resource type."},
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp"), Description: "When the resource type was first discovered by Turbot. (It may have been created earlier.)"},
			{Name: "icon", Type: proto.ColumnType_STRING, Description: "Icon of the resource type."},
			{Name: "mod_uri", Type: proto.ColumnType_STRING, Description: "URI of the mod that contains the resource type."},
			{Name: "parent_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ParentID"), Description: "ID for the parent of this resource type."},
			{Name: "path", Type: proto.ColumnType_JSON, Transform: transform.FromField("Turbot.Path").Transform(pathToArray), Description: "Hierarchy path with all identifiers of ancestors of the resource type."},
			{Name: "update_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.UpdateTimestamp"), Description: "When the resource type was last updated in Turbot."},
			{Name: "version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.VersionID"), Description: "Unique identifier for this version of the resource type."},
			{Name: "workspace", Type: proto.ColumnType_STRING, Hydrate: plugin.HydrateFunc(getTurbotWorkspace).WithCache(), Transform: transform.FromValue(), Description: "Specifies the workspace URL."},
		},
	}
}

const (
	queryResourceTypeList = `
query resourceTypeList($filter: [String!], $next_token: String) {
	resourceTypes(filter: $filter, paging: $next_token) {
		items {
			category {
				turbot {
					id
				}
			}
			categoryUri
			description
			icon
			modUri
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

	queryResourceTypeGet = `
query resourceGet($id: ID!) {
	resourceType(id: $id) {
		category {
			turbot {
				id
			}
		}
		categoryUri
		description
		icon
		modUri
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
			title
			updateTimestamp
			versionId
		}
		uri
	}
}
`
)

func listResourceType(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_resource_type.listResourceType", "connection_error", err)
		return nil, err
	}

	filter := "limit:5000"
	nextToken := ""

	// Adjust page limit, if less than default value
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < 5000 {
			filter = fmt.Sprintf("limit:%s", strconv.Itoa(int(*limit)))
		}
	}

	// Additional filters
	if d.KeyColumnQuals["uri"] != nil {
		filter = filter + fmt.Sprintf(" resourceTypeId:'%s' resourceTypeLevel:self", d.KeyColumnQuals["uri"].GetStringValue())
	}

	if d.KeyColumnQuals["category_uri"] != nil {
		filter = filter + fmt.Sprintf(" resourceCategory:'%s'", d.KeyColumnQuals["category_uri"].GetStringValue())
	}

	for {
		result := &ResourceTypesResponse{}
		err = conn.DoRequest(queryResourceTypeList, map[string]interface{}{"filter": filter, "next_token": nextToken}, result)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_resource_type.listResourceType", "query_error", err)
			return nil, err
		}
		for _, r := range result.ResourceTypes.Items {
			d.StreamListItem(ctx, r)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if plugin.IsCancelled(ctx) {
				return nil, nil
			}
		}
		if result.ResourceTypes.Paging.Next == "" {
			break
		}
		nextToken = result.ResourceTypes.Paging.Next
	}

	return nil, nil
}

func getResourceType(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_resource_type.getResourceType", "connection_error", err)
		return nil, err
	}
	quals := d.KeyColumnQuals
	id := quals["id"].GetInt64Value()
	result := &ResourceTypeResponse{}
	err = conn.DoRequest(queryResourceTypeGet, map[string]interface{}{"id": id}, result)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_resource_type.getResourceType", "query_error", err)
		return nil, err
	}
	return result.ResourceType, nil
}
