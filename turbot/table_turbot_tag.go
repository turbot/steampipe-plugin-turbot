package turbot

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableTurbotTag(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "turbot_tag",
		Description: "TODO",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("filter"),
			Hydrate:    listTag,
		},
		/*
			Get: &plugin.GetConfig{
				KeyColumns: plugin.SingleColumn("id"),
				Hydrate:    getTag,
			},
		*/
		Columns: []*plugin.Column{
			// Top columns
			{Name: "key", Type: proto.ColumnType_STRING, Description: "Tag key."},
			{Name: "value", Type: proto.ColumnType_STRING, Description: "Tag value."},
			//{Name: "resources", Type: proto.ColumnType_JSON, Transform: transform.FromField("Resources").Transform(tagResourcesToIdArray), Description: "Resources with this tag."},
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ID"), Description: "Unique identifier of the tag."},
			{Name: "version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.VersionID"), Description: "Unique identifier for this version of the tag."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.Timestamp"), Description: "Timestamp when the tag was last modified (created, updated or deleted)."},
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp"), Description: "When the tag was first discovered by Turbot. (It may have been created earlier.)"},
			//{Name: "update_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.UpdateTimestamp"), Description: "When the tag was last updated in Turbot."},
			{Name: "filter", Type: proto.ColumnType_STRING, Hydrate: filterString, Transform: transform.FromValue(), Description: "Filter used for this tag list."},
		},
	}
}

const (
	queryTagList = `
query tagList($filter: [String!], $next_token: String) {
	tags(filter: $filter, paging: $next_token) {
		items {
			key
			value
			turbot {
				id
				timestamp
				createTimestamp
				updateTimestamp
				versionId
			}
			#resources {
			#	items {
			#		turbot {
			#			id
			#		}
			#	}
			#}
		}
		paging {
			next
		}
	}
}
`
)

type TagsResponse struct {
	Tags struct {
		Items  []Tag
		Paging struct {
			Next string
		}
	}
}

type Tag struct {
	Key       string
	Value     string
	Resources TagResources
	Turbot    TurbotTagMetadata
}

type TagResources struct {
	Items []TagResource
}

type TagResource struct {
	Turbot struct {
		ID string
	}
}

type TurbotTagMetadata struct {
	ID              string
	VersionID       string
	Timestamp       string
	CreateTimestamp string
	DeleteTimestamp string
	UpdateTimestamp string
}

func listTag(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_tag.listTag", "connection_error", err)
		return nil, err
	}
	quals := d.KeyColumnQuals
	filter := quals["filter"].GetStringValue()
	plugin.Logger(ctx).Warn("listTag", "filter", filter, "d", d)

	result := &TagsResponse{}
	nextToken := ""

	for {
		err = conn.DoRequest(queryTagList, map[string]interface{}{"filter": filter, "paging": nextToken}, result)
		plugin.Logger(ctx).Warn("listTag", "result", result, "next", result.Tags.Paging.Next, "err", err)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_tag.listTag", "query_error", err)
			return nil, err
		}
		for _, r := range result.Tags.Items {
			d.StreamListItem(ctx, r)
		}
		if result.Tags.Paging.Next == "" {
			break
		}
		nextToken = result.Tags.Paging.Next
	}

	return nil, nil
}

func tagResourcesToIdArray(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	resources := d.Value.(TagResources)
	ids := []string{}
	for _, r := range resources.Items {
		ids = append(ids, r.Turbot.ID)
	}
	return ids, nil
}
