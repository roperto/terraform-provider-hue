variable "defaults" {
  type = object({
    dim_brightness = number
    delay_dim      = string
    delay_off      = string
    time_day       = string
  })
}
