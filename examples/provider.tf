terraform {
  required_providers {
    hue = {
      source = "github.com/roperto/hue"
    }
  }
}

provider "hue" {
  hostname = "192.168.xxx.yyy"
  username = "Check README.md"
}
