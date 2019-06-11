package provider

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceItem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"indexName": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the resource, also acts as it's unique ID",
				ForceNew:     true,
				ValidateFunc: validateIndex,
			},
			"query_warn": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the resource, also acts as it's unique ID",
				ForceNew:     true,
				ValidateFunc: validateIndex,
			},
			"query_info": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the resource, also acts as it's unique ID",
				ForceNew:     true,
				ValidateFunc: validateIndex,
			},
		},
		Create: resourceCreateItem,
		// Read:   resourceReadItem,
		// Update: resourceUpdateItem,
		// Delete: resourceDeleteItem,
		// Exists: resourceExistsItem,
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
	url := d.Get("elasticsearch_endpoint").(string)

	body := strings.NewReader(fmt.Sprintf(`
		{
			"index.search.slowlog.threshold.query.warn": "%s",
			"index.search.slowlog.threshold.query.info": "%s"
		}`,
		d.Get("query_warn").(string),
		d.Get("query_info").(string)))

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
