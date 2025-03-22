terraform {
  required_providers {
    xata = {
      source = "registry.terraform.io/hashicorp/xata"
    }
  }
  required_version = ">= 1.1.0"
}

provider "xata" {
  apikey = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
}

resource "xata_workspace" "markspace" {
  name = "markspace"
}