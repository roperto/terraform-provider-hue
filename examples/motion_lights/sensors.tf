locals {
  status_reset  = 0
  status_off    = 1
  status_on     = 2
  status_dimmed = 3
}

resource "hue_sensor" "status" {
  uniqueid = "motion_light_status_sensor_${var.motion_sensor.bridge_id}"
  name     = "ML${var.motion_sensor.bridge_id} ${var.name}"
  type     = "CLIPGenericStatus"
}
