package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImageDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccImageConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.clouding_image.test", "name", "Debian 11 (64 Bit)"),
				),
			},
		},
	})
}

const testAccImageConfig = `
data "clouding_image" "test" {
	id = "wLQbN5nvg829JaeZ"
}
`
