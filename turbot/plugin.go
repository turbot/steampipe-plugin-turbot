package turbot

import (
	"context"

	"github.com/turbot/steampipe-plugin-turbot/errors"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
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
			"turbot_control":        tableTurbotControl(ctx),
			"turbot_control_type":   tableTurbotControlType(ctx),
			"turbot_policy_setting": tableTurbotPolicySetting(ctx),
			"turbot_policy_type":    tableTurbotPolicyType(ctx),
			"turbot_resource":       tableTurbotResource(ctx),
			"turbot_resource_type":  tableTurbotResourceType(ctx),
			"turbot_smart_folder":   tableTurbotSmartFolder(ctx),
			"turbot_tag":            tableTurbotTag(ctx),
		},
	}
	return p
}
