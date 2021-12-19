package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSensors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSensorsRead,
		Schema: map[string]*schema.Schema{
			"sensors": {
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

func dataSourceSensorsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hue := m.(*Hue)

	got, err := hue.Read(ctx, ResourceTypeSensor, "")
	if err != nil {
		return diag.FromErr(err)
	}

	sensors := make([]map[string]interface{}, len(got))
	i := 0
	for bridgeId, sensor := range got {
		s := sensor.(map[string]interface{})
		sensors[i] = map[string]interface{}{
			"uniqueid":  s["uniqueid"],
			"bridge_id": bridgeId,
			"name":      s["name"],
			"type":      s["type"],
		}
		i++
	}
	d.SetId("sensors")
	d.Set("sensors", sensors)

	return nil
}
