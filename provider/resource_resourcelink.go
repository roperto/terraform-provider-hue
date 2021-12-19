package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceResourcelink() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResourcelinkCreate,
		ReadContext:   resourceResourcelinkRead,
		UpdateContext: resourceResourcelinkUpdate,
		DeleteContext: resourceResourcelinkDelete,
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
			"classid": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"links": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceResourcelinkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	links_got := d.Get("links").(*schema.Set).List()
	links := make([]string, len(links_got))
	for k, v := range links_got {
		links[k] = v.(string)
	}

	data := map[string]interface{}{
		"name":        d.Get("name"),
		"classid":       d.Get("classid"),
		"description": d.Get("description"),
		"links":       links,
	}

	hue := m.(*Hue)
	id, err := hue.Create(ctx, ResourceTypeResourcelink, data)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return resourceResourcelinkRead(ctx, d, m)
}

func resourceResourcelinkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hue := m.(*Hue)

	link, err := hue.Read(ctx, ResourceTypeResourcelink, d.Id())
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

	d.Set("name", link["name"])
	d.Set("classid", link["classid"])
	d.Set("description", link["description"])
	d.Set("links", link["links"])

	return nil
}

func resourceResourcelinkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	changes := map[string]interface{}{}
	if d.HasChange("name") {
		changes["name"] = d.Get("name").(string)
	}
	if d.HasChange("description") {
		changes["description"] = d.Get("description").(string)
	}
	if d.HasChange("classid") {
		changes["classid"] = d.Get("classid").(int)
	}
	if d.HasChange("links") {
		links_got := d.Get("links").(*schema.Set).List()
		links := make([]string, len(links_got))
		for k, v := range links_got {
			links[k] = v.(string)
		}
		changes["links"] = links
	}

	hue := m.(*Hue)
	err := hue.Update(ctx, ResourceTypeResourcelink, d.Id(), changes)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceResourcelinkRead(ctx, d, m)
}

func resourceResourcelinkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hue := m.(*Hue)

	err := hue.Delete(ctx, ResourceTypeResourcelink, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
