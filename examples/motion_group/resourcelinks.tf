resource "hue_resourcelink" "this" {
  name        = "Sensor ${var.motion_sensor.bridge_id}"
  description = "Resources of sensor ${var.motion_sensor.bridge_id}"
  classid     = 10020
  links = concat(
    [for r in local.rules : r.address],
    [
      hue_scene.recover.address,
      var.motion_sensor.address,
      var.light_sensor.address,
    ]
  )
}
