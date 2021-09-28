package turbot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableTurbotSmartFolder(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "turbot_smart_folder",
		Description: "Smart folders allow policy settings to be attached as groups to resources.",
		List: &plugin.ListConfig{
			Hydrate: listSmartFolder,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getSmartFolder,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ID"), Description: "Unique identifier of the smart folder."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Turbot.Title"), Description: "Title of the smart folder."},
			{Name: "trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Trunk.Title"), Description: "Title with full path of the smart folder."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("Data").TransformP(getMapValue, "description"), Description: "Description of the smart folder."},
			{Name: "tags", Type: proto.ColumnType_JSON, Transform: transform.FromField("Turbot.Tags").Transform(emptyMapIfNil), Description: "Tags for the smart folder."},
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("Turbot.Akas").Transform(emptyListIfNil), Description: "AKA (also known as) identifiers for the smart folder."},
			{Name: "attached_resource_ids", Type: proto.ColumnType_JSON, Transform: transform.FromField("AttachedResources.Items").Transform(attachedResourceIDs), Description: ""},
			// Other columns
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp"), Description: "When the smart folder was first discovered by Turbot. (It may have been created earlier.)"},
			{Name: "color", Type: proto.ColumnType_STRING, Transform: transform.FromField("Data").TransformP(getMapValue, "color"), Description: "Color of the smart folder in the UI."},
			{Name: "data", Type: proto.ColumnType_JSON, Description: "Resource data."},
			{Name: "metadata", Type: proto.ColumnType_JSON, Description: "Resource custom metadata."},
			{Name: "parent_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ParentID"), Description: "ID for the parent of this smart folder."},
			{Name: "path", Type: proto.ColumnType_JSON, Transform: transform.FromField("Turbot.Path").Transform(pathToArray), Description: "Hierarchy path with all identifiers of ancestors of the smart folder."},
			{Name: "resource_type_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ResourceTypeID"), Description: "ID of the resource type for this smart folder."},
			{Name: "resource_type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type.URI"), Description: "URI of the resource type for this smart folder."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.Timestamp"), Description: "Timestamp when the smart folder was last modified (created, updated or deleted)."},
			{Name: "update_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.UpdateTimestamp"), Description: "When the smart folder was last updated in Turbot."},
			{Name: "version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.VersionID"), Description: "Unique identifier for this version of the smart folder."},
			{Name: "workspace", Type: proto.ColumnType_STRING, Hydrate: plugin.HydrateFunc(getTurbotWorkspace).WithCache(), Transform: transform.FromValue(), Description: "Specifies the workspace URL."},
		},
	}
}

const (
	querySmartFolderList = `
query smartFolderList($filter: [String!], $next_token: String) {
	resources(filter: $filter, paging: $next_token) {
		items {
			attachedResources {
				items {
					turbot {
						id
					}
				}
			}
			data
			metadata
			trunk {
				title
			}
			turbot {
				id
				title
				tags
				akas
				timestamp
				createTimestamp
				updateTimestamp
				versionId
				parentId
				path
				resourceTypeId
			}
			type {
				uri
			}
		}
		paging {
			next
		}
	}
}
`

	querySmartFolderGet = `
query smartFolderGet($id: ID!) {
	resource(id: $id) {
		attachedResources {
			items {
				turbot {
					id
				}
			}
		}
		data
		metadata
		trunk {
			title
		}
		turbot {
			id
			title
			tags
			akas
			timestamp
			createTimestamp
			updateTimestamp
			versionId
			parentId
			path
			resourceTypeId
		}
		type {
			uri
		}
	}
}
`
)

func listSmartFolder(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_smart_folder.listSmartFolder", "connection_error", err)
		return nil, err
	}

	var pageLimit int64 = 5000

	// Adjust page limit, if less than default value
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < pageLimit {
			pageLimit = *limit
		}
	}
	filter := fmt.Sprintf("resourceTypeId:'tmod:@turbot/turbot#/resource/types/smartFolder' resourceTypeLevel:self limit:%s", strconv.Itoa(int(pageLimit)))

	nextToken := ""
	for {
		result := &ResourcesResponse{}
		err = conn.DoRequest(querySmartFolderList, map[string]interface{}{"filter": filter, "next_token": nextToken}, result)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_smart_folder.listSmartFolder", "query_error", err)
			return nil, err
		}
		for _, r := range result.Resources.Items {
			d.StreamListItem(ctx, r)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if result.Resources.Paging.Next == "" {
			break
		}
		nextToken = result.Resources.Paging.Next
	}

	return nil, nil
}

func getSmartFolder(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_smart_folder.getSmartFolder", "connection_error", err)
		return nil, err
	}
	quals := d.KeyColumnQuals
	id := quals["id"].GetInt64Value()
	result := &ResourceResponse{}
	err = conn.DoRequest(querySmartFolderGet, map[string]interface{}{"id": id}, result)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_smart_folder.getSmartFolder", "query_error", err)
		return nil, err
	}
	return result.Resource, nil
}
