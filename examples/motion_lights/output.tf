output "sensors_ids" {
  value = [hue_sensor.status.id]
}

output "rules_ids" {
  value = [for r in local.rules : r.id]
}

output "resourcelinks_ids" {
  value = [hue_resourcelink.this.id]
}
