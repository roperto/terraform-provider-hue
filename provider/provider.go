package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hostname": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HUE_HOSTNAME", nil),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HUE_USERNAME", nil),
			},
		},
		ConfigureContextFunc: providerConfigure,
		DataSourcesMap: map[string]*schema.Resource{
			"hue_bridge":        dataSourceBridge(),
			"hue_lights":        dataSourceLights(),
			"hue_groups":        dataSourceGroups(),
			"hue_schedules":     dataSourceSchedules(),
			"hue_scenes":        dataSourceScenes(),
			"hue_sensors":       dataSourceSensors(),
			"hue_rules":         dataSourceRules(),
			"hue_resourcelinks": dataSourceResourcelinks(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"hue_group":        resourceGroup(),
			"hue_light":        resourceLight(),
			"hue_scene":        resourceScene(),
			"hue_sensor":       resourceSensor(),
			"hue_sensor_light": resourceSensorLight(),
			"hue_rule":         resourceRule(),
			"hue_resourcelink": resourceResourcelink(),
		},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	address := d.Get("hostname").(string)
	username := d.Get("username").(string)
	hue := &Hue{
		Hostname: address,
		Username: username,
	}

	return hue, nil
}
