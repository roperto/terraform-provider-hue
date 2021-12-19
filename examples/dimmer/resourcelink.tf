resource "hue_resourcelink" "link" {
  name        = "Sensor ${hue_sensor.dimmer.bridge_id}"
  description = "Resources of sensor ${hue_sensor.dimmer.bridge_id}"
  classid     = 10011 # From Hue Essentials App
  links = concat(
  [for r in local.rules : r.address],
  [hue_sensor.dimmer.address],
  )
}
