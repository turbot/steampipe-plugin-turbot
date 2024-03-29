package turbot

import (
	"context"

	"github.com/turbot/steampipe-plugin-turbot/errors"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-turbot",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultGetConfig: &plugin.GetConfig{
			ShouldIgnoreError: errors.NotFoundError,
		},
		DefaultTransform: transform.FromGo(),
		TableMap: map[string]*plugin.Table{
			"turbot_active_grant":   tableTurbotActiveGrant(ctx),
			"turbot_control":        tableTurbotControl(ctx),
			"turbot_control_type":   tableTurbotControlType(ctx),
			"turbot_grant":          tableTurbotGrant(ctx),
			"turbot_mod_version":    tableTurbotModVersion(ctx),
			"turbot_notification":   tableTurbotNotification(ctx),
			"turbot_policy_setting": tableTurbotPolicySetting(ctx),
			"turbot_policy_type":    tableTurbotPolicyType(ctx),
			"turbot_policy_value":   tableTurbotPolicyValue(ctx),
			"turbot_resource":       tableTurbotResource(ctx),
			"turbot_resource_type":  tableTurbotResourceType(ctx),
			"turbot_smart_folder":   tableTurbotSmartFolder(ctx),
			"turbot_tag":            tableTurbotTag(ctx),
		},
	}
	return p
}
