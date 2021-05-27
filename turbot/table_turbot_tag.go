package turbot

import (
	"context"
	"regexp"
	"strconv"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableTurbotTag(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "turbot_tag",
		Description: "All tags discovered on cloud resources by Turbot.",
		List: &plugin.ListConfig{
			//KeyColumns: plugin.SingleColumn("filter"),
			Hydrate: listTag,
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
			{Name: "resources", Type: proto.ColumnType_JSON, Transform: transform.FromField("Resources").Transform(tagResourcesToIdArray), Description: "Resources with this tag."},
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ID"), Description: "Unique identifier of the tag."},
			{Name: "version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.VersionID"), Description: "Unique identifier for this version of the tag."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.Timestamp"), Description: "Timestamp when the tag was last modified (created, updated or deleted)."},
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp"), Description: "When the tag was first discovered by Turbot. (It may have been created earlier.)"},
			//{Name: "resources", Type: proto.ColumnType_JSON, Description: ""},
			{Name: "update_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.UpdateTimestamp"), Description: "When the tag was last updated in Turbot."},
			{Name: "filter", Type: proto.ColumnType_STRING, Hydrate: filterString, Transform: transform.FromValue(), Description: "Filter used for this tag list."},
		},
	}
}

const (
	queryTagList = `
query tagList($filter: [String!], $paging: String) {
	tags(filter: $filter, paging: $paging) {
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
			resources {
				items {
					turbot {
						id
					}
				}
			}
		}
		paging {
			next
		}
	}
}
`
)

func listTag(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_tag.listTag", "connection_error", err)
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
		result := &TagsResponse{}
		err = conn.DoRequest(queryTagList, map[string]interface{}{"filter": filter, "next_token": nextToken}, result)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_tag.listTag", "query_error", err)
			// TODO - this should not be necessary, but there is a bug where sometimes resource requests within the tags table fail
			//return nil, err
		}
		for _, r := range result.Tags.Items {
			d.StreamListItem(ctx, r)
		}
		if !pageResults || result.Tags.Paging.Next == "" {
			break
		}
		nextToken = result.Tags.Paging.Next
	}

	return nil, nil
}

func tagResourcesToIdArray(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	resources := d.Value.(TagResources)
	ids := []int64{}
	for _, r := range resources.Items {
		id, err := strconv.ParseInt(r.Turbot.ID, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
