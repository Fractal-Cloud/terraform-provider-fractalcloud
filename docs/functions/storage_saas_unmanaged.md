---
page_title: "storage_saas_unmanaged Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates an unmanaged Storage component for external resources.
---

# function: storage_saas_unmanaged

Creates an unmanaged Storage component for external resources. Use this to represent storage systems that are managed outside of Fractal Cloud but need to be referenced by other blueprint components.

## Example Usage

```terraform
locals {
  external_db = provider::fc::storage_saas_unmanaged({
    id           = "external-db"
    display_name = "Legacy Database"
    description  = "Externally managed database"
  })
}
```

## Signature

```text
storage_saas_unmanaged(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
