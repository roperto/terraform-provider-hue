resource "hue_rule" "off-press" {
  name = "DS${hue_sensor.dimmer.bridge_id} off-p ${var.name}"
  conditions = jsonencode([
    {
      address  = "${hue_sensor.dimmer.address}/state/buttonevent"
      operator = "eq"
      value    = "4000"
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
        on = false
      }
    },
  ])
}
