resource "hue_rule" "reset" {
  name = "ML${var.motion_sensor.bridge_id} reset ${var.name}"
  conditions = jsonencode(concat(
    [
      for address in var.lights_addresses :
      {
        address  = "${address}/state/on"
        operator = "eq"
        value    = "false"
      }
    ],
    [
      {
        address  = "${var.motion_sensor.address}/state/presence"
        operator = "eq"
        value    = "false"
      },
      {
        address  = "${var.motion_sensor.address}/state/presence"
        operator = "dx"
      },
    ]
  ))
  actions = jsonencode([
    {
      address = "/sensors/${hue_sensor.status.bridge_id}/state"
      method  = "PUT"
      body = {
        status = local.status_reset
      }
    },
  ])
}
