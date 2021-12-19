package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSchedules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSchedulesRead,
		Schema: map[string]*schema.Schema{
			"schedules": {
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
					},
				},
			},
		},
	}
}

func dataSourceSchedulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hue := m.(*Hue)

	got, err := hue.Read(ctx, ResourceTypeSchedule, "")
	if err != nil {
		return diag.FromErr(err)
	}

	schedules := make([]map[string]interface{}, len(got))
	i := 0
	for bridgeId, schedule := range got {
		s := schedule.(map[string]interface{})
		schedules[i] = map[string]interface{}{
			"bridge_id": bridgeId,
			"name":      s["name"],
		}
		i++
	}
	d.SetId("schedules")
	d.Set("schedules", schedules)

	return nil
}
