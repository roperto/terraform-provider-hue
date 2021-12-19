output "sensors_ids" {
  value = [hue_sensor.status.id]
}

output "scenes_ids" {
  value = [hue_scene.recover.id]
}

output "rules_ids" {
  value = [for r in local.rules : r.id]
}

output "resourcelinks_ids" {
  value = [hue_resourcelink.this.id]
}
