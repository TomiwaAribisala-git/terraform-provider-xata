terraform {
  required_providers {
    xata = {
      source = "registry.terraform.io/hashicorp/xata"
    }
  }
}

provider "xata" {
  apikey = "32xekt495b3nivt435f5f5f5"
}

data "xata_workspaces" "example" {}

output "example_workspaces" {
  value = data.xata_workspaces.example
}
