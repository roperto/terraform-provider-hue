locals {
  status_off    = 0
  status_on     = 1
  status_dimmed = 2
}

resource "hue_sensor" "status" {
  uniqueid = "motion_group_status_sensor_${var.motion_sensor.bridge_id}"
  name     = "MG${var.motion_sensor.bridge_id} ${var.name}"
  type     = "CLIPGenericStatus"
}
