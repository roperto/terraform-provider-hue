package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLight() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLightCreate,
		ReadContext:   resourceLightRead,
		UpdateContext: resourceLightUpdate,
		DeleteContext: resourceLightDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"uniqueid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bridge_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"modelid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"manufacturername": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"productname": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"swversion": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLightFindBridgeId(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) error {
	hue := m.(*Hue)

	lights, err := hue.List(ctx, ResourceTypeLight)
	if err != nil {
		return err
	}

	uniqueId := d.Get("uniqueid").(string)
	for bridgeId, light := range lights {
		if light["uniqueid"] == uniqueId {
			d.Set("bridge_id", bridgeId)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Cannot find light by uniqueid: %s", uniqueId))
}

func resourceLightCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	uniqueId := d.Get("uniqueid").(string)
	d.SetId(uniqueId)

	err := resourceLightFindBridgeId(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}

	diags := resourceLightRead(ctx, d, m)

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Light imported (cannot be created) -- run terraform again to configure it.",
	})

	return diags
}

func resourceLightRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId(d.Get("uniqueid").(string))
	bridgeId := d.Get("bridge_id").(string)

	hue := m.(*Hue)
	light, err := hue.Read(ctx, ResourceTypeLight, bridgeId)
	if err != nil {
		switch e := err.(type) {
		case *HueError:
			if e.Type == HueErrorResourceNotAvailable {
				return resourceLightDelete(ctx, d, m)
			}
		}
		return diag.FromErr(err)
	}
	if light["uniqueid"] != d.Id() {
		err := resourceLightFindBridgeId(ctx, d, m)
		if err != nil {
			return diag.FromErr(err)
		}
		return resourceLightRead(ctx, d, m)
	}

	d.Set("address", fmt.Sprintf("/%s/%s", ResourceTypeLight, bridgeId))
	d.Set("name", light["name"])
	d.Set("type", light["type"])
	d.Set("modelid", light["modelid"])
	d.Set("manufacturername", light["manufacturername"])
	d.Set("productname", light["productname"])
	d.Set("swversion", light["swversion"])
	return nil
}

func resourceLightUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	changes := map[string]interface{}{}
	if d.HasChange("name") {
		changes["name"] = d.Get("name").(string)
	}

	hue := m.(*Hue)
	err := hue.Update(ctx, ResourceTypeLight, d.Get("bridge_id").(string), changes)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLightRead(ctx, d, m)
}

func resourceLightDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Light removed from state but not deleted.",
		},
	}
}
