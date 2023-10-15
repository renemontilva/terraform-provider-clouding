package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFirewallDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testFirewallDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.clouding_firewall.test", "id", "L1qX02j9agnW9ary"),
					resource.TestCheckResourceAttr("data.clouding_firewall.test", "name", "default"),
					resource.TestCheckResourceAttr("data.clouding_firewall.test", "description", "Default security group"),
				),
			},
		},
	})
}

const testFirewallDataSourceConfig = `
data "clouding_firewall" "test" {
  id = "L1qX02j9agnW9ary"
}
`
