
##############################
# Resource: clouding_sshkey  #
##############################

resource "clouding_sshkey" "example" {
  name       = "sshkey"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDmDTWnz2hR4byvUr9a2vaOW5BuArZorY60Sk7CfFgeay4oIMDTRWURQaFKWc5NqiE/Q/cvWO8MOo6v0ji7OzNysERRic6NoaS0kEY7gjFvyvvojU6jHN8yBogEmLKCdt4OY3LqJ1FV4ptqRovOJyxanNnEpJBrbkFxzPP5N3n/WGuXRN9KFSJXp76NTVQ68tfCB4bmkXQyhWbFKKkVKqUyPlVVEGVuCMGVw6GvSdz/meIaVdDpJSmhEm5KX5Mv4mg6udRJS5N+Bzq4iVkBDQUSf5nMZwH32volP07nnCvgGNENmcJMiMkUV4L5uUFOqUhgPyj/6kxwkEyzG974C6K5"
}

################################
# Data source: clouding_sshkey #
################################

data "clouding_sshkey" "test" {
  id = clouding_sshkey.test.id
}
