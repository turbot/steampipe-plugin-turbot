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

func tableTurbotActiveGrant(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "turbot_active_grant",
		Description: "All active grants of resources by Turbot.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "id", Require: plugin.Optional},
			},
			Hydrate: listActiveGrants,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ID"), Description: "Unique identifier of the grant."},
			{Name: "status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Grant.Identity.Status"), Description: "Status of the grant."},
			{Name: "display_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Grant.Identity.DisplayName"), Description: "Display name of the grant."},
			{Name: "email", Type: proto.ColumnType_STRING, Transform: transform.FromField("Grant.Identity.Email"), Description: "Email identity for the grant."},
			{Name: "family_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Grant.Identity.FamilyName"), Description: "Family name of the grant."},
			{Name: "given_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Grant.Identity.GivenName"), Description: "Given name of the grant."},
			{Name: "last_login_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Grant.Identity.LastLoginTimestamp"), Description: "Last login timestamp for the login."},
			{Name: "profile_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Grant.Identity.ProfileID"), Description: "Profile id of the grant."},
			{Name: "grant_identity_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Grant.Identity.Trunk.Title"), Description: "Full title (including ancestor trunk) of the grant identity."},
			{Name: "grant_type_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Grant.Type.Trunk.Title"), Description: "Full title (including ancestor trunk) of the grant type."},
			{Name: "resource_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Resource.Trunk.Title"), Description: "Full title (including ancestor trunk) of the resource."},
			{Name: "grant_type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Grant.Type.URI"), Description: "URI of the grant type."},
			{Name: "resource_type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Resource.Type.URI"), Description: "URI of the resource type."},
			{Name: "akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("Grant.Identity.Akas"), Description: "AKA (also known as) identifiers for the grant"},
			// Other columns
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp"), Description: "The create time of grant."},
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Filter used for this grant list."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.Timestamp"), Description: "Timestamp when the grant was last modified (created, updated or deleted)."},
			{Name: "update_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.UpdateTimestamp"), Description: "When the tag grant last updated in Turbot."},
			{Name: "version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.VersionID"), Description: "Unique identifier for this version of the grant."},
			{Name: "workspace", Type: proto.ColumnType_STRING, Hydrate: plugin.HydrateFunc(getTurbotWorkspace).WithCache(), Transform: transform.FromValue(), Description: "Specifies the workspace URL."},
		},
	}
}

const (
	grants = `
query activeGrantList($filter: [String!], $paging: String) {
	activeGrants(filter: $filter, paging: $paging) {
    items {
      resource {
        akas
        title
        trunk {
          title
        }
				type {
          uri
        }
      }
      grant {
        identity {
          title
          trunk {
            title
          }
          akas
          email: get(path: "email")
          status: get(path: "status")
          givenName: get(path: "givenName")
          profileId: get(path: "profileId")
          familyName: get(path: "familyName")
          displayName: get(path: "displayName")
          lastLoginTimestamp: get(path: "lastLoginTimestamp")
        }
        type {
          categoryUri
          category
          modUri
          trunk {
            title
          }
					uri
        }
      }
      turbot {
        id
        createTimestamp
        deleteTimestamp
        versionId
        timestamp
      }
    }
		paging {
      next
    }
  }
}
`
)

func listActiveGrants(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_active_grants.listActiveGrants", "connection_error", err)
		return nil, err
	}

	filters := []string{}
	quals := d.KeyColumnQuals

	filter := ""
	if quals["filter"] != nil {
		filter = quals["filter"].GetStringValue()
		filters = append(filters, filter)
	}

	// Additional filters
	if quals["id"] != nil {
		filters = append(filters, fmt.Sprintf("id:%s", getQualListValues(ctx, quals, "id", "int64")))
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
		var pageLimit int64 = 5000

		// Adjust page limit, if less than default value
		limit := d.QueryContext.Limit
		if d.QueryContext.Limit != nil {
			if *limit < pageLimit {
				pageLimit = *limit
			}
		}
		filters = append(filters, fmt.Sprintf("limit:%s", strconv.Itoa(int(pageLimit))))
	}

	plugin.Logger(ctx).Trace("turbot_active_grants.listActiveGrants", "quals", quals)
	plugin.Logger(ctx).Trace("turbot_active_grants.listActiveGrants", "filters", filters)

	nextToken := ""
	for {
		result := &ActiveGrantResponse{}
		err = conn.DoRequest(grants, map[string]interface{}{"filter": filters, "next_token": nextToken}, result)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_active_grants.listActiveGrants", "query_error", err)
		}
		for _, r := range result.ActiveGrants.Items {
			d.StreamListItem(ctx, r)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if !pageResults || result.ActiveGrants.Paging.Next == "" {
			break
		}
		nextToken = result.ActiveGrants.Paging.Next
	}

	return nil, nil
}
