package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupsRead,
		Schema: map[string]*schema.Schema{
			"groups": {
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
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"lights": &schema.Schema{
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

func dataSourceGroupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hue := m.(*Hue)

	got, err := hue.Read(ctx, ResourceTypeGroup, "")
	if err != nil {
		return diag.FromErr(err)
	}

	groups := make([]map[string]interface{}, len(got))
	i := 0
	for bridgeId, group := range got {
		g := group.(map[string]interface{})
		groups[i] = map[string]interface{}{
			"bridge_id": bridgeId,
			"name":      g["name"],
			"type":      g["type"],
			"lights":    g["lights"],
		}
		i++
	}
	d.SetId("groups")
	d.Set("groups", groups)

	return nil
}
