package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CannotFindSensorError struct {
	UniqueID string
}

func (e *CannotFindSensorError) Error() string {
	return fmt.Sprintf("Cannot find sensor by uniqueid: %s", e.UniqueID)
}

func resourceSensor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSensorCreate,
		ReadContext:   resourceSensorRead,
		UpdateContext: resourceSensorUpdate,
		DeleteContext: resourceSensorDelete,
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
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSensorFindBridgeId(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) error {
	hue := m.(*Hue)

	sensors, err := hue.List(ctx, ResourceTypeSensor)
	if err != nil {
		return err
	}

	for bridgeId, sensor := range sensors {
		if sensor["uniqueid"] == d.Id() {
			d.Set("bridge_id", bridgeId)
			return nil
		}
	}
	return &CannotFindSensorError{UniqueID: d.Id()}
}

func resourceSensorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	uniqueId := d.Get("uniqueid").(string)
	d.SetId(uniqueId)

	if d.Get("type") != "CLIPGenericStatus" {
		err := resourceSensorFindBridgeId(ctx, d, m)
		if err != nil {
			return diag.FromErr(err)
		}

		diags := resourceSensorRead(ctx, d, m)

		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Sensor imported (cannot be created) -- run terraform again to configure it.",
		})

		return diags
	}

	data := map[string]interface{}{
		"uniqueid":         d.Get("uniqueid").(string),
		"name":             d.Get("name").(string),
		"type":             d.Get("type").(string),
		"modelid":          "Custom State Sensor",
		"swversion":        "1.0",
		"manufacturername": "Terraform Hue",
	}

	hue := m.(*Hue)
	id, err := hue.Create(ctx, ResourceTypeSensor, data)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("brigde_id", id)
	return resourceSensorRead(ctx, d, m)
}

func resourceSensorReadData(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) (map[string]interface{}, diag.Diagnostics) {
	if d.Get("uniqueid").(string) == "" {
		d.Set("uniqueid", d.Id())
	} else {
		d.SetId(d.Get("uniqueid").(string))
	}

	hue := m.(*Hue)

	bridgeId := d.Get("bridge_id").(string)
	if bridgeId != "" {
		sensor, err := hue.Read(ctx, ResourceTypeSensor, bridgeId)
		if err != nil {
			switch e := err.(type) {
			case *HueError:
				if e.Type == HueErrorResourceNotAvailable {
					d.SetId("")
					return nil, nil
				}
			}
			return nil, diag.FromErr(err)
		}

		if sensor["uniqueid"] == d.Id() {
			d.Set("name", sensor["name"])
			d.Set("type", sensor["type"])
			d.Set("address", fmt.Sprintf("/%s/%s", ResourceTypeSensor, bridgeId))
			return sensor, nil
		}
	}

	err := resourceSensorFindBridgeId(ctx, d, m)
	if err != nil {
		switch err.(type) {
		case *CannotFindSensorError:
			d.SetId("")
			return nil, nil
		}
		return nil, diag.FromErr(err)
	}
	return resourceSensorReadData(ctx, d, m)
}

func resourceSensorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, diags := resourceSensorReadData(ctx, d, m)
	return diags
}

func resourceSensorUpdateData(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
	changes map[string]interface{},
) diag.Diagnostics {
	if d.HasChange("name") {
		changes["name"] = d.Get("name").(string)
	}

	hue := m.(*Hue)
	err := hue.Update(ctx, ResourceTypeSensor, d.Get("bridge_id").(string), changes)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSensorRead(ctx, d, m)
}
func resourceSensorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceSensorUpdateData(ctx, d, m, map[string]interface{}{})
}

func resourceSensorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Get("type") != "CLIPGenericStatus" {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Sensor removed from state but not deleted.",
			},
		}
	}

	hue := m.(*Hue)
	err := hue.Delete(ctx, ResourceTypeSensor, d.Get("bridge_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
