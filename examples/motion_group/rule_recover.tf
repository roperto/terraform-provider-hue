resource "hue_rule" "recover" {
  name = "MG${var.motion_sensor.bridge_id} recover ${var.name}"
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
  actions = jsonencode([
    {
      method  = "PUT"
      address = "${var.group.address}/action"
      body = {
        scene = hue_scene.recover.id
      }
    },
    {
      method  = "PUT"
      address = "${hue_sensor.status.address}/state"
      body = {
        status = local.status_on
      }
    },
  ])
}
