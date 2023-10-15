
####################################
# Resource: clouding_firewall_rule #
#################################### 

resource "clouding_firewall" "example" {
  name        = "example test"
  description = "Firewall example"
}

resource "clouding_firewall_rule" "example_rule" {
  firewall_id    = clouding_firewall.example.id
  source_ip      = "0.0.0.0/0"
  protocol       = "tcp"
  description    = "Allow http connections"
  port_range_min = 80
  port_range_max = 80
}
