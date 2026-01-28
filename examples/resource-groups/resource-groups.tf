resource "fractalcloud_resource_group" "audi_production" {
  type = "Personal"
  owner_id = "xxx"
  display_name = "Audi production"
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