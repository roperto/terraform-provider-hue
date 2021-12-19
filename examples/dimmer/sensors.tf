resource "hue_sensor" "dimmer" {
  uniqueid = var.uniqueid
  type     = "ZLLSwitch"
  name     = "${var.name} Dimmer"
}

resource "hue_sensor" "status" {
  uniqueid = "dimmer_switch_status_sensor_${hue_sensor.dimmer.bridge_id}"
  name     = "Dimmer switch ${hue_sensor.dimmer.bridge_id} 1000 State"
  type     = "CLIPGenericStatus"
}
