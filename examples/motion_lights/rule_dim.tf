resource "hue_rule" "dim" {
  name = "ML${var.motion_sensor.bridge_id} dim ${var.name}"
  conditions = jsonencode(
    [
      {
        address  = "${var.motion_sensor.address}/state/presence"
        operator = "eq"
        value    = "false"
      },
      {
        address  = "${var.motion_sensor.address}/state/presence"
        operator = "ddx"
        value    = "PT${var.config.delay_dim}"
      },
      {
        address  = "${hue_sensor.status.address}/state/status"
        operator = "eq"
        value    = tostring(local.status_on)
      },
    ]
  )
  actions = jsonencode(concat(
    [
      {
        method  = "PUT"
        address = "${hue_sensor.status.address}/state"
        body = {
          status = local.status_dimmed
        }
      },
    ],
    [
      for address in var.lights_addresses :
      {
        method  = "PUT"
        address = "${address}/state"
        body = {
          bri_inc = -local.dim_brightness_by
        }
      }
    ],
  ))
}
