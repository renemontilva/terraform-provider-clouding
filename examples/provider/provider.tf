# It is a good practice to use environment variables to store sensitive data
# like API keys, tokens, etc. This way, you can avoid committing them to your repository.

# First option, add a value to the token attribute
provider "clouding" {
  token = "token123"
}

# Second option, use an environment variable
# export CLOUDING_TOKEN=token123, it is not necessary to add the provider block
provider "clouding" {}
