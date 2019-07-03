package client

import (
	"encoding/json"
	"errors"
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

	err = c.awsSigner(req, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
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
	err = c.awsSigner(req, body)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func getItem(res *http.Response, indexName string) (*GetItemResponse, error) {
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &GetItemResponse{}, err
	}

	item := &GetItemResponse{}

	var jsonInterface interface{}
	if err = json.Unmarshal(resBody, &jsonInterface); err != nil {
		return &GetItemResponse{}, err
	}

	if qwt, ok := jsonInterface.(map[string]interface{})[indexName].(map[string]interface{})["settings"].(map[string]interface{})["index"].(map[string]interface{})["search"].(map[string]interface{})["slowlog"].(map[string]interface{})["threshold"].(map[string]interface{})["query"].(map[string]interface{})["warn"].(string); ok {
		item.Query_warn_threshold = qwt
	} else {
		apiError()
	}
	if qit, ok := jsonInterface.(map[string]interface{})[indexName].(map[string]interface{})["settings"].(map[string]interface{})["index"].(map[string]interface{})["search"].(map[string]interface{})["slowlog"].(map[string]interface{})["threshold"].(map[string]interface{})["query"].(map[string]interface{})["info"].(string); ok {
		item.Query_info_threshold = qit
	} else {
		apiError()
	}
	if qdt, ok := jsonInterface.(map[string]interface{})[indexName].(map[string]interface{})["settings"].(map[string]interface{})["index"].(map[string]interface{})["search"].(map[string]interface{})["slowlog"].(map[string]interface{})["threshold"].(map[string]interface{})["query"].(map[string]interface{})["debug"].(string); ok {
		item.Query_debug_threshold = qdt
	} else {
		apiError()
	}
	if qtt, ok := jsonInterface.(map[string]interface{})[indexName].(map[string]interface{})["settings"].(map[string]interface{})["index"].(map[string]interface{})["search"].(map[string]interface{})["slowlog"].(map[string]interface{})["threshold"].(map[string]interface{})["query"].(map[string]interface{})["trace"].(string); ok {
		item.Query_trace_threshold = qtt
	} else {
		apiError()
	}

	if fwt, ok := jsonInterface.(map[string]interface{})[indexName].(map[string]interface{})["settings"].(map[string]interface{})["index"].(map[string]interface{})["search"].(map[string]interface{})["slowlog"].(map[string]interface{})["threshold"].(map[string]interface{})["fetch"].(map[string]interface{})["warn"].(string); ok {
		item.Fetch_warn_threshold = fwt
	} else {
		apiError()
	}
	if fit, ok := jsonInterface.(map[string]interface{})[indexName].(map[string]interface{})["settings"].(map[string]interface{})["index"].(map[string]interface{})["search"].(map[string]interface{})["slowlog"].(map[string]interface{})["threshold"].(map[string]interface{})["fetch"].(map[string]interface{})["info"].(string); ok {
		item.Fetch_info_threshold = fit
	} else {
		apiError()
	}
	if fdt, ok := jsonInterface.(map[string]interface{})[indexName].(map[string]interface{})["settings"].(map[string]interface{})["index"].(map[string]interface{})["search"].(map[string]interface{})["slowlog"].(map[string]interface{})["threshold"].(map[string]interface{})["fetch"].(map[string]interface{})["debug"].(string); ok {
		item.Fetch_debug_threshold = fdt
	} else {
		apiError()
	}
	if ftt, ok := jsonInterface.(map[string]interface{})[indexName].(map[string]interface{})["settings"].(map[string]interface{})["index"].(map[string]interface{})["search"].(map[string]interface{})["slowlog"].(map[string]interface{})["threshold"].(map[string]interface{})["fetch"].(map[string]interface{})["trace"].(string); ok {
		item.Fetch_trace_threshold = ftt
	} else {
		apiError()
	}

	return item, nil
}

func (c clientImpl) awsSigner(req *http.Request, body io.ReadSeeker) error {
	signer := v4.NewSigner(c.Creds)
	_, err := signer.Sign(req, body, "es", c.Region, time.Now())
	return err
}

func apiError() (*GetItemResponse, error) {
	return nil, errors.New("Provider error: An outdated version of Elasticsearch's slow log API contract is being used.")
}
