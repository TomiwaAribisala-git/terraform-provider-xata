terraform {
  required_providers {
    xata = {
      source = "registry.terraform.io/hashicorp/xata"
    }
  }
}

provider "xata" {
  apikey = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
}

data "xata_workspaces" "example" {}

output "example_workspaces" {
  value = data.xata_workspaces.example
}
