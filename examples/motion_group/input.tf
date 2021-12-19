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

variable "group" {
  type = any
}

variable "day_scene" {
  type = string
}

variable "night_scene" {
  type = string
}
