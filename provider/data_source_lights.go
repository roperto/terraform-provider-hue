package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLights() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLightsRead,
		Schema: map[string]*schema.Schema{
			"lights": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uniqueid": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"bridge_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceLightsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hue := m.(*Hue)

	got, err := hue.Read(ctx, ResourceTypeLight, "")
	if err != nil {
		return diag.FromErr(err)
	}

	lights := make([]map[string]interface{}, len(got))
	i := 0
	for bridgeId, light := range got {
		g := light.(map[string]interface{})
		lights[i] = map[string]interface{}{
			"bridge_id": bridgeId,
			"uniqueid":  g["uniqueid"],
			"name":      g["name"],
			"type":      g["type"],
		}
		i++
	}
	d.SetId("lights")
	d.Set("lights", lights)

	return nil
}
