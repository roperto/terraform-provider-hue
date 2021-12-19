package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func resourceScene() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSceneCreate,
		ReadContext:   resourceSceneRead,
		UpdateContext: resourceSceneUpdate,
		DeleteContext: resourceSceneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"recycle": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"appdata_version": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"appdata_data": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "tfhue_rXX_dXX",
			},
			"lightstate": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"on": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"bri": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"x": &schema.Schema{
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"y": &schema.Schema{
							Type:     schema.TypeFloat,
							Optional: true,
						},
					},
				},
			},
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSceneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	scene := map[string]interface{}{
		"name":        d.Get("name").(string),
		"type":        "GroupScene",
		"group":       d.Get("group").(string),
		"recycle":     false,
		"lightstates": unflattenSceneLightStates(d.Get("lightstate").(*schema.Set)),
		"appdata": map[string]interface{}{
			"version": d.Get("appdata_version"),
			"data":    d.Get("appdata_data"),
		},
	}

	hue := m.(*Hue)
	id, err := hue.Create(ctx, ResourceTypeScene, scene)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return resourceSceneRead(ctx, d, m)
}

func resourceSceneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hue := m.(*Hue)
	scene, err := hue.Read(ctx, ResourceTypeScene, d.Id())
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

	appdata, ok := scene["appdata"]
	if ok {
		appdata := appdata.(map[string]interface{})
		d.Set("appdata_version", appdata["version"])
		d.Set("appdata_data", appdata["data"])
	} else {
		d.Set("appdata_version", nil)
		d.Set("appdata_data", nil)
	}

	d.Set("name", scene["name"])
	d.Set("type", scene["type"])
	d.Set("group", scene["group"])
	d.Set("recycle", scene["recycle"])
	d.Set("address", fmt.Sprintf("/%s/%s", ResourceTypeScene, d.Id()))

	lightstateMap := scene["lightstates"].(map[string]interface{})
	lightstates := make([]map[string]interface{}, len(lightstateMap))
	i := 0
	for bridgeId, values := range lightstateMap {
		values := values.(map[string]interface{})
		lightstate := map[string]interface{}{}

		var v interface{}
		var ok bool

		lightstate["id"] = bridgeId

		items := []string{"on", "bri"}
		for _, item := range items {
			v, ok := values[item]
			if ok {
				lightstate[item] = v
			} else {
				lightstate[item] = ""
			}
		}

		v, ok = values["xy"]
		if ok {
			lightstate["x"] = v.([]interface{})[0].(float64)
			lightstate["y"] = v.([]interface{})[1].(float64)
		}

		lightstates[i] = lightstate
		i++
	}

	d.Set("lightstate", lightstates)

	return nil
}

func resourceSceneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	changes := map[string]interface{}{}
	if d.HasChange("name") {
		changes["name"] = d.Get("name").(string)
	}
	if d.HasChange("lightstate") {
		changes["lightstates"] = unflattenSceneLightStates(d.Get("lightstate").(*schema.Set))
	}
	if d.HasChange("appdata_version") || d.HasChange("appdata_data") {
		changes["appdata"] = map[string]interface{}{
			"version": d.Get("appdata_version"),
			"data":    d.Get("appdata_data"),
		}
	}

	hue := m.(*Hue)
	err := hue.Update(ctx, ResourceTypeScene, d.Id(), changes)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSceneRead(ctx, d, m)
}

func resourceSceneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hue := m.(*Hue)

	err := hue.Delete(ctx, ResourceTypeScene, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func unflattenSceneLightStates(data *schema.Set) map[int]map[string]interface{} {
	lightstates := map[int]map[string]interface{}{}

	for _, entrymap := range data.List() {
		entry := entrymap.(map[string]interface{})
		state := map[string]interface{}{}

		var ok bool
		var v interface{}

		v, ok = entry["on"]
		if ok {
			state["on"] = v.(bool)
		}

		vX, okX := entry["x"]
		vY, okY := entry["y"]
		if okX && okY {
			state["xy"] = []float64{vX.(float64), vY.(float64),
			}
		}

		v, ok = entry["bri"]
		if ok {
			state["bri"] = parseInt(v)
		}

		lightId, _ := strconv.Atoi(entry["id"].(string))
		lightstates[lightId] = state
	}
	return lightstates
}
