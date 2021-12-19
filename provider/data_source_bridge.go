package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBridge() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBridgeRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"datastoreversion": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"swversion": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"apiversion": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"mac": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"factorynew": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"replacesbridgeid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"modelid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"starterkitid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBridgeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hue := m.(*Hue)
	config, err := hue.Read(ctx, "config", "")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(config["bridgeid"].(string))
	d.Set("name", config["name"])
	d.Set("datastoreversion", config["datastoreversion"])
	d.Set("swversion", config["swversion"])
	d.Set("apiversion", config["apiversion"])
	d.Set("mac", config["mac"])
	d.Set("factorynew", config["factorynew"])
	d.Set("replacesbridgeid", config["replacesbridgeid"])
	d.Set("modelid", config["modelid"])
	d.Set("starterkitid", config["starterkitid"])

	return nil
}
