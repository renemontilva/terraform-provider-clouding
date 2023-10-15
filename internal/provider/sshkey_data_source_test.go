package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSshKeyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccSshKeyDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.clouding_sshkey.test", "name", "test"),
					resource.TestCheckResourceAttr("data.clouding_sshkey.test", "public_key", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDmDTWnz2hR4byvUr9a2vaOW5BuArZorY60Sk7CfFgeay4oIMDTRWURQaFKWc5NqiE/Q/cvWO8MOo6v0ji7OzNysERRic6NoaS0kEY7gjFvyvvojU6jHN8yBogEmLKCdt4OY3LqJ1FV4ptqRovOJyxanNnEpJBrbkFxzPP5N3n/WGuXRN9KFSJXp76NTVQ68tfCB4bmkXQyhWbFKKkVKqUyPlVVEGVuCMGVw6GvSdz/meIaVdDpJSmhEm5KX5Mv4mg6udRJS5N+Bzq4iVkBDQUSf5nMZwH32volP07nnCvgGNENmcJMiMkUV4L5uUFOqUhgPyj/6kxwkEyzG974C6K5"),
				),
			},
		},
	})
}

const testAccSshKeyDataSourceConfig = `
resource "clouding_sshkey" "test" {
	name = "test"
	public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDmDTWnz2hR4byvUr9a2vaOW5BuArZorY60Sk7CfFgeay4oIMDTRWURQaFKWc5NqiE/Q/cvWO8MOo6v0ji7OzNysERRic6NoaS0kEY7gjFvyvvojU6jHN8yBogEmLKCdt4OY3LqJ1FV4ptqRovOJyxanNnEpJBrbkFxzPP5N3n/WGuXRN9KFSJXp76NTVQ68tfCB4bmkXQyhWbFKKkVKqUyPlVVEGVuCMGVw6GvSdz/meIaVdDpJSmhEm5KX5Mv4mg6udRJS5N+Bzq4iVkBDQUSf5nMZwH32volP07nnCvgGNENmcJMiMkUV4L5uUFOqUhgPyj/6kxwkEyzG974C6K5"
}
  
data "clouding_sshkey" "test" {
	id = clouding_sshkey.test.id
}
`
