package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRuleCreate,
		ReadContext:   resourceRuleRead,
		UpdateContext: resourceRuleUpdate,
		DeleteContext: resourceRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "enabled",
			},
			"conditions": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"actions": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error

	var conditions []interface{}
	err = json.Unmarshal([]byte(d.Get("conditions").(string)), &conditions)
	if err != nil {
		return diag.FromErr(err)
	}

	var actions []interface{}
	err = json.Unmarshal([]byte(d.Get("actions").(string)), &actions)
	if err != nil {
		return diag.FromErr(err)
	}

	data := map[string]interface{}{
		"name":       d.Get("name"),
		"conditions": conditions,
		"actions":    actions,
		"status":    d.Get("status"),
	}

	hue := m.(*Hue)
	id, err := hue.Create(ctx, ResourceTypeRule, data)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return resourceRuleRead(ctx, d, m)
}

func resourceRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hue := m.(*Hue)

	rule, err := hue.Read(ctx, ResourceTypeRule, d.Id())
	if err != nil {
		switch e := err.(type) {
		case *HueError:
			if e.Type == HueErrorResourceNotAvailable {
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	conditions, err := json.Marshal(rule["conditions"])
	if err != nil {
		return diag.FromErr(err)
	}

	actions, err := json.Marshal(rule["actions"])
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", rule["name"])
	d.Set("status", rule["status"])
	d.Set("conditions", string(conditions))
	d.Set("actions", string(actions))
	d.Set("address", fmt.Sprintf("/%s/%s", ResourceTypeRule, d.Id()))

	return nil
}

func resourceRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	changes := map[string]interface{}{}
	if d.HasChange("name") {
		changes["name"] = d.Get("name")
	}
	if d.HasChange("enabled") {
		changes["enabled"] = d.Get("enabled")
	}
	if d.HasChange("conditions") {
		var conditions []interface{}
		err := json.Unmarshal([]byte(d.Get("conditions").(string)), &conditions)
		if err != nil {
			return diag.FromErr(err)
		}
		changes["conditions"] = conditions
	}
	if d.HasChange("actions") {
		var actions []interface{}
		err := json.Unmarshal([]byte(d.Get("actions").(string)), &actions)
		if err != nil {
			return diag.FromErr(err)
		}
		changes["actions"] = actions
	}

	hue := m.(*Hue)
	err := hue.Update(ctx, ResourceTypeRule, d.Id(), changes)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRuleRead(ctx, d, m)
}

func resourceRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hue := m.(*Hue)

	err := hue.Delete(ctx, ResourceTypeRule, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
