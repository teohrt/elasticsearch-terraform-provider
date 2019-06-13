package provider

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/teohrt/terraform-provider-esslowlogconfig/client"
)

func resourceIndexConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the resource, also acts as it's unique ID",
				ForceNew:     true,
				ValidateFunc: validateIndex,
			},
			"query_warn_threshold": {
				Type:         schema.TypeString,
				Default:      -1,
				Optional:     true,
				ForceNew:     false,
				Description:  "",
				ValidateFunc: validateIndex,
			},
			"query_info_threshold": {
				Type:         schema.TypeString,
				Default:      -1,
				Optional:     true,
				ForceNew:     false,
				Description:  "",
				ValidateFunc: validateIndex,
			},
			"query_debug_threshold": {
				Type:         schema.TypeString,
				Default:      -1,
				Optional:     true,
				ForceNew:     false,
				Description:  "",
				ValidateFunc: validateIndex,
			},
			"query_trace_threshold": {
				Type:         schema.TypeString,
				Default:      -1,
				Optional:     true,
				ForceNew:     false,
				Description:  "",
				ValidateFunc: validateIndex,
			},
			"fetch_warn_threshold": {
				Type:         schema.TypeString,
				Default:      -1,
				Optional:     true,
				ForceNew:     false,
				Description:  "",
				ValidateFunc: validateIndex,
			},
			"fetch_info_threshold": {
				Type:         schema.TypeString,
				Default:      -1,
				Optional:     true,
				ForceNew:     false,
				Description:  "",
				ValidateFunc: validateIndex,
			},
			"fetch_debug_threshold": {
				Type:         schema.TypeString,
				Default:      -1,
				Optional:     true,
				ForceNew:     false,
				Description:  "",
				ValidateFunc: validateIndex,
			},
			"fetch_trace_threshold": {
				Type:         schema.TypeString,
				Default:      -1,
				Optional:     true,
				ForceNew:     false,
				Description:  "",
				ValidateFunc: validateIndex,
			},
		},
		Create: resourceCreateItem,
		Read:   resourceReadItem,
		Update: resourceCreateItem,
		Delete: resourceDeleteItem,
	}
}

func validateIndex(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("Expected index to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("index cannot contain whitespace. Got %s", value))
		return warns, errs
	}
	return warns, errs
}

func resourceCreateItem(d *schema.ResourceData, m interface{}) error {
	client := m.(client.Client)
	reqBody := strings.NewReader(fmt.Sprintf(`
		{
			"index.search.slowlog.threshold.query.warn": "%s",
			"index.search.slowlog.threshold.query.info": "%s",
			"index.search.slowlog.threshold.query.debug": "%s",
			"index.search.slowlog.threshold.query.trace": "%s",

			"index.search.slowlog.threshold.fetch.warn": "%s",
			"index.search.slowlog.threshold.fetch.info": "%s",
			"index.search.slowlog.threshold.fetch.debug": "%s",
			"index.search.slowlog.threshold.fetch.trace": "%s"
		}`,
		d.Get("query_warn_threshold").(string),
		d.Get("query_info_threshold").(string),
		d.Get("query_debug_threshold").(string),
		d.Get("query_trace_threshold").(string),

		d.Get("fetch_warn_threshold").(string),
		d.Get("fetch_info_threshold").(string),
		d.Get("fetch_debug_threshold").(string),
		d.Get("fetch_trace_threshold").(string),
	))

	res, err := client.Put(d.Get("name").(string), reqBody)
	if err != nil {
		return err
	}

	err = handleResponse(res)
	if err != nil {
		return errors.New(fmt.Sprintf("Bad response in resourceCreateItem: %s", err.Error()))
	}

	d.SetId(d.Get("name").(string))
	return nil
}

func resourceReadItem(d *schema.ResourceData, m interface{}) error {
	client := m.(client.Client)

	item, err := client.Get(d.Get("name").(string))
	if err != nil {
		return err
	}

	d.Set("name", d.Id())
	d.Set("query_warn_threshold", item.Query_warn_threshold)
	d.Set("query_info_threshold", item.Query_info_threshold)
	d.Set("query_debug_threshold", item.Query_debug_threshold)
	d.Set("query_trace_threshold", item.Query_trace_threshold)

	d.Set("fetch_warn_threshold", item.Fetch_warn_threshold)
	d.Set("fetch_info_threshold", item.Fetch_info_threshold)
	d.Set("fetch_debug_threshold", item.Fetch_debug_threshold)
	d.Set("fetch_trace_threshold", item.Fetch_trace_threshold)

	return nil
}

func resourceDeleteItem(d *schema.ResourceData, m interface{}) error {
	client := m.(client.Client)
	body := strings.NewReader(`
		{
			"index.search.slowlog.threshold.query.warn": -1,
			"index.search.slowlog.threshold.query.info": -1,
			"index.search.slowlog.threshold.query.debug": -1,
			"index.search.slowlog.threshold.query.trace": -1,

			"index.search.slowlog.threshold.fetch.warn": -1,
			"index.search.slowlog.threshold.fetch.info": -1,
			"index.search.slowlog.threshold.fetch.debug": -1,
			"index.search.slowlog.threshold.fetch.trace": -1
		}`)

	res, err := client.Put(d.Get("name").(string), body)
	if err != nil {
		return err
	}

	err = handleResponse(res)
	if err != nil {
		return errors.New(fmt.Sprintf("Bad response in resourceDeleteItem: %s", err.Error()))
	}

	return err
}

func handleResponse(res *http.Response) error {
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if string(resBody) != `{"acknowledged":true}` {
		return errors.New(string(resBody))
	}

	return nil
}
