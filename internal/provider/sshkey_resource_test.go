package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/renemontilva/terraform-provider-clouding/internal/provider"
)

func TestAccSshKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"clouding": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSshKeyConfig("testacc"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("clouding_sshkey.test", "name", "testacc"),
					resource.TestCheckResourceAttr("clouding_sshkey.test", "public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDmDTWnz2hR4byvUr9a2vaOW5BuArZorY60Sk7CfFgeay4oIMDTRWURQaFKWc5NqiE/Q/cvWO8MOo6v0ji7OzNysERRic6NoaS0kEY7gjFvyvvojU6jHN8yBogEmLKCdt4OY3LqJ1FV4ptqRovOJyxanNnEpJBrbkFxzPP5N3n/WGuXRN9KFSJXp76NTVQ68tfCB4bmkXQyhWbFKKkVKqUyPlVVEGVuCMGVw6GvSdz/meIaVdDpJSmhEm5KX5Mv4mg6udRJS5N+Bzq4iVkBDQUSf5nMZwH32volP07nnCvgGNENmcJMiMkUV4L5uUFOqUhgPyj/6kxwkEyzG974C6K5"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "clouding_sshkey.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: testAccSshKeyConfig("testacc2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("clouding_sshkey.test", "name", "testacc2"),
					resource.TestCheckResourceAttr("clouding_sshkey.test", "public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDmDTWnz2hR4byvUr9a2vaOW5BuArZorY60Sk7CfFgeay4oIMDTRWURQaFKWc5NqiE/Q/cvWO8MOo6v0ji7OzNysERRic6NoaS0kEY7gjFvyvvojU6jHN8yBogEmLKCdt4OY3LqJ1FV4ptqRovOJyxanNnEpJBrbkFxzPP5N3n/WGuXRN9KFSJXp76NTVQ68tfCB4bmkXQyhWbFKKkVKqUyPlVVEGVuCMGVw6GvSdz/meIaVdDpJSmhEm5KX5Mv4mg6udRJS5N+Bzq4iVkBDQUSf5nMZwH32volP07nnCvgGNENmcJMiMkUV4L5uUFOqUhgPyj/6kxwkEyzG974C6K5"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccSshKeyConfig(name string) string {
	return fmt.Sprintf(`
resource "clouding_sshkey" "test" {
  name       = "%s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDmDTWnz2hR4byvUr9a2vaOW5BuArZorY60Sk7CfFgeay4oIMDTRWURQaFKWc5NqiE/Q/cvWO8MOo6v0ji7OzNysERRic6NoaS0kEY7gjFvyvvojU6jHN8yBogEmLKCdt4OY3LqJ1FV4ptqRovOJyxanNnEpJBrbkFxzPP5N3n/WGuXRN9KFSJXp76NTVQ68tfCB4bmkXQyhWbFKKkVKqUyPlVVEGVuCMGVw6GvSdz/meIaVdDpJSmhEm5KX5Mv4mg6udRJS5N+Bzq4iVkBDQUSf5nMZwH32volP07nnCvgGNENmcJMiMkUV4L5uUFOqUhgPyj/6kxwkEyzG974C6K5"
}
`, name)
}
