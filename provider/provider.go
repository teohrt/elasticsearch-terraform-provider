package provider

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/teohrt/terraform-provider-esslowlogconfig/client"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"es_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"esslowlogconfig_index": resourceIndexConfig(),
		},
		ConfigureFunc: providerClient,
	}
}

func providerClient(d *schema.ResourceData) (interface{}, error) {
	return client.New(&client.Config{
		Endpoint: d.Get("es_endpoint").(string),
		Region:   d.Get("region").(string),
		Creds:    credentials.NewEnvCredentials(),
	}), nil
}
