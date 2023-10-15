package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFirewallRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFirewallRuleConfig("0.0.0.0/0", "tcp", "Allow http connections", "80", "80"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("clouding_firewall_rule.test", "source_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("clouding_firewall_rule.test", "description", "Allow http connections"),
					resource.TestCheckResourceAttr("clouding_firewall_rule.test", "port_range_min", "80"),
					resource.TestCheckResourceAttr("clouding_firewall_rule.test", "port_range_max", "80"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "clouding_firewall_rule.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: testAccFirewallRuleConfig("10.0.0.0/0", "tcp", "Allow http connections two", "8080", "8080"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("clouding_firewall_rule.test", "source_ip", "10.0.0.0/0"),
					resource.TestCheckResourceAttr("clouding_firewall_rule.test", "description", "Allow http connections two"),
					resource.TestCheckResourceAttr("clouding_firewall_rule.test", "port_range_min", "8080"),
					resource.TestCheckResourceAttr("clouding_firewall_rule.test", "port_range_max", "8080"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccFirewallRuleConfig(source_ip, protocol, description, portMin, portMax string) string {
	return fmt.Sprintf(`
resource "clouding_firewall" "test" {
	name = "testacc-firewall-rule"
	description = "testacc-firewall-rule description"
}
resource "clouding_firewall_rule" "test" {
	firewall_id = clouding_firewall.test.id 
	source_ip = "%s"
	protocol = "%s"
	description = "%s"
	port_range_min = "%s" 
	port_range_max = "%s" 
}
`, source_ip, protocol, description, portMin, portMax)
}
