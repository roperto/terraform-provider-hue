resource "hue_rule" "on-presence-nighttime" {
  name = "MG${var.motion_sensor.bridge_id} on-pn ${var.name}"
  conditions = jsonencode(
    [
      {
        address  = "/config/localtime"
        operator = "not in"
        value    = var.defaults.time_day
      },
      {
        address  = "${hue_sensor.status.address}/state/status"
        operator = "lt"
        value    = tostring(local.status_on)
      },
      {
        address  = "${var.light_sensor.address}/state/dark"
        operator = "eq"
        value    = "true"
      },
      {
        address  = "${var.motion_sensor.address}/state/presence"
        operator = "eq"
        value    = "true"
      },
      {
        address  = "${var.motion_sensor.address}/state/presence"
        operator = "dx"
      },
    ]
  )
  actions = jsonencode([
    {
      method  = "PUT"
      address = "${hue_sensor.status.address}/state"
      body = {
        status = local.status_on
      }
    },
    {
      method  = "PUT"
      address = "${var.group.address}/action"
      body = {
        scene = var.night_scene
      }
    },
  ])
}
