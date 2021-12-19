package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceScenes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceScenesRead,
		Schema: map[string]*schema.Schema{
			"scenes": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scene_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"scene_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceScenesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hue := m.(*Hue)

	groups, err := hue.Read(ctx, ResourceTypeGroup, "")
	if err != nil {
		return diag.FromErr(err)
	}

	got, err := hue.Read(ctx, ResourceTypeScene, "")
	if err != nil {
		return diag.FromErr(err)
	}

	scenes := make([]map[string]interface{}, len(got))
	i := 0
	for bridgeId, scene := range got {
		s := scene.(map[string]interface{})
		group_id := s["group"].(string)
		group_data := groups[group_id].(map[string]interface{})
		scenes[i] = map[string]interface{}{
			"scene_id":   bridgeId,
			"scene_name": s["name"],
			"group_id":   group_id,
			"group_name": group_data["name"],
		}
		i++
	}
	d.SetId("scenes")
	d.Set("scenes", scenes)

	return nil
}
