resource "fractalcloud_resource_group" "production_audi" {
  type = "Personal"
  owner_id = "xxx"
  display_name = "Audi Production"
}

resource "fractalcloud_resource_group" "production_toyota" {
  type = "Personal"
  owner_id = "xxx"
  display_name = "Toyota Production"
}

resource "fractalcloud_resource_group" "production" {
  type = "Personal"
  owner_id = "xxx"
  display_name = "Production"
}
