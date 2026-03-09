data "fractalcloud_fractal" "existing_fractal" {
  bounded_context_id = data.fractalcloud_personal_bounded_context.existing_bounded_context.id
  name    = "existing-fractal"
  version = "1.0"
}

# Fractal using provider functions — no magic strings needed
resource "fractalcloud_fractal" "iaas_fractal" {
  bounded_context_id = data.fractalcloud_personal_bounded_context.existing_bounded_context.id
  name        = "basic-iaas"
  version     = "1.0"
  description = "Basic IaaS Fractal with VPC, Subnet, and VM"

  components = [
    provider::fractalcloud::virtual_network({
      id           = "main-vpc"
      display_name = "Main VPC"
      description  = "Primary VPC for the IaaS workload"
      cidr_block   = "10.0.0.0/16"
    }),
    provider::fractalcloud::subnet({
      id                = "public-subnet"
      display_name      = "Public Subnet"
      cidr_block        = "10.0.1.0/24"
      availability_zone = "eu-central-1a"
      vpc_id            = "main-vpc"
    }),
    provider::fractalcloud::security_group({
      id           = "web-sg"
      display_name = "Web Security Group"
      description  = "Allow HTTPS traffic"
      vpc_id       = "main-vpc"
      ingress_rules = [
        {
          from_port   = 443
          source_cidr = "0.0.0.0/0"
        }
      ]
    }),
    provider::fractalcloud::virtual_machine({
      id           = "web-server"
      display_name = "Web Server"
      description  = "Application server"
      subnet_id    = "public-subnet"
    }),
  ]
}

# Fractal using provider functions for container workloads
resource "fractalcloud_fractal" "container_fractal" {
  bounded_context_id = data.fractalcloud_personal_bounded_context.existing_bounded_context.id
  name        = "microservice"
  version     = "1.0"
  description = "Containerized Microservice Fractal"

  components = [
    provider::fractalcloud::container_platform({
      id           = "k8s-cluster"
      display_name = "Kubernetes Cluster"
      description  = "Managed container orchestration platform"
    }),
    provider::fractalcloud::workload({
      id              = "api-service"
      display_name    = "API Service"
      description     = "Main API workload"
      container_image = "my-registry/api-service:latest"
      container_port  = 8080
      cpu             = "512"
      memory          = "1024"
      desired_count   = 2
      platform_id     = "k8s-cluster"
    }),
  ]
}

output "fractal" {
  value = data.fractalcloud_fractal.existing_fractal
}
