resource "hue_rule" "down-hold" {
  name = "DS${hue_sensor.dimmer.bridge_id} dn-h ${var.name}"
  conditions = jsonencode([
    {
      address  = "${hue_sensor.dimmer.address}/state/buttonevent"
      operator = "eq"
      value    = "3001"
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
        bri_inc        = -56
        transitiontime = 9
      }
    },
  ])
}
