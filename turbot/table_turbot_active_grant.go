package turbot

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableTurbotActiveGrant(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "turbot_active_grant",
		Description: "All active grants of resources by Turbot.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "id", Require: plugin.Optional},
			},
			Hydrate: listActivegrants,
		},
		Columns: grantColumns(),
	}
}

//// LIST FUNCTION
func listActivegrants(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listActivegrants")
	grants, activeGrants, err := listGrants(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("listActivegrants", "Error", err)
		return nil, err
	}

	plugin.Logger(ctx).Trace("Activegrants length", len(activeGrants))
	plugin.Logger(ctx).Trace("Grants Length", len(grants))


	for _, grant := range grants {
		grantStatus := getGrantStatus(grant, activeGrants)
		if grantStatus == "Active" {
			d.StreamListItem(ctx, grant)
		}
	}
	return nil, nil
}

//// TRANSFORM FUNCTION

func getGrantStatus(grant Grant, activeGrants []ActiveGrant) (status string) {
	status = "InActive"
	for _, activeGrantDetails := range activeGrants {
		if grant.Turbot.ID == activeGrantDetails.Grant.Turbot.ID {
			status = "Active"
			break
		}
	}
	return status
}
