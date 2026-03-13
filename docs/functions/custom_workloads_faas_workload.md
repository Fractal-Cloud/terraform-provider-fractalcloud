---
page_title: "custom_workloads_faas_workload Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a FaaS (serverless) Workload blueprint component.
---

# function: custom_workloads_faas_workload

Creates a FaaS (Function as a Service) Workload blueprint component for serverless deployments. Includes additional serverless-specific parameters such as `runtime`, `handler`, `memory_mb`, and `timeout_seconds`. The `platform` dependency is type-validated to ensure it is a ContainerPlatform component, and `subnet` is validated to ensure it is a Subnet component.

## Example Usage

```terraform
locals {
  workload = provider::fc::custom_workloads_faas_workload({
    id              = "event-handler"
    display_name    = "Event Handler Function"
    container_image = "myregistry/handler:latest"
    container_port  = 8080
    container_name  = "handler"
    cpu             = "256"
    memory          = "512"
    desired_count   = 1
    runtime         = "nodejs20.x"
    handler         = "index.handler"
    memory_mb       = 256
    timeout_seconds = 30
    platform        = local.k8s
    subnet          = local.subnet
  })
}
```

## Signature

```text
custom_workloads_faas_workload(config object) object
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
| `runtime` | String | No | Serverless runtime (e.g. `"nodejs20.x"`, `"python3.12"`, `"java21"`). |
| `handler` | String | No | Function entry point (e.g. `"index.handler"`). |
| `memory_mb` | Number | No | Memory allocation for the function in MB. |
| `timeout_seconds` | Number | No | Function execution timeout in seconds. |
| `platform` | Component Object | No | A ContainerPlatform component to add as a dependency. Must be a component returned by `network_and_compute_paas_container_platform`. |
| `subnet` | Component Object | No | A Subnet component to add as a dependency. Must be a component returned by `network_and_compute_iaas_subnet`. |
| `links` | List of Object | No | Runtime relationship links to other components. Each link has a `target` (component object) and optional `settings` (map of string key-value pairs). |
| `security_groups` | List of Component Object | No | SecurityGroup components for SG membership. |
