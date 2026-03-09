resource "fractalcloud_personal_bounded_context" "new_bc" {
  short_name = "new-bc"
  display_name = "New Bounded Context"
  description = "Bounded Context created through Terraform provider"
}

data "fractalcloud_personal_bounded_context" "existing_bounded_context" {
  short_name = "existing-resource-group"
}

output sandbox_bc {
  value = data.fractalcloud_personal_bounded_context.existing_bounded_context
}
