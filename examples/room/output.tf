output "lights" {
  value = hue_light.lights
}

output "scenes" {
  value = hue_scene.this
}

output "group" {
  value = hue_group.this
}

output "dimmer" {
  value = length(module.dimmer) == 1 ? module.dimmer[0].this : null
}

output "lights_ids" {
  value = toset([for l in hue_light.lights : l.id])
}

output "scenes_ids" {
  value = toset(concat(
    [for s in hue_scene.this : s.id],
    concat([], [for m in module.motion_group : m.scenes_ids]...),
  ))
}

output "sensors_ids" {
  value = toset(concat(
    concat([], [for p in module.presence : p.sensors_ids]...),
    concat([], [for m in module.motion_group : m.sensors_ids]...),
    concat([], [for m in module.dimmer : m.sensors_ids]...),
  ))
}

output "rules_ids" {
  value = toset(concat(
    concat([], [for m in module.motion_group : m.rules_ids]...),
    concat([], [for m in module.dimmer : m.rules_ids]...),
  ))
}

output "resourcelinks_ids" {
  value = toset(concat(
    concat([], [for m in module.motion_group : m.resourcelinks_ids]...),
    concat([], [for m in module.dimmer : m.resourcelinks_ids]...),
  ))
}
