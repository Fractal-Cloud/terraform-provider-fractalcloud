---
page_title: "custom_workloads_caas_workload Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a CaaS Workload (container) blueprint component.
---

# function: custom_workloads_caas_workload

Creates a CaaS Workload (Container as a Service) blueprint component. The `platform` dependency is type-validated to ensure it is a ContainerPlatform component, and `subnet` is validated to ensure it is a Subnet component. Use `links` for port-based traffic rules and `security_groups` for SG membership.

## Example Usage

```terraform
locals {
  k8s = provider::fc::network_and_compute_paas_container_platform({
    id = "k8s-cluster"
  })

  subnet = provider::fc::network_and_compute_iaas_subnet({
    id         = "app-subnet"
    cidr_block = "10.0.1.0/24"
    vpc        = local.vpc
  })

  workload = provider::fc::custom_workloads_caas_workload({
    id              = "api-service"
    display_name    = "API Service"
    container_image = "myregistry/api:latest"
    container_port  = 8080
    container_name  = "api"
    cpu             = "256"
    memory          = "512"
    desired_count   = 3
    platform        = local.k8s
    subnet          = local.subnet
  })
}
```

## Signature

```text
custom_workloads_caas_workload(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `container_image` | String | No | Container image to deploy (e.g. `"myregistry/api:latest"`). |
| `container_port` | Number | No | Port the container listens on. |
| `container_name` | String | No | Name for the container. |
| `cpu` | String | No | CPU allocation (e.g. `"256"` for 0.25 vCPU). |
| `memory` | String | No | Memory allocation in MB (e.g. `"512"`). |
| `desired_count` | Number | No | Desired number of running instances. |
| `platform` | Component Object | No | A ContainerPlatform component to add as a dependency. Must be a component returned by `network_and_compute_paas_container_platform`. |
| `subnet` | Component Object | No | A Subnet component to add as a dependency. Must be a component returned by `network_and_compute_iaas_subnet`. |
| `links` | List of Object | No | Runtime relationship links to other components. Each link has a `target` (component object) and optional `settings` (map of string key-value pairs). |
| `security_groups` | List of Component Object | No | SecurityGroup components for SG membership. |
