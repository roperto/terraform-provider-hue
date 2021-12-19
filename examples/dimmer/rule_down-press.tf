resource "hue_rule" "down-press" {
  name = "DS${hue_sensor.dimmer.bridge_id} dn-p ${var.name}"
  conditions = jsonencode([
    {
      address  = "${hue_sensor.dimmer.address}/state/buttonevent"
      operator = "eq"
      value    = "3000"
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
        bri_inc        = -30
        transitiontime = 9
      }
    },
  ])
}
