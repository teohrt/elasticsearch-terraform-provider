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

	res, err := c.client.Do(req)
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
	qdt := jsonInterface.(map[string]interface{})[indexName].(map[string]interface{})["settings"].(map[string]interface{})["index"].(map[string]interface{})["search"].(map[string]interface{})["slowlog"].(map[string]interface{})["threshold"].(map[string]interface{})["query"].(map[string]interface{})["debug"].(string)
	qtt := jsonInterface.(map[string]interface{})[indexName].(map[string]interface{})["settings"].(map[string]interface{})["index"].(map[string]interface{})["search"].(map[string]interface{})["slowlog"].(map[string]interface{})["threshold"].(map[string]interface{})["query"].(map[string]interface{})["trace"].(string)

	fwt := jsonInterface.(map[string]interface{})[indexName].(map[string]interface{})["settings"].(map[string]interface{})["index"].(map[string]interface{})["search"].(map[string]interface{})["slowlog"].(map[string]interface{})["threshold"].(map[string]interface{})["fetch"].(map[string]interface{})["warn"].(string)
	fit := jsonInterface.(map[string]interface{})[indexName].(map[string]interface{})["settings"].(map[string]interface{})["index"].(map[string]interface{})["search"].(map[string]interface{})["slowlog"].(map[string]interface{})["threshold"].(map[string]interface{})["fetch"].(map[string]interface{})["info"].(string)
	fdt := jsonInterface.(map[string]interface{})[indexName].(map[string]interface{})["settings"].(map[string]interface{})["index"].(map[string]interface{})["search"].(map[string]interface{})["slowlog"].(map[string]interface{})["threshold"].(map[string]interface{})["fetch"].(map[string]interface{})["debug"].(string)
	ftt := jsonInterface.(map[string]interface{})[indexName].(map[string]interface{})["settings"].(map[string]interface{})["index"].(map[string]interface{})["search"].(map[string]interface{})["slowlog"].(map[string]interface{})["threshold"].(map[string]interface{})["fetch"].(map[string]interface{})["trace"].(string)

	return &GetItemResponse{
		Query_warn_threshold:  qwt,
		Query_info_threshold:  qit,
		Query_debug_threshold: qdt,
		Query_trace_threshold: qtt,

		Fetch_warn_threshold:  fwt,
		Fetch_info_threshold:  fit,
		Fetch_debug_threshold: fdt,
		Fetch_trace_threshold: ftt,
	}, nil
}

func (c clientImpl) awsSigner(req *http.Request, body io.ReadSeeker) error {
	signer := v4.NewSigner(c.Creds)
	_, err := signer.Sign(req, body, "es", c.Region, time.Now())
	return err
}
