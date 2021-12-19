resource "hue_rule" "on-dark-daytime" {
  name = "ML${var.motion_sensor.bridge_id} on-dd ${var.name}"
  conditions = jsonencode([
    {
      address  = "/config/localtime"
      operator = "in"
      value    = var.config.time_day
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
      address  = "${var.light_sensor.address}/state/dark"
      operator = "dx"
    },
  ])
  actions = jsonencode(concat(
    [
      {
        address = "${hue_sensor.status.address}/state"
        method  = "PUT"
        body = {
          status = local.status_on
        }
      },
    ],
    [
      for address in var.lights_addresses :
      {
        address = "${address}/state"
        method  = "PUT"
        body    = merge({on = true}, local.light_day)
      }
    ],
  ))
}
