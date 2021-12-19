resource "hue_rule" "on-press-night" {
  name = "DS${hue_sensor.dimmer.bridge_id} on-pn ${var.name}"
  conditions = jsonencode([
    {
      address  = "/config/localtime"
      operator = "not in"
      value    = var.config.time_day
    },
    {
      address  = "${hue_sensor.dimmer.address}/state/buttonevent"
      operator = "eq"
      value    = "1000"
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
        scene = var.night_scene
      }
    },
  ])
}
