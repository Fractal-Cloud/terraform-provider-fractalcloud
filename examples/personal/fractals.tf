data "fractalcloud_fractal" "existing_fractal" {
  bounded_context_id = data.fractalcloud_personal_bounded_context.existing_bounded_context.id
  name    = "existing-fractal"
  version = "1.0"
}

# IaaS Fractal with VPC, subnets, security groups, and VMs linked together
resource "fractalcloud_fractal" "iaas_fractal" {
  bounded_context_id = data.fractalcloud_personal_bounded_context.existing_bounded_context.id
  name        = "basic-iaas"
  version     = "1.0"
  description = "IaaS Fractal with network, security, and compute"

  components = [
    provider::fractalcloud::virtual_network({
      id           = "main-vpc"
      display_name = "Main VPC"
      cidr_block   = "10.0.0.0/16"
    }),

    provider::fractalcloud::subnet({
      id                = "public-subnet"
      display_name      = "Public Subnet"
      cidr_block        = "10.0.1.0/24"
      availability_zone = "eu-central-1a"
      vpc_id            = "main-vpc" # auto-wires dependency on VPC
    }),

    provider::fractalcloud::security_group({
      id           = "web-sg"
      display_name = "Web Security Group"
      description  = "Allow HTTPS from the internet"
      vpc_id       = "main-vpc" # auto-wires dependency on VPC
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
      subnet_id    = "public-subnet" # auto-wires dependency on subnet
      security_groups = ["web-sg"]   # link: SG membership (no settings)
      links = [
        {
          target_id = "api-server"   # link: port-based traffic rule
          from_port = 8080
        }
      ]
    }),

    provider::fractalcloud::virtual_machine({
      id           = "api-server"
      display_name = "API Server"
      subnet_id    = "public-subnet"
      security_groups = ["web-sg"]
    }),
  ]
}

# Container Fractal with workloads linked for inter-service traffic
resource "fractalcloud_fractal" "container_fractal" {
  bounded_context_id = data.fractalcloud_personal_bounded_context.existing_bounded_context.id
  name        = "microservice"
  version     = "1.0"
  description = "Containerized Microservice Fractal"

  components = [
    provider::fractalcloud::container_platform({
      id           = "k8s-cluster"
      display_name = "Kubernetes Cluster"
    }),

    provider::fractalcloud::workload({
      id              = "api-service"
      display_name    = "API Service"
      container_image = "my-registry/api-service:latest"
      container_port  = 8080
      cpu             = "512"
      memory          = "1024"
      desired_count   = 2
      platform_id     = "k8s-cluster" # auto-wires dependency on platform
      links = [
        {
          target_id = "db-service"    # link: traffic rule to database
          from_port = 5432
          protocol  = "tcp"
        }
      ]
    }),

    provider::fractalcloud::workload({
      id              = "db-service"
      display_name    = "Database Service"
      container_image = "postgres:16"
      container_port  = 5432
      cpu             = "1024"
      memory          = "2048"
      platform_id     = "k8s-cluster"
    }),
  ]
}

output "fractal" {
  value = data.fractalcloud_fractal.existing_fractal
}
