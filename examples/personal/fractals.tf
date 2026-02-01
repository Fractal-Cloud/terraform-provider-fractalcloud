data "fractalcloud_fractal" "existing_fractal" {
  resource_group_id = data.fractalcloud_personal_resource_group.existing_resource_group.id
  name = "existing-fractal"
  version = "1.0"
}

output fractal {
  value = data.fractalcloud_fractal.existing_fractal
}