locals {
  rules = [
    hue_rule.on-press-day,
    hue_rule.on-press-night,
    hue_rule.off-hold,
    hue_rule.up-press,
    hue_rule.up-hold,
    hue_rule.down-press,
    hue_rule.down-hold,
    hue_rule.off-press,
  ]
}

output "this" {
  value = hue_sensor.dimmer
}

output "sensors_ids" {
  value = [
    hue_sensor.dimmer.id,
    hue_sensor.status.id,
  ]
}

output "rules_ids" {
  value = [for r in local.rules : r.id]
}

output "resourcelinks_ids" {
  value = [hue_resourcelink.link.id]
}
