package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/teohrt/terraform-provider-esdynamiconfig/client"
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
			"aws_access_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("aws_access_key", "defaultValue"),
				Description: "The access key for use with AWS ES Service domains",
			},
			"aws_secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("aws_secret_key", "defaultValue"),
				Description: "The secret key for use with AWS ES Service domains",
			},
			"aws_token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("aws_token", "defaultValue"),
				Description: "The session token for use with AWS ES Service domains",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"esdynamiconfig_index": resourceIndexConfig(),
		},
		ConfigureFunc: providerClient,
	}
}

func providerClient(d *schema.ResourceData) (interface{}, error) {
	config := client.Config{
		Endpoint:        d.Get("es_endpoint").(string),
		Region:          d.Get("region").(string),
		AccessKeyID:     d.Get("aws_access_key").(string),
		SecretAccessKey: d.Get("aws_secret_key").(string),
		SessionToken:    d.Get("aws_token").(string),
	}

	return client.New(config), nil
}
