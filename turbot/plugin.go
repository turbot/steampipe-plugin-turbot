package turbot

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-turbot",
		DefaultTransform: transform.FromGo().NullIfZero(),
		TableMap: map[string]*plugin.Table{
			"turbot_control":          tableTurbotControl(ctx),
			"turbot_resource":         tableTurbotResource(ctx),
			"turbot_tag":              tableTurbotTag(ctx),
			"turbot_aws_ec2_instance": tableTurbotAwsEc2Instance(ctx),
		},
	}
	return p
}
