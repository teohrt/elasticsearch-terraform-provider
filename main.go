package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/teohrt/terraform-provider-esslowlogconfig/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
