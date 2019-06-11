package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/teohrt/terraform-provider-elasticsearch-index-thresholds/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
