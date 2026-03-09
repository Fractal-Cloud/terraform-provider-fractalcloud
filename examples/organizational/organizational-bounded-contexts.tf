data "fractalcloud_organization" "existing_org" {
  id = "a15f2627-c927-4125-a3a5-7141977135b1"
}

resource "fractalcloud_organizational_bounded_context" "new_bc" {
  short_name = "new-bc"
  organization_id = data.fractalcloud_organization.existing_org.id
  display_name = "New Bounded Context"
  description = "Bounded Context created through Terraform provider"
}

data "fractalcloud_organizational_bounded_context" "existing_bounded_context" {
  short_name = "existing-resource-group"
  organization_id = data.fractalcloud_organization.existing_org.id
}

output "production_org" {
  value = data.fractalcloud_organization.existing_org
}

output "production_bc" {
  value = data.fractalcloud_organizational_bounded_context.existing_bounded_context
}
