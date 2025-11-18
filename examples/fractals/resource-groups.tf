resource "fractalcloud_resource_group" "production_audi" {
  type = "Personal"
  owner_id = "29b195ed-ac8b-45bb-b8c5-5ee0fd542b11"
  display_name = "Audi Production"
}

resource "fractalcloud_resource_group" "production_toyota" {
  type = "Personal"
  owner_id = "29b195ed-ac8b-45bb-b8c5-5ee0fd542b11"
  display_name = "Toyota Production"
}

resource "fractalcloud_resource_group" "production" {
  type = "Personal"
  owner_id = "29b195ed-ac8b-45bb-b8c5-5ee0fd542b11"
  display_name = "Production"
}
