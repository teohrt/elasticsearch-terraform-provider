package client

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

type Config struct {
	Endpoint string
	Region   string
	Creds    *credentials.Credentials
}

type Client interface {
	Get(index string) (*GetItemResponse, error)
	Put(index string, body io.ReadSeeker) (*http.Response, error)
}

type clientImpl struct {
	Config
	client *http.Client
}

func New(config *Config) Client {
	return clientImpl{
		Config: *config,
		client: &http.Client{},
	}
}

func (c clientImpl) Get(index string) (*GetItemResponse, error) {
	url := fmt.Sprintf("%s/%s/_settings?pretty", c.Endpoint, index)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	signedReq, err := c.signRequest(req, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(signedReq)
	if err != nil {
		return nil, err
	}

	return getItem(res, index)
}

func (c clientImpl) Put(index string, body io.ReadSeeker) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/_settings", c.Endpoint, index)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	signedReq, err := c.signRequest(req, body)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(signedReq)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func getItem(res *http.Response, indexName string) (*GetItemResponse, error) {
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &GetItemResponse{}, err
	}

	var jsonInterface interface{}
	if err = json.Unmarshal(resBody, &jsonInterface); err != nil {
		return &GetItemResponse{}, err
	}

	// TODO: validation required
	qwt := jsonInterface.(map[string]interface{})[indexName].(map[string]interface{})["settings"].(map[string]interface{})["index"].(map[string]interface{})["search"].(map[string]interface{})["slowlog"].(map[string]interface{})["threshold"].(map[string]interface{})["query"].(map[string]interface{})["warn"].(string)
	qit := jsonInterface.(map[string]interface{})[indexName].(map[string]interface{})["settings"].(map[string]interface{})["index"].(map[string]interface{})["search"].(map[string]interface{})["slowlog"].(map[string]interface{})["threshold"].(map[string]interface{})["query"].(map[string]interface{})["info"].(string)

	return &GetItemResponse{
		Query_warn_threshold: qwt,
		Query_info_threshold: qit,
	}, nil
}

func (c clientImpl) signRequest(req *http.Request, body io.ReadSeeker) (*http.Request, error) {
	signer := v4.NewSigner(c.Creds)

	_, err := signer.Sign(req, body, "es", c.Region, time.Now())
	if err != nil {
		return nil, err
	}

	return req, nil
}
