package provider

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/teohrt/terraform-provider-esdynamiconfig/client"
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
	body := strings.NewReader(fmt.Sprintf(`
		{
			"index.search.slowlog.threshold.query.warn": "%s",
			"index.search.slowlog.threshold.query.info": "%s"
		}`,
		d.Get("query_warn_threshold").(string),
		d.Get("query_info_threshold").(string)))

	res, err := client.Put(d.Get("name").(string), body)
	if err != nil {
		return err
	}
	if res == nil { // TODO: validate res
		return errors.New("Bad response from endpoint during CreateItem")
	}

	d.SetId(d.Get("name").(string))
	return nil
}

func resourceReadItem(d *schema.ResourceData, m interface{}) error {
	client := m.(client.Client)

	res, err := client.Get(d.Get("name").(string))
	if err != nil {
		return err
	}
	if res == nil { // TODO: validate res
		return errors.New("Bad response from endpoint during ReadItem")
	}

	d.Set("name", d.Id())
	d.Set("query_warn_threshold", "TODO")
	d.Set("query_info_threshold", "TODO")

	return nil
}

func resourceDeleteItem(d *schema.ResourceData, m interface{}) error {
	client := m.(client.Client)
	body := strings.NewReader(`
		{
			"index.search.slowlog.threshold.query.warn": -1,
			"index.search.slowlog.threshold.query.info": -1
		}`)

	res, err := client.Put(d.Get("name").(string), body)
	if err != nil {
		return err
	}
	if res == nil { // TODO: validate res
		return errors.New("Bad response from endpoint during CreateItem")
	}

	return err
}
