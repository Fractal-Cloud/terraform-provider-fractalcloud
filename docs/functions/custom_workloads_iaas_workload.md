---
page_title: "custom_workloads_iaas_workload Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates an IaaS Workload blueprint component running on VMs.
---

# function: custom_workloads_iaas_workload

Creates an IaaS Workload blueprint component running on virtual machines. The `vm` dependency is type-validated to ensure it is a VirtualMachine component, and `subnet` is validated to ensure it is a Subnet component. Use `links` for port-based traffic rules and `security_groups` for SG membership.

## Example Usage

```terraform
locals {
  vm = provider::fc::network_and_compute_iaas_virtual_machine({
    id     = "app-server"
    subnet = local.subnet
  })

  workload = provider::fc::custom_workloads_iaas_workload({
    id              = "app-workload"
    display_name    = "Application Workload"
    container_image = "myregistry/app:latest"
    container_port  = 3000
    container_name  = "app"
    cpu             = "512"
    memory          = "1024"
    desired_count   = 2
    vm              = local.vm
    subnet          = local.subnet
  })
}
```

## Signature

```text
custom_workloads_iaas_workload(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `container_image` | String | No | Container image to deploy. |
| `container_port` | Number | No | Port the container listens on. |
| `container_name` | String | No | Name for the container. |
| `cpu` | String | No | CPU allocation. |
| `memory` | String | No | Memory allocation in MB. |
| `desired_count` | Number | No | Desired number of running instances. |
| `vm` | Component Object | No | A VirtualMachine component to add as a dependency. Must be a component returned by `network_and_compute_iaas_virtual_machine`. |
| `subnet` | Component Object | No | A Subnet component to add as a dependency. Must be a component returned by `network_and_compute_iaas_subnet`. |
| `links` | List of Object | No | Runtime relationship links to other components. Each link has a `target` (component object) and optional `settings` (map of string key-value pairs). |
| `security_groups` | List of Component Object | No | SecurityGroup components for SG membership. |
