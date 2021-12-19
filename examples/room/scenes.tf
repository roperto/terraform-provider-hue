variable "scenes" {
  type    = bool
  default = true
}

locals {
  appdata = {
    Relax       = "tfhue_r01_d01",
    Read        = "tfhue_r01_d02",
    Concentrate = "tfhue_r01_d03",
    Energise    = "tfhue_r01_d04",
    Bright      = "tfhue_r01_d05",
    Dimmed      = "tfhue_r01_d06",
    Nightlight  = "tfhue_r01_d07",
  }
  use_scenes = var.scenes ? var.config.scenes : {}
}

resource "hue_scene" "this" {
  for_each = local.use_scenes

  group        = hue_group.this.id
  name         = each.key
  appdata_data = local.appdata[each.value.icon]

  dynamic "lightstate" {
    for_each = hue_light.lights
    content {
      id  = lightstate.value.bridge_id
      on  = true
      x   = each.value.xy[0]
      y   = each.value.xy[1]
      bri = each.value.bri
    }
  }
}
