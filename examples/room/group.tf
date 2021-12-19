variable "name" {
  type = string
}

resource "hue_group" "this" {
  name   = var.name
  type   = "Room"
  lights = [for l in hue_light.lights : l.bridge_id]
}
