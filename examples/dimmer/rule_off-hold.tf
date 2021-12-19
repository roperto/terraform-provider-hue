resource "hue_rule" "off-hold" {
  name = "DS${hue_sensor.dimmer.bridge_id} of-h ${var.name}"
  conditions = jsonencode([
    {
      address  = "${hue_sensor.dimmer.address}/state/buttonevent"
      operator = "eq"
      value    = "4001"
    },
    {
      address  = "${hue_sensor.dimmer.address}/state/lastupdated"
      operator = "dx"
    },
  ])
  actions = jsonencode([
    {
      method  = "PUT"
      address = "${var.group.address}/action"
      body = {
        scene = var.nightlight_scene
      }
    },
  ])
}
