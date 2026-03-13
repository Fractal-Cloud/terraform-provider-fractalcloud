---
page_title: "observability_saas_unmanaged Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates an unmanaged Observability component for external resources.
---

# function: observability_saas_unmanaged

Creates an unmanaged Observability component for external resources. Use this to reference observability infrastructure that is managed outside of Fractal Cloud, allowing other components to declare dependencies or links to it.

## Example Usage

```terraform
locals {
  external_monitoring = provider::fc::observability_saas_unmanaged({
    id           = "external-datadog"
    display_name = "External Datadog"
    description  = "Datadog instance managed outside Fractal Cloud"
  })
}
```

## Signature

```text
observability_saas_unmanaged(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
