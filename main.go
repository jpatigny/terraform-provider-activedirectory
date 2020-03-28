package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/jpatigny/terraform-provider-activedirectory/activedirectory"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: activedirectory.Provider,
	})
}
