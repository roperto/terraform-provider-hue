variable "name" {
  type = string

  validation {
    condition     = length(var.name) <= 18
    error_message = "The name is too long."
  }
}

variable "motion_sensor" {
  type = any
}

variable "light_sensor" {
  type = any
}

variable "lights_addresses" {
  type = list(string)
}

variable "light_day" {
  type = object({
    xy   = list(number)
    bri  = number
    icon = string
  })
  default = null
}

variable "light_night" {
  type = object({
    xy   = list(number)
    bri  = number
    icon = string
  })
  default = null
}

locals {
  light_day   = coalesce(var.light_day, var.config.scenes["Cool"])
  light_night = coalesce(var.light_night, var.config.scenes["Warm"])

  // Ensure dimming does not go below 1 so it can be recovered properly.
  dim_brightness_by = min(
    var.config.dim_brightness,
    local.light_day.bri - 1,
    local.light_night.bri - 1,
  )
}
