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

	return c.client.Do(req)
}

func getItem(res *http.Response, indexName string) (*GetItemResponse, error) {
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &GetItemResponse{}, err
	}

	/*
		All of this marshalling and unmarshalling is necessary to
		circumvent a dynamic field name in the response body. The
		field depends on the index name. Because of this, unmarshalling
		directly into a generic struct doesn't work because not all of
		the fields are static.

		What this does is basically dive into the response body JSON
		just past the indexName field, so we're able to unmarshal the
		rest of the JSON into a struct that will match our field names
		accurately.
	*/
	var jsonInterface interface{}
	if err = json.Unmarshal(resBody, &jsonInterface); err != nil {
		return nil, err
	}

	unmarshalledIndex, ok := jsonInterface.(map[string]interface{})[indexName]
	if !ok {
		return nil, err
	}

	bytes, err := json.Marshal(unmarshalledIndex)
	if err != nil {
		return nil, err
	}

	var container Container
	if err = json.Unmarshal(bytes, &container); err != nil {
		return nil, err
	}

	return &GetItemResponse{
		Query_warn_threshold:  container.Settings.Index.Search.Slowlog.Threshold.Query.Warn,
		Query_info_threshold:  container.Settings.Index.Search.Slowlog.Threshold.Query.Info,
		Query_debug_threshold: container.Settings.Index.Search.Slowlog.Threshold.Query.Debug,
		Query_trace_threshold: container.Settings.Index.Search.Slowlog.Threshold.Query.Trace,

		Fetch_warn_threshold:  container.Settings.Index.Search.Slowlog.Threshold.Fetch.Warn,
		Fetch_info_threshold:  container.Settings.Index.Search.Slowlog.Threshold.Fetch.Info,
		Fetch_debug_threshold: container.Settings.Index.Search.Slowlog.Threshold.Fetch.Debug,
		Fetch_trace_threshold: container.Settings.Index.Search.Slowlog.Threshold.Fetch.Trace,
	}, nil
}

func (c clientImpl) awsSigner(req *http.Request, body io.ReadSeeker) error {
	signer := v4.NewSigner(c.Creds)
	_, err := signer.Sign(req, body, "es", c.Region, time.Now())
	return err
}
