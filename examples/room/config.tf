variable "config" {
  type = object({
    dim_brightness = number
    delay_dim      = string
    delay_off      = string
    time_day       = string

    scenes = map(object({
      xy   = list(number)
      bri  = number
      icon = string
    }))
  })
}
