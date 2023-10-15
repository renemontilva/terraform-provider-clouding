
#############################
# Resource: clouding_server #
#############################

##### Before creating a server, you can create a firewall and retrieve an image value from the data source. 

### Firewall
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

### Image
data "clouding_image" "example" {
  id = "wLQbN5nvg829JaeZ"
}


### Server
resource "clouding_server" "example" {
  name        = "webserver"
  hostname    = "webserver.example.com"
  flavor_id   = "0.5x1"
  firewall_id = clouding_firewall.example.id

  access_configuration = {
    password      = "password123"
    save_password = true
  }

  volume = {
    source = "image"
    id     = data.clouding_image.example.id
    ssd_gb = 5
  }

  backup_preference = {
    slots     = 3
    frequency = "ThreeDays"
  }
}
