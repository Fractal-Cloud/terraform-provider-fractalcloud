---
page_title: "network_and_compute_saas_unmanaged Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates an unmanaged NetworkAndCompute component for external resources.
---

# function: network_and_compute_saas_unmanaged

Creates an unmanaged NetworkAndCompute component for external resources. Use this to represent infrastructure that exists outside of Fractal Cloud management but needs to be referenced as a dependency or link target by other blueprint components.

## Example Usage

```terraform
locals {
  external_network = provider::fc::network_and_compute_saas_unmanaged({
    id           = "external-vpc"
    display_name = "External VPC"
    description  = "Pre-existing VPC managed outside Fractal Cloud"
  })
}
```

## Signature

```text
network_and_compute_saas_unmanaged(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
