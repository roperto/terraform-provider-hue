resource "hue_rule" "reset" {
  name = "MG${var.motion_sensor.bridge_id} reset ${var.name}"
  conditions = jsonencode([
    {
      address  = "${var.group.address}/state/any_on"
      operator = "eq"
      value    = "false"
    },
    {
      address  = "${var.motion_sensor.address}/state/presence"
      operator = "eq"
      value    = "false"
    },
  ])
  actions = jsonencode([
    {
      method  = "PUT"
      address = "${hue_sensor.status.address}/state"
      body = {
        status = local.status_off
      }
    }
  ])
}
