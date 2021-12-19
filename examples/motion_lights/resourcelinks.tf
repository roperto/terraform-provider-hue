resource "hue_resourcelink" "this" {
  name        = "ML${var.motion_sensor.bridge_id} ${var.name} "
  description = "Motion Lights on Sensor ${var.motion_sensor.bridge_id} for ${var.name}"
  classid     = 10020
  links = concat(
    [for r in local.rules : r.address],
    [
      var.motion_sensor.address,
      var.light_sensor.address,
    ]
  )
}
