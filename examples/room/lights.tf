variable "lights" {
  type    = list(string)
  default = []
}

variable "named_lights" {
  type    = map(list(string))
  default = {}
}

locals {
  all_lights = merge([
    for name, ids in merge(
      { "" : var.lights },
      var.named_lights
    ) :
    {
      for index, id in ids :
      id => format(
        "${var.name}%s%s%s",
        name == "" ? "" : " ${name}",
        lower(substr(name, -5, -1)) == "light" ? "" : " Light",
        length(ids) == 1 ? "" : format(" %d", index + 1)
      )
    }
  ]...)
}

resource "hue_light" "lights" {
  for_each = local.all_lights

  name     = each.value
  uniqueid = each.key
}
