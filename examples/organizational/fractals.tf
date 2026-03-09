data "fractalcloud_fractal" "existing_fractal" {
  bounded_context_id = data.fractalcloud_organizational_bounded_context.existing_bounded_context.id
  name    = "existing-fractal"
  version = "1.0"
}

resource "fractalcloud_fractal" "new_fractal" {
  bounded_context_id = data.fractalcloud_organizational_bounded_context.existing_bounded_context.id
  name        = "new-fractal"
  version     = "1.0"
  description = "Fractal Created with Terraform in Organization"
  components = [
    {
      id           = "container-platform-1"
      type         = "NetworkAndCompute.PaaS.ContainerPlatform"
      display_name = "Container Platform"
      description  = "Container Platform component"
      version      = "1.0"
      parameters   = {}
    },
    {
      id               = "search-cluster-1"
      type             = "Storage.CaaS.Search"
      display_name     = "Containerized Search Platform"
      description      = "Containerized Search Platform component"
      version          = "1.0"
      parameters       = {}
      dependencies_ids = ["container-platform-1"]
    }
  ]
}

output "fractal" {
  value = data.fractalcloud_fractal.existing_fractal
}
