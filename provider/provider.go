package provider

import (
	awscredentials "github.com/aws/aws-sdk-go/aws/credentials"
	awssigv4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/deoxxa/aws_signing_client"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
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
				Optional:    true,
				Default:     schema.EnvDefaultFunc("aws_access_key", nil),
				Description: "The access key for use with AWS ES Service domains",
			},
			"aws_secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     schema.EnvDefaultFunc("aws_secret_key", nil),
				Description: "The secret key for use with AWS ES Service domains",
			},
			"aws_token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     schema.EnvDefaultFunc("aws_token", nil),
				Description: "The session token for use with AWS ES Service domains",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"esdynamiconfig_index": resourceElasticsearchDynamicIndexConfig(),
		},
		ConfigureFunc: providerClient,
	}
}

// returns a client of type *http.Client
func providerClient(d *schema.ResourceData) (interface{}, error) {
	creds := awscredentials.NewChainCredentials([]awscredentials.Provider{
		&awscredentials.StaticProvider{
			Value: awscredentials.Value{
				AccessKeyID:     d.Get("aws_access_key").(string),
				SecretAccessKey: d.Get("aws_secret_key").(string),
				SessionToken:    d.Get("aws_token").(string),
			},
		},
		&awscredentials.EnvProvider{},
		&awscredentials.SharedCredentialsProvider{},
	})
	signer := awssigv4.NewSigner(creds)
	client, _ := aws_signing_client.New(signer, nil, "es", d.Get("region").(string))

	return client, nil
}
