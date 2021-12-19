variable "name" {
  type = string
}

variable "uniqueid" {
  type = string
}

variable "sensor_light_threshold" {
  type    = number
  default = 16000
}

locals {
  prefix_cropped = substr(var.name, 0, 32 - length(" Motion Sensor"))

  motion_sensor_light       = "0400"
  motion_sensor_temperature = "0402"
  motion_sensor_motion      = "0406"
}

resource "hue_sensor" "motion" {
  uniqueid = "${var.uniqueid}-${local.motion_sensor_motion}"
  name     = "${local.prefix_cropped} Motion Sensor"
  type     = "ZLLPresence"
}

resource "hue_sensor" "temperature" {
  uniqueid = "${var.uniqueid}-${local.motion_sensor_temperature}"
  name     = "${local.prefix_cropped} Temp Sensor"
  type     = "ZLLTemperature"
}

resource "hue_sensor_light" "light" {
  uniqueid = "${var.uniqueid}-${local.motion_sensor_light}"
  name     = "${local.prefix_cropped} Light Sensor"
  type     = "ZLLLightLevel"

  config_tholddark = var.sensor_light_threshold
}

output "motion_sensor" {
  value = hue_sensor.motion
}

output "temperature_sensor" {
  value = hue_sensor.temperature
}

output "light_sensor" {
  value = hue_sensor_light.light
}

output "sensors_ids" {
  value = [
    hue_sensor.motion.id,
    hue_sensor.temperature.id,
    hue_sensor_light.light.id,
  ]
}