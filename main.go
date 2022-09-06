package main

import (
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-turbot/turbot"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: turbot.Plugin})
}
