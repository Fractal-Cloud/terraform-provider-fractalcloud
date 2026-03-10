data "fractalcloud_fractal" "existing_fractal" {
  bounded_context_id = data.fractalcloud_organizational_bounded_context.existing_bounded_context.id
  name    = "existing-fractal"
  version = "1.0"
}

# Multi-tier IaaS Fractal with full dependency and link wiring
resource "fractalcloud_fractal" "org_iaas" {
  bounded_context_id = data.fractalcloud_organizational_bounded_context.existing_bounded_context.id
  name        = "org-iaas"
  version     = "1.0"
  description = "Organization multi-tier IaaS architecture"

  components = [
    provider::fractalcloud::virtual_network({
      id           = "main-vpc"
      display_name = "Main VPC"
      cidr_block   = "10.0.0.0/16"
    }),

    provider::fractalcloud::subnet({
      id                = "web-subnet"
      display_name      = "Web Tier Subnet"
      cidr_block        = "10.0.1.0/24"
      availability_zone = "eu-central-1a"
      vpc_id            = "main-vpc"
    }),

    provider::fractalcloud::subnet({
      id                = "app-subnet"
      display_name      = "App Tier Subnet"
      cidr_block        = "10.0.2.0/24"
      availability_zone = "eu-central-1b"
      vpc_id            = "main-vpc"
    }),

    provider::fractalcloud::security_group({
      id          = "web-sg"
      description = "Allow HTTPS from internet"
      vpc_id      = "main-vpc"
      ingress_rules = [
        {
          from_port   = 443
          source_cidr = "0.0.0.0/0"
        }
      ]
    }),

    provider::fractalcloud::security_group({
      id          = "app-sg"
      description = "Allow traffic from web tier only"
      vpc_id      = "main-vpc"
      ingress_rules = [
        {
          from_port           = 8080
          source_component_id = "web-server"
        }
      ]
    }),

    # Web server: depends on web-subnet, member of web-sg, links to app-server on port 8080
    provider::fractalcloud::virtual_machine({
      id              = "web-server"
      display_name    = "Web Server"
      subnet_id       = "web-subnet"
      security_groups = ["web-sg"]
      links = [
        {
          target_id = "app-server"
          from_port = 8080
        }
      ]
    }),

    # App server: depends on app-subnet, member of app-sg
    provider::fractalcloud::virtual_machine({
      id              = "app-server"
      display_name    = "Application Server"
      subnet_id       = "app-subnet"
      security_groups = ["app-sg"]
    }),
  ]
}

output "fractal" {
  value = data.fractalcloud_fractal.existing_fractal
}
