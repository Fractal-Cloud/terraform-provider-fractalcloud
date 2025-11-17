terraform {
  required_providers {
    fractalcloud = {
      source = "registry.terraform.io/fractalcloud/fc"
    }
  }
  required_version = ">= 1.1.0"
}

provider "fractalcloud" {
  service_account_id = "xxx"
  service_account_secret = "xxx"
}
