resource "fractalcloud_personal_resource_group" "new_rg" {
  short_name = "new-rg"
  display_name = "New Resource Group"
  description = "Resource group created through Terraform provider"
}

data "fractalcloud_personal_resource_group" "existing_resource_group" {
  short_name = "existing-resource-group"
}

output sandbox_rg {
  value = data.fractalcloud_personal_resource_group.existing_resource_group
}