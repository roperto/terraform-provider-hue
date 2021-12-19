variable "config" {
  type = object({
    time_day       = string
    delay_dim      = string
    delay_off      = string
    dim_brightness = number

    scenes = map(object({
      xy   = list(number)
      bri  = number
      icon = string
    }))
  })
}
