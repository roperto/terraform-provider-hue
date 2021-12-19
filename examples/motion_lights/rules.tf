locals {
  rules = [
    hue_rule.dim,
    hue_rule.recover,
    hue_rule.off,
    hue_rule.reset,
    hue_rule.on-presence-daytime,
    hue_rule.on-dark-daytime,
    hue_rule.on-presence-nighttime,
    hue_rule.on-dark-nighttime,
  ]
}
