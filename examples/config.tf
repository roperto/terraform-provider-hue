locals {
  config = {
    time_day       = "T07:00:00/T19:00:00"
    delay_dim      = "00:05:00"
    delay_off      = "00:00:30"
    dim_brightness = 127
    scenes = {
      Warm = {
        xy   = [0.445, 0.4067]
        bri  = 254
        icon = "Bright"
      }
      Cool = {
        xy   = [0.3143, 0.3301]
        bri  = 254
        icon = "Concentrate"
      }
      Nightlight = {
        xy   = [0.5612, 0.4042]
        bri  = 1
        icon = "Nightlight"
      }
    }
  }
}
