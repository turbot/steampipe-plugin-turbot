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

func tableTurbotGrant(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "turbot_grant",
		Description: "All grants of resources by Turbot.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "id", Require: plugin.Optional},
			},
			Hydrate: listActiveAndInactiveGrants,
		},
		Columns: grantColumns(),
	}
}

func grantColumns() []*plugin.Column {
	return []*plugin.Column{
		// Top columns
		{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.ID"), Description: "Unique identifier of the identity."},
		{Name: "identity_status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Identity.Status"), Description: "Status of the identity."},
		// {Name: "grant_status", Type: proto.ColumnType_STRING, Transform: transform.FromField("Status"), Description: "Status of the grant."},
		{Name: "identity_display_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Identity.DisplayName"), Description: "Display name of the identity."},
		{Name: "identity_email", Type: proto.ColumnType_STRING, Transform: transform.FromField("Identity.Email"), Description: "Email of the identity."},
		{Name: "ientity_family_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Identity.FamilyName"), Description: "Family name of the identity."},
		{Name: "identity_given_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Identity.GivenName"), Description: "Given name of the identity."},
		{Name: "identity_last_login_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Identity.LastLoginTimestamp"), Description: "Last login timestamp."},
		{Name: "identity_profile_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Identity.ProfileID"), Description: "Profile id of the identity."},
		{Name: "identity_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Identity.Trunk.Title"), Description: "Full title (including ancestor trunk) of the grant identity."},
		{Name: "level_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Level.Title"), Description: "The title of the level."},
		{Name: "level_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Level.Trunk.Title"), Description: "Full title (including ancestor trunk) of the level."},
		{Name: "level_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Level.URI"), Description: "The URI of the level."},
		{Name: "type_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type.Trunk.Title"), Description: "Full title (including ancestor trunk) of the grant type."},
		{Name: "type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Type.URI"), Description: "URI of the grant type."},
		{Name: "resource_trunk_title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Resource.Trunk.Title"), Description: "Full title (including ancestor trunk) of the resource."},
		{Name: "resource_type_uri", Type: proto.ColumnType_STRING, Transform: transform.FromField("Resource.Type.URI"), Description: "URI of the resource type."},
		{Name: "identity_akas", Type: proto.ColumnType_JSON, Transform: transform.FromField("Identity.Akas"), Description: "AKA (also known as) identifiers for the identity"},
		// Other columns
		{Name: "create_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.CreateTimestamp").NullIfEqual(""), Description: "The create time of the grant."},
		{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Filter used for this grant list."},
		{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.Timestamp").NullIfEqual(""), Description: "Timestamp when the grant was last modified (created, updated or deleted)."},
		{Name: "update_timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Turbot.UpdateTimestamp"), Description: "When the grant was last updated in Turbot."},
		{Name: "version_id", Type: proto.ColumnType_INT, Transform: transform.FromField("Turbot.VersionID").NullIfEqual(""), Description: "Unique identifier for this version of the identity."},
		{Name: "workspace", Type: proto.ColumnType_STRING, Hydrate: plugin.HydrateFunc(getTurbotWorkspace).WithCache(), Transform: transform.FromValue(), Description: "Specifies the workspace URL."},
	}
}

const (
	grants = `
query PermissionsByIdentity($filter: [String!], $paging: String) {
	permissionsDetails: permissionsDetailsByIdentity(filter: $filter, paging: $paging) {
		items {
			permissions {
				grants {
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
					}
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
					type {
						category
						categoryUri
						modUri
						trunk {
							title
						}
						uri
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
            createTimestamp
            deleteTimestamp
            profileId
            timestamp
            updateTimestamp
            versionId
					}
				}
				activeGrants {
					grant {
						turbot {
							id
						}
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

func listGrants(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) ([]Grant, []ActiveGrant, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_grants.listGrants", "connection_error", err)
		return nil, nil, err
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

	plugin.Logger(ctx).Trace("turbot_grants.listGrants", "quals", quals)
	plugin.Logger(ctx).Trace("turbot_grants.listGrants", "filters", filters)

	grantsInfo := []Grant{}
	activeGrantsDetails := []ActiveGrant{}
	nextToken := ""
	for {
		result := &PermissionByIdentityResponse{}
		err = conn.DoRequest(grants, map[string]interface{}{"filter": filters, "next_token": nextToken}, result)
		if err != nil {
			plugin.Logger(ctx).Error("turbot_grants.listGrants", "query_error", err)
		}
		for _, grantDetails := range result.PermissionsDetails.Items {
			for _, permission := range grantDetails.Permissions {
				activeGrantsDetails = append(activeGrantsDetails, permission.ActiveGrants...)
				grantsInfo = append(grantsInfo, permission.Grants...)
			}

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil, nil
			}
		}
		if !pageResults || result.PermissionsDetails.Paging.Next == "" {
			break
		}
		nextToken = result.PermissionsDetails.Paging.Next
	}

	return grantsInfo, activeGrantsDetails, err
}

func listActiveAndInactiveGrants(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	grants, _, err := listGrants(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("listActiveAndInactiveGrants", "Error", err)
		return nil, err
	}

	for _, grant := range grants {
		d.StreamListItem(ctx, grant)
	}
	return nil, nil
}
