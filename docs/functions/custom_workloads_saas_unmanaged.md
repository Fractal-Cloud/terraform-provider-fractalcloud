---
page_title: "custom_workloads_saas_unmanaged Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates an unmanaged CustomWorkloads component for external resources.
---

# function: custom_workloads_saas_unmanaged

Creates an unmanaged CustomWorkloads component for external resources. Use this to represent workloads that exist outside of Fractal Cloud management but need to be referenced as a dependency or link target by other blueprint components.

## Example Usage

```terraform
locals {
  external_service = provider::fc::custom_workloads_saas_unmanaged({
    id           = "legacy-api"
    display_name = "Legacy API Service"
    description  = "Pre-existing API managed outside Fractal Cloud"
  })
}
```

## Signature

```text
custom_workloads_saas_unmanaged(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
