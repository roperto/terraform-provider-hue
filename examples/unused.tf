# You can use this strategy if you want to remove extra stuff from your bridge.
# It will not do it for you, you will still need to delete stuff using (example):
# curl -X DELETE http://hue/api/[username]/[resource]/[id]

locals {
  rooms = [
    module.bedroom,
    module.hallway,
    module.living,
    module.study,
  ]
}

locals {
  groups_ids = [for r in local.rooms : r.group.id]

  lights_ids = setunion(
    [for r in local.rooms : r.lights_ids]...
  )

  scenes_ids = setunion(
    setunion([for r in local.rooms : r.scenes_ids]...),
  )

  sensors_ids = setunion(
    concat(
      [for r in local.rooms : r.sensors_ids],
    )...
  )

  rules_ids = setunion(
    concat(
      [for r in local.rooms : r.rules_ids],
    )...
  )

  resourcelinks_ids = setunion(
    concat(
      [for r in local.rooms : r.resourcelinks_ids],
    )...
  )
}

data "hue_groups" "found" {}
locals {
  unused_groups_1 = { for g in data.hue_groups.found.groups : g.bridge_id => g }
  unused_groups_2 = toset([for g in data.hue_groups.found.groups : g.bridge_id])
  unused_groups_3 = setsubtract(local.unused_groups_2, local.groups_ids)
}
output "unused_groups" {
  value = [
    for id in local.unused_groups_3 :
    "Group #${id}: ${local.unused_groups_1[id].name}"
  ]
}

data "hue_lights" "found" {}
locals {
  unused_lights_1 = { for l in data.hue_lights.found.lights : l.uniqueid => l }
  unused_lights_2 = toset([for l in data.hue_lights.found.lights : l.uniqueid])
  unused_lights_3 = setsubtract(local.unused_lights_2, local.lights_ids)
}
output "unused_lights" {
  value = [
    for id in local.unused_lights_3 :
    "Light '${id}' #${local.unused_lights_1[id].bridge_id}: ${local.unused_lights_1[id].name}"
  ]
}


data "hue_resourcelinks" "found" {}
locals {
  unused_resourcelinks_1 = { for r in data.hue_resourcelinks.found.resourcelinks : r.bridge_id => r }
  unused_resourcelinks_2 = toset([for r in data.hue_resourcelinks.found.resourcelinks : r.bridge_id])
  unused_resourcelinks_3 = setsubtract(local.unused_resourcelinks_2, local.resourcelinks_ids)
}
output "unused_resourcelinks" {
  value = [
    for id in local.unused_resourcelinks_3 :
    "ResourceLink #${id} '${local.unused_resourcelinks_1[id].name}': ${local.unused_resourcelinks_1[id].description}"
  ]
}

data "hue_rules" "found" {}
locals {
  unused_rules_1 = { for r in data.hue_rules.found.rules : r.bridge_id => r }
  unused_rules_2 = toset([for r in data.hue_rules.found.rules : r.bridge_id])
  unused_rules_3 = setsubtract(local.unused_rules_2, local.rules_ids)
}
output "unused_rules" {
  value = [
    for id in local.unused_rules_3 :
    "Rule #${id}: ${local.unused_rules_1[id].name}"
  ]
}

data "hue_scenes" "found" {}
locals {
  unused_scenes_1 = { for s in data.hue_scenes.found.scenes : s.scene_id => s }
  unused_scenes_2 = toset([for s in data.hue_scenes.found.scenes : s.scene_id])
  unused_scenes_3 = setsubtract(local.unused_scenes_2, local.scenes_ids)
}
output "unused_scenes" {
  value = [
    for id in local.unused_scenes_3 :
    "Scene '${id}' @ '${local.unused_scenes_1[id].group_name}' (${local.unused_scenes_1[id].group_id}): ${local.unused_scenes_1[id].scene_name}"
  ]
}

data "hue_schedules" "found" {}
output "unused_schedules" {
  value = tolist(data.hue_schedules.found.schedules)
}

data "hue_sensors" "found" {}
locals {
  unused_sensors_1 = { for l in data.hue_sensors.found.sensors : l.uniqueid => l }
  unused_sensors_2 = toset([for l in data.hue_sensors.found.sensors : l.uniqueid])
  unused_sensors_3 = setsubtract(local.unused_sensors_2, local.sensors_ids)
}
output "unused_sensors" {
  value = [
    for id in local.unused_sensors_3 :
    "Sensor ${local.unused_sensors_1[id].type} '${id}' #${local.unused_sensors_1[id].bridge_id}: ${local.unused_sensors_1[id].name}"
    if local.unused_sensors_1[id].type != "Daylight"
  ]
}