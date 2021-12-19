resource "hue_rule" "recover" {
  name = "ML${var.motion_sensor.bridge_id} recover ${var.name}"
  conditions = jsonencode([
    {
      address  = "${hue_sensor.status.address}/state/status"
      operator = "gt"
      value    = tostring(local.status_on)
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
  ])
  actions = jsonencode(concat(
    [
      {
        method  = "PUT"
        address = "/sensors/${hue_sensor.status.bridge_id}/state"
        body = {
          status = local.status_on
        }
      },
    ],
    [
      for address in var.lights_addresses :
      {
        method  = "PUT"
        address = "${address}/state"
        body = {
          bri_inc = local.dim_brightness_by
        }
      }
    ],
  ))
}