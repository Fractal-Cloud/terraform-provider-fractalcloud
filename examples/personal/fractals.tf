data "fractalcloud_fractal" "existing_fractal" {
  bounded_context_id = data.fractalcloud_personal_bounded_context.existing_bounded_context.id
  name = "existing-fractal"
  version = "1.0"
}

resource "fractalcloud_fractal" "new_fractal" {
  bounded_context_id = data.fractalcloud_personal_bounded_context.existing_bounded_context.id
  name = "new-fractal"
  version = "1.0"
  description = "Fractal Created with Terraform"
  components = [
    {
      "dependencies_ids" = []
      "description" = "Container Platform component"
      "display_name" = "Container Platform"
      "id" = "container-platform-1"
      "links" = []
      "output_fields" = []
      "parameters" = {}
      "recreate_on_failure" = false
      "type" = "NetworkAndCompute.PaaS.ContainerPlatform"
      "version" = "1.0"
    },
    {
      "dependencies_ids" = tolist([
        "container-platform-1",
      ])
      "description" = "Containerized Search Platform component"
      "display_name" = "Containerized Search Platform"
      "id" = "search-cluster-1"
      "is_locked" = false
      "links" = []
      "output_fields" = []
      "parameters" = {}
      "recreate_on_failure" = false
      "type" = "Storage.CaaS.Search"
      "version" = "1.0"
    }
  ]
}

output fractal {
  value = data.fractalcloud_fractal.existing_fractal
}
