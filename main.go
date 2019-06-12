package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/teohrt/terraform-provider-esdynamiconfig/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
