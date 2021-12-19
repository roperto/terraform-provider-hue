resource "hue_rule" "off" {
  name = "MG${var.motion_sensor.bridge_id} off ${var.name}"
  conditions = jsonencode([
    {
      address  = "${var.motion_sensor.address}/state/presence"
      operator = "eq"
      value    = "false"
    },
    {
      address  = "${hue_sensor.status.address}/state/status"
      operator = "gt"
      value    = tostring(local.status_on)
    },
    {
      address  = "${hue_sensor.status.address}/state/status"
      operator = "ddx"
      value    = "PT${var.defaults.delay_off}"
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
    {
      method  = "PUT"
      address = "${hue_sensor.status.address}/state"
      body = {
        status = local.status_off
      }
    },
  ])
}
