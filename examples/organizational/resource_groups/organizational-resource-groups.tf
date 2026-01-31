data "fractalcloud_organization" "existing_org" {
  id = "a15f2627-c927-4125-a3a5-7141977135b1"
}

resource "fractalcloud_organizational_resource_group" "new_rg" {
  short_name = "new-rg"
  organization_id = data.fractalcloud_organization.existing_org.id
  display_name = "New Resource Group"
  description = "Resource group created through Terraform provider"
}

data "fractalcloud_organizational_resource_group" "existing_resource_group" {
  short_name = "existing-resource-group"
  organization_id = data.fractalcloud_organization.existing_org.id
}

output "production_rg" {
  value = data.fractalcloud_organizational_resource_group.existing_resource_group
}