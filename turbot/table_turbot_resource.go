package turbot

import (
	"context"
	//"fmt"
	"strconv"
	//"strings"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableTurbotResource(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "turbot_resource",
		Description: "TODO",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("filter"),
			Hydrate:    listResource,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getResource,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ID"), Description: "Unique identifier of the resource."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Turbot.Title"), Description: "Title of the resource."},
			{Name: "tags", Type: proto.ColumnType_JSON, Transform: transform.FromField("Turbot.Tags"), Description: "Tags for the resource."},
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("Turbot.Akas"), Description: "AKA (also known as) identifiers for the resource."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.Timestamp"), Description: "Timestamp when the resource was last modified (created, updated or deleted)."},
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp"), Description: "When the resource was first discovered by Turbot. (It may have been created earlier.)"},
			{Name: "update_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.UpdateTimestamp"), Description: "When the resource was last updated in Turbot."},
			{Name: "version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.VersionID"), Description: "Unique identifier for this version of the resource."},
			{Name: "parent_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ParentID"), Description: "ID for the parent of this resource. For the Turbot root resource this is null."},
			{Name: "path", Type: proto.ColumnType_STRING, Transform: transform.FromField("Turbot.Path"), Description: "Hierarchy path with all identifiers of ancestors of the resource."},
			{Name: "resource_type_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ResourceTypeID"), Description: "ID of the resource type for this resource."},
			//{Name: "delete_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.DeleteTimestamp"), Description: "When the resource was deleted from Turbot."},
			{Name: "data", Type: proto.ColumnType_JSON, Description: "Resource data."},
			{Name: "metadata", Type: proto.ColumnType_JSON, Description: "Resource custom metadata."},
			{Name: "resource_type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type.URI"), Description: "URI of the resource type for this resource."},
			{Name: "filter", Type: proto.ColumnType_STRING, Hydrate: filterString, Transform: transform.FromValue(), Description: "Filter used for this resource list."},
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
			type {
				uri
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
		type {
			uri
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
	}
}
`
)

type ResourcesResponse struct {
	Resources struct {
		Items  []Resource
		Paging struct {
			Next string
		}
	}
}

type ResourceResponse struct {
	Resource Resource
}

type Resource struct {
	Turbot   TurbotResourceMetadata
	Data     map[string]interface{}
	Metadata map[string]interface{}
	Type     struct {
		URI string
	}
}

type TurbotResourceMetadata struct {
	ID                string
	ParentID          string
	Akas              []string
	Custom            map[string]interface{}
	Metadata          map[string]interface{}
	Tags              map[string]interface{}
	Title             string
	VersionID         string
	ActorIdentityID   string
	ActorPersonaID    string
	ActorRoleID       string
	ResourceParentAka string
	Timestamp         string
	CreateTimestamp   string
	DeleteTimestamp   string
	UpdateTimestamp   string
	Path              string
	ResourceGroupIDs  []string
	ResourceTypeID    string
	State             string
	Terraform         map[string]interface{}
}

func listResource(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_resource.listResource", "connection_error", err)
		return nil, err
	}

	/*
		filters := []string{"limit:5000"}
		allQuals := d.QueryContext.Quals
		plugin.Logger(ctx).Warn("listResource", "allQuals", allQuals)
		if allQuals["resource_type_id"] != nil {
			rti := allQuals["resource_type_id"]
			for _, q := range rti.Quals {
				plugin.Logger(ctx).Warn("listResource", "q", q)
				plugin.Logger(ctx).Warn("listResource", "q.GetFieldName()", q.GetFieldName())
				plugin.Logger(ctx).Warn("listResource", "q.GetOperator()", q.GetOperator())
				plugin.Logger(ctx).Warn("listResource", "q.GetValue().GetInt64Value()", q.GetValue().GetInt64Value())
				filters = append(filters, fmt.Sprintf("resourceTypeId:%d", q.GetValue().GetInt64Value()))
			}
		}
		filter := strings.Join(filters, " ")
		plugin.Logger(ctx).Warn("listResource", "filter", filter)
	*/

	quals := d.KeyColumnQuals
	filter := quals["filter"].GetStringValue()
	plugin.Logger(ctx).Warn("listResource", "filter", filter, "d", d)

	nextToken := ""

	for {
		result := &ResourcesResponse{}
		err = conn.DoRequest(queryResourceList, map[string]interface{}{"filter": filter, "next_token": nextToken}, result)
		plugin.Logger(ctx).Warn("listResource", "result", result, "next", result.Resources.Paging.Next, "err", err)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_resource.listResource", "query_error", err)
			return nil, err
		}
		for _, r := range result.Resources.Items {
			d.StreamListItem(ctx, r)
		}
		if result.Resources.Paging.Next == "" {
			break
		}
		nextToken = result.Resources.Paging.Next
	}

	return nil, nil
}

func getResource(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_resource.getResource", "connection_error", err)
		return nil, err
	}
	quals := d.KeyColumnQuals
	id := quals["id"].GetInt64Value()
	idStr := strconv.FormatInt(id, 10)
	plugin.Logger(ctx).Warn("getResource", "id", id, "id.str", idStr, "d", d)

	result := &ResourceResponse{}

	start := time.Now()
	err = conn.DoRequest(queryResourceGet, map[string]interface{}{"id": id}, result)
	plugin.Logger(ctx).Warn("getResource", "time", time.Since(start).Milliseconds(), "result", result, "err", err)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_resource.getResource", "query_error", err)
		return nil, err
	}
	return result.Resource, nil
}
