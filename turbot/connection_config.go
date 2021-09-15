package turbot

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type turbotConfig struct {
	Profile   *string `cty:"profile"`
	AccessKey *string `cty:"access_key"`
	SecretKey *string `cty:"secret_key"`
	Workspace *string `cty:"workspace_url"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"profile": {
		Type: schema.TypeString,
	},
	"access_key": {
		Type: schema.TypeString,
	},
	"secret_key": {
		Type: schema.TypeString,
	},
	"workspace_url": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &turbotConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) turbotConfig {
	if connection == nil || connection.Config == nil {
		return turbotConfig{}
	}
	config, _ := connection.Config.(turbotConfig)
	return config
}
