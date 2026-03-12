data "fractalcloud_fractal" "existing_fractal" {
  bounded_context_id = data.fractalcloud_personal_bounded_context.existing_bounded_context.id
  name    = "existing-fractal"
  version = "1.0"
}

# IaaS Fractal with VPC, subnets, security groups, and VMs linked together.
# Components are defined as locals so that dependencies and links use
# direct object references (type-checked) instead of string IDs.
locals {
  main_vpc = provider::fractalcloud::network_and_compute_iaas_virtual_network({
    id           = "main-vpc"
    display_name = "Main VPC"
    cidr_block   = "10.0.0.0/16"
  })

  public_subnet = provider::fractalcloud::network_and_compute_iaas_subnet({
    id                = "public-subnet"
    display_name      = "Public Subnet"
    cidr_block        = "10.0.1.0/24"
    availability_zone = "eu-central-1a"
    vpc               = local.main_vpc # type-checked dependency on VPC
  })

  web_sg = provider::fractalcloud::network_and_compute_iaas_security_group({
    id           = "web-sg"
    display_name = "Web Security Group"
    description  = "Allow HTTPS from the internet"
    vpc          = local.main_vpc # type-checked dependency on VPC
    ingress_rules = [
      {
        from_port   = 443
        source_cidr = "0.0.0.0/0"
      }
    ]
  })

  api_server = provider::fractalcloud::network_and_compute_iaas_virtual_machine({
    id              = "api-server"
    display_name    = "API Server"
    subnet          = local.public_subnet
    security_groups = [local.web_sg]
  })

  web_server = provider::fractalcloud::network_and_compute_iaas_virtual_machine({
    id              = "web-server"
    display_name    = "Web Server"
    subnet          = local.public_subnet   # type-checked dependency on subnet
    security_groups = [local.web_sg]        # type-checked SG membership link
    links = [
      {
        target    = local.api_server        # type-checked port-based traffic rule
        from_port = 8080
      }
    ]
  })
}

resource "fractalcloud_fractal" "iaas_fractal" {
  bounded_context_id = data.fractalcloud_personal_bounded_context.existing_bounded_context.id
  name        = "basic-iaas"
  version     = "1.0"
  description = "IaaS Fractal with network, security, and compute"

  components = [
    local.main_vpc,
    local.public_subnet,
    local.web_sg,
    local.web_server,
    local.api_server,
  ]
}

# Container Fractal with workloads linked for inter-service traffic
locals {
  k8s_cluster = provider::fractalcloud::network_and_compute_paas_container_platform({
    id           = "k8s-cluster"
    display_name = "Kubernetes Cluster"
  })

  db_service = provider::fractalcloud::custom_workloads_caas_workload({
    id              = "db-service"
    display_name    = "Database Service"
    container_image = "postgres:16"
    container_port  = 5432
    cpu             = "1024"
    memory          = "2048"
    platform        = local.k8s_cluster
  })

  api_service = provider::fractalcloud::custom_workloads_caas_workload({
    id              = "api-service"
    display_name    = "API Service"
    container_image = "my-registry/api-service:latest"
    container_port  = 8080
    cpu             = "512"
    memory          = "1024"
    desired_count   = 2
    platform        = local.k8s_cluster # type-checked dependency on platform
    links = [
      {
        target    = local.db_service    # type-checked traffic rule to database
        from_port = 5432
        protocol  = "tcp"
      }
    ]
  })
}

resource "fractalcloud_fractal" "container_fractal" {
  bounded_context_id = data.fractalcloud_personal_bounded_context.existing_bounded_context.id
  name        = "microservice"
  version     = "1.0"
  description = "Containerized Microservice Fractal"

  components = [
    local.k8s_cluster,
    local.api_service,
    local.db_service,
  ]
}

output "fractal" {
  value = data.fractalcloud_fractal.existing_fractal
}
