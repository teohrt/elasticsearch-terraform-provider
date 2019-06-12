package client

import (
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	Endpoint        string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

type Client interface {
	Get(index string) (*http.Response, error)
	Put(index string, body io.Reader) (*http.Response, error)
}

type clientImpl struct {
	endpoint string
	client   *http.Client
}

func New(config Config) Client {
	// creds := awscredentials.NewChainCredentials([]awscredentials.Provider{
	// 	&awscredentials.StaticProvider{
	// 		Value: awscredentials.Value{
	// 			AccessKeyID:     config.AccessKeyID,
	// 			SecretAccessKey: config.SecretAccessKey,
	// 			SessionToken:    config.SessionToken,
	// 		},
	// 	},
	// 	&awscredentials.EnvProvider{},
	// 	&awscredentials.SharedCredentialsProvider{},
	// })
	// signer := awssigv4.NewSigner(creds)
	// signedClient, _ := aws_signing_client.New(signer, nil, "es", config.Region)

	signedClient := &http.Client{} // mocking a correctly signed client for testing

	return clientImpl{
		endpoint: config.Endpoint,
		client:   signedClient,
	}
}

func (c clientImpl) Get(index string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/_settings?pretty", c.endpoint, index)

	res, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c clientImpl) Put(index string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/_settings?pretty", c.endpoint, index)

	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
