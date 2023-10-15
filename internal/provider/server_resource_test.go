package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccServerConfig("testacc", "testacc01"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("clouding_server.test", "name", "testacc"),
					resource.TestCheckResourceAttr("clouding_server.test", "hostname", "testacc01"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "clouding_server.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: testAccServerConfig("testacc2", "testacc02"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("clouding_server.test", "name", "testacc2"),
					resource.TestCheckResourceAttr("clouding_server.test", "hostname", "testacc02"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccServerConfig(name, hostname string) string {
	return fmt.Sprintf(`
resource "clouding_server" "test" {
  name = "%s"
  hostname = "%s"
  flavor_id = "0.5x1"
  firewall_id = "L1qX02j9agnW9ary"

  access_configuration = {
    password = "test1234"
    save_password = true 
  }

  volume = {
    source = "image"
    id = "wLQbN5nvg829JaeZ" #debian 11
    ssd_gb = 5
  }

  backup_preference = {
    slots = 3
    frequency = "ThreeDays"
  }
}
`, name, hostname)
}
