resource "hue_scene" "recover" {
  name         = "Motion sensor ${var.motion_sensor.bridge_id} recover"
  group        = var.group.id
  appdata_data = "tfhue_ms_rec_h01"

  dynamic "lightstate" {
    for_each = var.group.lights
    content {
      id  = lightstate.value
      on  = false
      bri = 1
      x   = 0.5
      y   = 0.5
    }
  }

  lifecycle {
    ignore_changes = [lightstate]
  }
}
