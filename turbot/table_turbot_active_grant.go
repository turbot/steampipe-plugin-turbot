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
				{Name: "grant_id", Require: plugin.Optional},
			},
			Hydrate: listActiveGrants,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "grant_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Grant.Turbot.ID"), Description: "Unique identifier of the grant."},
			{Name: "resource_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Resource.Turbot.ID"), Description: "Unique identifier of the resource."},
			{Name: "identity_status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Grant.Identity.Status"), Description: "Status of the identity."},
			{Name: "identity_display_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Grant.Identity.DisplayName"), Description: "Display name of the identity."},
			{Name: "identity_email", Type: proto.ColumnType_STRING, Transform: transform.FromField("Grant.Identity.Email"), Description: "Email identity for the identity."},
			{Name: "ientity_family_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Grant.Identity.FamilyName"), Description: "Family name of the identity."},
			{Name: "identity_given_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Grant.Identity.GivenName"), Description: "Given name of the identity."},
			{Name: "identity_last_login_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Identity.LastLoginTimestamp"), Description: "Last login timestamp."},
			{Name: "identity_profile_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Grant.Identity.ProfileID"), Description: "Profile id of the identity."},
			{Name: "identity_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Grant.Identity.Trunk.Title"), Description: "Full title (including ancestor trunk) of the grant identity."},
			{Name: "level_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Grant.Level.Title"), Description: "The title of the level."},
			{Name: "level_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Grant.Level.Trunk.Title"), Description: "Full title (including ancestor trunk) of the level."},
			{Name: "level_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Grant.Level.URI"), Description: "The URI of the level."},
			{Name: "resource_type_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Resource.Type.Trunk.Title"), Description: "Full title (including ancestor trunk) of the grant type."},
			{Name: "resource_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Resource.Trunk.Title"), Description: "Full title (including ancestor trunk) of the resource."},
			{Name: "resource_type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Resource.Type.URI"), Description: "URI of the resource type."},
			{Name: "identity_akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("Grant.Identity.Akas"), Description: "AKA (also known as) identifiers for the identity"},
			// Other columns
			{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp").NullIfEqual(""), Description: "The create time of the grant."},
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Filter used for this grant list."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.Timestamp").NullIfEqual(""), Description: "Timestamp when the grant was last modified (created, updated or deleted)."},
			{Name: "update_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.UpdateTimestamp"), Description: "When the grant was last updated in Turbot."},
			{Name: "version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.VersionID").NullIfEqual(""), Description: "Unique identifier for this version of the identity."},
			{Name: "workspace", Type: proto.ColumnType_STRING, Hydrate: plugin.HydrateFunc(getTurbotWorkspace).WithCache(), Transform: transform.FromValue(), Description: "Specifies the workspace URL."},
		},
	}
}

const (
	activeGrants = `
	query MyQuery($filter: [String!], $paging: String) {
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
				trunk {
				  title
				}
			  }
			  turbot {
				id
				createTimestamp
				deleteTimestamp
				timestamp
				versionId
				updateTimestamp
			  }
			}
			grant {
			  identity {
				akas
				email: get(path: "email")
				status: get(path: "status")
				givenName: get(path: "givenName")
				profileId: get(path: "profileId")
				familyName: get(path: "familyName")
				displayName: get(path: "displayName")
				lastLoginTimestamp: get(path: "lastLoginTimestamp")
				trunk {
				  title
				}
			  }
			  level {
				title
				uri
				trunk {
				  title
				}
			  }
			  turbot {
				id
			  }
			}
			turbot {
			  createTimestamp
			  deleteTimestamp
			  updateTimestamp
			  title
			  timestamp
			  versionId
			}
		  }
		  paging {
			next
		  }
		}
	  }		 
`
)

func listActiveGrants(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
	if quals["grant_id"] != nil {
		filters = append(filters, fmt.Sprintf("id:%s", getQualListValues(ctx, quals, "grant_id", "int64")))
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

	nextToken := ""
	for {
		result := &ActiveGrantInfo{}
		err = conn.DoRequest(activeGrants, map[string]interface{}{"filter": filters, "next_token": nextToken}, result)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_active_grants.listActiveGrants", "query_error", err)
		}
		for _, ActiveGrantDetails := range result.ActiveGrants.Items {

			d.StreamListItem(ctx, ActiveGrantDetails)
			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if !pageResults || result.ActiveGrants.Paging.Next== "" {
			break
		}
		nextToken = result.ActiveGrants.Paging.Next
	}

	return nil, err
}
