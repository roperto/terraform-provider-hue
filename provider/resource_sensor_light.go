package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSensorLight() *schema.Resource {
	s := resourceSensor()
	s.CreateContext = resourceSensorLightCreate
	s.ReadContext = resourceSensorLightRead
	s.UpdateContext = resourceSensorLightUpdate
	s.DeleteContext = resourceSensorLightDelete
	s.Schema["config_on"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	}
	s.Schema["config_battery"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
	s.Schema["config_reachable"] = &schema.Schema{
		Type:     schema.TypeBool,
		Computed: true,
	}
	s.Schema["config_alert"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	s.Schema["config_tholddark"] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  16000,
	}
	s.Schema["config_tholdoffset"] = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  7000,
	}
	s.Schema["config_ledindication"] = &schema.Schema{
		Type:     schema.TypeBool,
		Computed: true,
	}
	s.Schema["config_usertest"] = &schema.Schema{
		Type:     schema.TypeBool,
		Computed: true,
	}
	return s
}

func resourceSensorLightCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	diags := resourceSensorCreate(ctx, d, m)
	if diags.HasError() {
		return diags
	}

	moreDiags := resourceSensorLightUpdate(ctx, d, m)
	for _, v := range moreDiags {
		diags = append(diags, v)
	}

	return diags
}

func resourceSensorLightRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sensor, diags := resourceSensorReadData(ctx, d, m)
	if diags.HasError() {
		return diags
	}
	if sensor == nil {
		return diags
	}

	config := sensor["config"].(map[string]interface{})
	d.Set("config_on", config["on"])
	d.Set("config_battery", config["battery"])
	d.Set("config_reachable", config["reachable"])
	d.Set("config_alert", config["alert"])
	d.Set("config_tholddark", config["tholddark"])
	d.Set("config_tholdoffset", config["tholdoffset"])
	d.Set("config_ledindication", config["ledindication"])
	d.Set("config_usertest", config["usertest"])

	return diags
}

func resourceSensorLightUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := map[string]interface{}{}
	if d.HasChange("config_on") {
		config["on"] = d.Get("config_on")
	}
	if d.HasChange("config_tholddark") {
		config["tholddark"] = d.Get("config_tholddark")
	}
	if d.HasChange("config_tholdoffset") {
		config["tholdoffset"] = d.Get("config_tholdoffset")
	}

	changed := map[string]interface{}{}
	if len(config) > 0 {
		changed["config"] = config
	}

	diags := resourceSensorUpdateData(ctx, d, m, changed)
	return diags
}

func resourceSensorLightDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceSensorDelete(ctx, d, m)
}
