package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/renemontilva/terraform-provider-clouding/internal/provider"
)

func TestAccFirewallResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"clouding": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFirewallConfig("firewall-one", "testacc description one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("clouding_firewall.test", "name", "firewall-one"),
					resource.TestCheckResourceAttr("clouding_firewall.test", "description", "testacc description one"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "clouding_firewall.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: testAccFirewallConfig("testacc-firewall-two", "testacc description two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("clouding_firewall.test", "name", "testacc-firewall-two"),
					resource.TestCheckResourceAttr("clouding_firewall.test", "description", "testacc description two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccFirewallConfig(name, description string) string {
	return fmt.Sprintf(`
resource "clouding_firewall" "test" {
	name = "%s"
	description = "%s"
}
`, name, description)
}
