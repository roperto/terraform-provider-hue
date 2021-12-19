# A room with one light.
module "bedroom" {
  source = "./room"

  name   = "Bedroom"
  config = local.config
  lights          = ["00:17:88:01:xx:xx:xx:xx-0b"]
}

# A room with two named lights.
module "study" {
  source = "./room"

  name   = "Study"
  config = local.config
  named_lights = {
    Ceiling = ["00:17:88:01:xx:xx:xx:xx-0b"],
    Desk = ["00:17:88:01:xx:xx:xx:xx-0b"],
  }
}

# A room using motion sensor.
module "hallway" {
  source = "./room"

  name            = "Hallway"
  config          = local.config
  lights          = ["00:17:88:01:xx:xx:xx:xx-0b"]
  presence_sensor = "00:17:88:01:xx:xx:xx:xx-02"
}

# A room with a dimmer switch.
module "living" {
  source = "./room"

  name   = "Lounge"
  config = local.config
  lights          = ["00:17:88:01:xx:xx:xx:xx-0b"]
  dimmer = "00:17:88:01:xx:xx:xx:xx-02-fc00"
}
