package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRulesRead,
		Schema: map[string]*schema.Schema{
			"rules": {
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

func dataSourceRulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hue := m.(*Hue)

	got, err := hue.Read(ctx, ResourceTypeRule, "")
	if err != nil {
		return diag.FromErr(err)
	}

	rules := make([]map[string]interface{}, len(got))
	i := 0
	for bridgeId, rule := range got {
		r := rule.(map[string]interface{})

		rules[i] = map[string]interface{}{
			"bridge_id":  bridgeId,
			"name":       r["name"],
		}
		i++
	}
	d.SetId("rules")
	err = d.Set("rules", rules)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
