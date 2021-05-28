package turbot

import (
	"context"
	"fmt"
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
			KeyColumns: plugin.AnyColumn([]string{"id", "key", "value", "filter"}),
			Hydrate:    listTag,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ID"), Description: "Unique identifier of the tag."},
			{Name: "key", Type: proto.ColumnType_STRING, Description: "Tag key."},
			{Name: "value", Type: proto.ColumnType_STRING, Description: "Tag value."},
			{Name: "resource_ids", Type: proto.ColumnType_JSON, Transform: transform.FromField("Resources").Transform(tagResourcesToIdArray), Description: "Turbot IDs of resources with this tag."},
			// Other columns
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp"), Description: "When the tag was first discovered by Turbot. (It may have been created earlier.)"},
			{Name: "filter", Type: proto.ColumnType_STRING, Hydrate: filterString, Transform: transform.FromValue(), Description: "Filter used for this tag list."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.Timestamp"), Description: "Timestamp when the tag was last modified (created, updated or deleted)."},
			{Name: "update_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.UpdateTimestamp"), Description: "When the tag was last updated in Turbot."},
			{Name: "version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.VersionID"), Description: "Unique identifier for this version of the tag."},
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
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_tag.listTag", "connection_error", err)
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
		filters = append(filters, fmt.Sprintf("id:%d", quals["id"].GetInt64Value()))
	}
	if quals["key"] != nil {
		filters = append(filters, fmt.Sprintf("key:'%s'", escapeQualString(ctx, quals, "key")))
	}
	if quals["value"] != nil {
		filters = append(filters, fmt.Sprintf("value:'%s'", escapeQualString(ctx, quals, "value")))
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

	plugin.Logger(ctx).Trace("turbot_tag.listTag", "quals", quals)
	plugin.Logger(ctx).Trace("turbot_tag.listTag", "filters", filters)

	nextToken := ""
	for {
		result := &TagsResponse{}
		err = conn.DoRequest(queryTagList, map[string]interface{}{"filter": filters, "next_token": nextToken}, result)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_tag.listTag", "query_error", err)
			// TODO - this is a bit risk and should not be necessary, but there is a
			// bug in Turbot where sometimes resource requests within the tags table fail
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
