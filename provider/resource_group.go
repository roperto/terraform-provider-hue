package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: resourceGroupValidateType,
			},
			"lights": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceGroupValidateType(i interface{}, path cty.Path) diag.Diagnostics {
	allowed := []string{"LightGroup", "Room", "Luminaire", "LightSource"}
	for _, v := range allowed {
		if v == i {
			return nil
		}
	}
	return diag.Errorf(
		"Invalid type '%s'. Allowed values: %s",
		i,
		allowed,
	)
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	lights_got := d.Get("lights").(*schema.Set).List()
	lights := make([]string, len(lights_got))
	for k, v := range lights_got {
		lights[k] = v.(string)
	}

	data := map[string]interface{}{
		"name":   d.Get("name").(string),
		"type":   d.Get("type").(string),
		"lights": lights,
	}

	hue := m.(*Hue)
	id, err := hue.Create(ctx, ResourceTypeGroup, data)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return resourceGroupRead(ctx, d, m)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hue := m.(*Hue)

	group, err := hue.Read(ctx, ResourceTypeGroup, d.Id())
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

	d.Set("address", fmt.Sprintf("/%s/%s", ResourceTypeGroup, d.Id()))
	d.Set("name", group["name"])
	d.Set("type", group["type"])
	d.Set("lights", group["lights"])

	return nil
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	changes := map[string]interface{}{}
	if d.HasChange("name") {
		changes["name"] = d.Get("name").(string)
	}
	if d.HasChange("lights") {
		lights_got := d.Get("lights").(*schema.Set).List()
		lights := make([]string, len(lights_got))
		for k, v := range lights_got {
			lights[k] = v.(string)
		}
		changes["lights"] = lights
	}

	hue := m.(*Hue)
	err := hue.Update(ctx, ResourceTypeGroup, d.Id(), changes)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGroupRead(ctx, d, m)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hue := m.(*Hue)

	err := hue.Delete(ctx, ResourceTypeGroup, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
