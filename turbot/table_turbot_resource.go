package turbot

import (
	"context"
	"fmt"
	"regexp"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableTurbotResource(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "turbot_resource",
		Description: "Resources from the Turbot CMDB.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AnyColumn([]string{"id", "resource_type_id", "resource_type_uri", "filter"}),
			Hydrate:    listResource,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ID"), Description: "Unique identifier of the resource."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Turbot.Title"), Description: "Title of the resource."},
			{Name: "trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Trunk.Title"), Description: "Title with full path of the resource."},
			{Name: "tags", Type: proto.ColumnType_JSON, Transform: transform.FromField("Turbot.Tags"), Description: "Tags for the resource."},
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("Turbot.Akas"), Description: "AKA (also known as) identifiers for the resource."},
			// Other columns
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp"), Description: "When the resource was first discovered by Turbot. (It may have been created earlier.)"},
			{Name: "data", Type: proto.ColumnType_JSON, Description: "Resource data."},
			{Name: "filter", Type: proto.ColumnType_STRING, Hydrate: filterString, Transform: transform.FromValue(), Description: "Filter used for this resource list."},
			{Name: "metadata", Type: proto.ColumnType_JSON, Description: "Resource custom metadata."},
			{Name: "parent_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ParentID"), Description: "ID for the parent of this resource. For the Turbot root resource this is null."},
			{Name: "path", Type: proto.ColumnType_JSON, Transform: transform.FromField("Turbot.Path").Transform(pathToArray), Description: "Hierarchy path with all identifiers of ancestors of the resource."},
			{Name: "resource_type_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ResourceTypeID"), Description: "ID of the resource type for this resource."},
			{Name: "resource_type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type.URI"), Description: "URI of the resource type for this resource."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.Timestamp"), Description: "Timestamp when the resource was last modified (created, updated or deleted)."},
			{Name: "update_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.UpdateTimestamp"), Description: "When the resource was last updated in Turbot."},
			{Name: "version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.VersionID"), Description: "Unique identifier for this version of the resource."},
		},
	}
}

const (
	queryResourceList = `
query resourceList($filter: [String!], $next_token: String) {
	resources(filter: $filter, paging: $next_token) {
		items {
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

	queryResourceGet = `
query resourceGet($id: ID!) {
	resource(id: $id) {
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

func listResource(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_resource.listResource", "connection_error", err)
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
		filters = append(filters, fmt.Sprintf("resourceId:%d level:self", quals["id"].GetInt64Value()))
	}
	if quals["resource_type_id"] != nil {
		filters = append(filters, fmt.Sprintf("resourceTypeId:%d resourceTypeLevel:self", quals["resource_type_id"].GetInt64Value()))
	}
	if quals["resource_type_uri"] != nil {
		filters = append(filters, fmt.Sprintf("resourceTypeId:'%s' resourceTypeLevel:self", escapeQualString(ctx, quals, "resource_type_uri")))
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

	nextToken := ""
	for {
		result := &ResourcesResponse{}
		err = conn.DoRequest(queryResourceList, map[string]interface{}{"filter": filters, "next_token": nextToken}, result)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_resource.listResource", "query_error", err)
			return nil, err
		}
		for _, r := range result.Resources.Items {
			d.StreamListItem(ctx, r)
		}
		if !pageResults || result.Resources.Paging.Next == "" {
			break
		}
		nextToken = result.Resources.Paging.Next
	}

	return nil, nil
}
