variable "presence_sensor" {
  type    = string
  default = ""
}

variable "use_presence_sensor" {
  type    = bool
  default = true
}

locals {
  use_presence_sensor = var.use_presence_sensor && var.presence_sensor != ""
}

module "presence" {
  source = "../presence"
  count  = var.presence_sensor == "" ? 0 : 1

  name     = var.name
  uniqueid = var.presence_sensor
}

module "motion_group" {
  source = "../motion_group"
  count  = local.use_presence_sensor ? 1 : 0

  name          = substr(var.name, 0, 18)
  defaults      = var.config
  group         = hue_group.this
  motion_sensor = module.presence[0].motion_sensor
  light_sensor  = module.presence[0].light_sensor
  day_scene     = hue_scene.this["Cool"].id
  night_scene   = hue_scene.this["Warm"].id
}

variable "dimmer" {
  type    = string
  default = ""
}

module "dimmer" {
  source = "../dimmer"
  count  = var.dimmer == "" ? 0 : 1

  name             = var.name
  uniqueid         = var.dimmer
  group            = hue_group.this
  day_scene        = hue_scene.this["Cool"].id
  night_scene      = hue_scene.this["Warm"].id
  nightlight_scene = hue_scene.this["Nightlight"].id
  config           = var.config
}
