package turbot

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableTurbotModVersion(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "turbot_mod_version",
		Description: "Module versions in turbot organization.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "name", Require: plugin.Optional},
				{Name: "status", Require: plugin.Optional},
				{Name: "filter", Require: plugin.Optional},
				{Name: "org_name", Require: plugin.Optional},
			},
			Hydrate: listModVersion,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the mod."},
			{Name: "identity_name", Type: proto.ColumnType_STRING, Description: "The indentity name of the mod."},
			{Name: "org_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("org_name"), Description: "The name of the organization."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the mod version."},
			{Name: "version", Type: proto.ColumnType_STRING, Description: "The version of the mod."},
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Filter used to search for mod versions."},
			{Name: "mod_peer_dependency", Type: proto.ColumnType_JSON, Transform: transform.FromField("Head.PeerDependencies"), Description: "Peer dependencies of the mod."},
			// Other columns
			{Name: "workspace", Type: proto.ColumnType_STRING, Hydrate: plugin.HydrateFunc(getTurbotWorkspace).WithCache(), Transform: transform.FromValue(), Description: "Specifies the workspace URL."},
		},
	}
}

type ModVersionInfo struct {
	IdentityName string
	Name         string
	Status       string
	Version      string
	Head         ModVersionHead
}

const (
	queryModVersions = `
query modVersionSearchByName($search: String, $modName: String, $orgName: String, $status: [ModVersionStatus!]) {
  modVersionSearches(search: $search, modName: $modName, orgName: $orgName, status: $status) {
	items {
	identityName
	name
	versions {
		version
		status
		head
	}
	}
	paging {
	next
	}
  }
}
`
)

func listModVersion(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("turbot_mod_version.listModVersion", "connection_error", err)
		return nil, err
	}

	var status interface{}
	var modName, searchText, orgName string

	quals := d.KeyColumnQuals

	// Additional filters
	if quals["status"] != nil {
		status = strings.ToUpper(quals["status"].GetStringValue())
	}
	if quals["name"] != nil {
		modName = quals["name"].GetStringValue()
	}
	if quals["filter"] != nil {
		searchText = quals["filter"].GetStringValue()
	}
	if quals["org_name"] != nil {
		orgName = quals["org_name"].GetStringValue()
	}

	plugin.Logger(ctx).Trace("turbot_mod_version.listModVersion", "quals", quals)
	nextToken := ""

	for {
		result := &ModVersionResponse{}
		if status != nil {
			err = conn.DoRequest(queryModVersions, map[string]interface{}{"search": searchText, "orgName": orgName, "modName": modName, "status": status, "next_token": nextToken}, result)
		} else {
			err = conn.DoRequest(queryModVersions, map[string]interface{}{"search": searchText, "orgName": orgName, "modName": modName, "next_token": nextToken}, result)
		}

		if err != nil {
			plugin.Logger(ctx).Error("turbot_mod_version.listModVersion", "query_error", err)
			return nil, err
		}
		for _, r := range result.ModVersionSearches.Items {

			for _, resp := range r.Versions {
				d.StreamListItem(ctx, ModVersionInfo{
					IdentityName: r.IdentityName,
					Name:         r.Name,
					Status:       resp.Status,
					Version:      resp.Version,
					Head:         resp.Head,
				})
			}

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if result.ModVersionSearches.Paging.Next == "" {
			break
		}
	}

	return nil, nil
}
