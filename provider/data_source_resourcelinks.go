package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceResourcelinks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourcelinksRead,
		Schema: map[string]*schema.Schema{
			"resourcelinks": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bridge_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"links": &schema.Schema{
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceResourcelinksRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hue := m.(*Hue)

	got, err := hue.Read(ctx, ResourceTypeResourcelink, "")
	if err != nil {
		return diag.FromErr(err)
	}

	resourcelinks := make([]map[string]interface{}, len(got))
	i := 0
	for bridgeId, link := range got {
		g := link.(map[string]interface{})
		resourcelinks[i] = map[string]interface{}{
			"bridge_id":   bridgeId,
			"name":        g["name"],
			"description": g["description"],
			"links":       g["links"],
		}
		i++
	}
	d.SetId("resourcelinks")
	d.Set("resourcelinks", resourcelinks)

	return nil
}
