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
			"elasticsearch_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"elasticsearch_dynamic_index_config": resourceItem(),
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
