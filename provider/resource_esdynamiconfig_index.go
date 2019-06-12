package provider

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceElasticsearchDynamicIndexConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"indexName": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the resource, also acts as it's unique ID",
				ForceNew:    true,
				// ValidateFunc: validateIndex,
			},
			"query_warn_threshold": {
				Type:         schema.TypeString,
				Default:      "-1",
				Optional:     true,
				ForceNew:     false,
				Description:  "",
				ValidateFunc: validateIndex,
			},
			"query_info_threshold": {
				Type:         schema.TypeString,
				Default:      "-1",
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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
	client := m.(*http.Client)
	url := fmt.Sprintf("%s/%s/_settings?pretty", d.Get("es_endpoint").(string), d.Get("indexName"))

	body := strings.NewReader(fmt.Sprintf(`
		{
			"index.search.slowlog.threshold.query.warn": "%s",
			"index.search.slowlog.threshold.query.info": "%s"
		}`,
		d.Get("query_warn_threshold").(string),
		d.Get("query_info_threshold").(string)))

	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if res != nil {
		return err
	}

	d.SetId(d.Get("indexName").(string))
	return nil
}

func resourceReadItem(d *schema.ResourceData, m interface{}) error {
	client := m.(*http.Client)
	url := fmt.Sprintf("%s/%s", d.Get("es_endpoint").(string), d.Get("indexName"))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	d.Set("indexName", d.Id())
	d.Set("query_warn_threshold", "TODO")
	d.Set("query_info_threshold", "TODO")

	return nil
}

func resourceDeleteItem(d *schema.ResourceData, m interface{}) error {
	client := m.(*http.Client)
	url := fmt.Sprintf("%s/%s", d.Get("es_endpoint").(string), d.Get("indexName"))

	body := strings.NewReader(`
		{
			"index.search.slowlog.threshold.query.warn": -1,
			"index.search.slowlog.threshold.query.info": -1
		}`)

	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return err
	}

	_, err = client.Do(req)

	return err
}

func resourceExistsItem(d *schema.ResourceData, m interface{}) error {
	client := m.(*http.Client)
	url := fmt.Sprintf("%s/%s", d.Get("es_endpoint").(string), d.Get("indexName"))

	body := strings.NewReader(`
		{
			"index.search.slowlog.threshold.query.warn": -1,
			"index.search.slowlog.threshold.query.info": -1
		}`)

	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return err
	}

	_, err = client.Do(req)

	return err
}
