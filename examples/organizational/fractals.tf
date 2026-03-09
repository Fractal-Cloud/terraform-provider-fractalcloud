data "fractalcloud_fractal" "existing_fractal" {
  bounded_context_id = data.fractalcloud_organizational_bounded_context.existing_bounded_context.id
  name    = "existing-fractal"
  version = "1.0"
}

# IaaS Fractal using provider functions
resource "fractalcloud_fractal" "org_iaas" {
  bounded_context_id = data.fractalcloud_organizational_bounded_context.existing_bounded_context.id
  name        = "org-iaas"
  version     = "1.0"
  description = "Organization IaaS Fractal"

  components = [
    provider::fractalcloud::virtual_network({
      id           = "main-vpc"
      display_name = "Main VPC"
      cidr_block   = "10.0.0.0/16"
    }),
    provider::fractalcloud::subnet({
      id                = "app-subnet"
      display_name      = "Application Subnet"
      cidr_block        = "10.0.1.0/24"
      availability_zone = "eu-central-1a"
      vpc_id            = "main-vpc"
    }),
    provider::fractalcloud::subnet({
      id                = "data-subnet"
      display_name      = "Data Subnet"
      cidr_block        = "10.0.2.0/24"
      availability_zone = "eu-central-1b"
      vpc_id            = "main-vpc"
    }),
    provider::fractalcloud::security_group({
      id           = "app-sg"
      display_name = "Application Security Group"
      description  = "Allow application traffic"
      vpc_id       = "main-vpc"
      ingress_rules = [
        {
          from_port   = 443
          source_cidr = "0.0.0.0/0"
        },
        {
          from_port           = 8080
          source_component_id = "web-server"
        }
      ]
    }),
    provider::fractalcloud::virtual_machine({
      id           = "web-server"
      display_name = "Web Server"
      subnet_id    = "app-subnet"
    }),
  ]
}

output "fractal" {
  value = data.fractalcloud_fractal.existing_fractal
}
