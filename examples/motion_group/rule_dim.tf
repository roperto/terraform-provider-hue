resource "hue_rule" "dim" {
  name = "MG${var.motion_sensor.bridge_id} dim ${var.name}"
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
        value    = "PT${var.defaults.delay_dim}"
      },
      {
        address  = "${hue_sensor.status.address}/state/status"
        operator = "eq"
        value    = tostring(local.status_on)
      },
    ]
  )
  actions = jsonencode(
    [
      {
        method  = "PUT"
        address = hue_scene.recover.address
        body = {
          storelightstate = true
        }
      },
      {
        method  = "PUT"
        address = "${hue_sensor.status.address}/state"
        body = {
          status = local.status_dimmed
        }
      },
      {
        method  = "PUT"
        address = "${var.group.address}/action"
        body = {
          bri_inc = -var.defaults.dim_brightness
        }
      },
    ]
  )
}
