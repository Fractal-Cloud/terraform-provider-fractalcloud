resource "fractalcloud_resource_group" "audi_tf_production" {
  id = {
    type = "Personal"
    owner_id = "xxx"
    shortname = "audi-tf-production"
  }
  display_name = "Audi TF Production"
  description = "Resource group created through Terraform provider"
}

data "fractalcloud_resource_group" "production" {
  id = {
    type = "Personal"
    owner_id = "xxx"
    shortname = "production"
  }
}

output "production_rg" {
  value = data.fractalcloud_resource_group.production
}