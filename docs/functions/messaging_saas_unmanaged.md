---
page_title: "messaging_saas_unmanaged Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates an unmanaged Messaging component for external resources.
---

# function: messaging_saas_unmanaged

Creates an unmanaged Messaging component for external resources. Use this to represent messaging systems that are managed outside of Fractal Cloud but need to be referenced by other blueprint components.

## Example Usage

```terraform
locals {
  external_mq = provider::fc::messaging_saas_unmanaged({
    id           = "external-mq"
    display_name = "Legacy Message Queue"
    description  = "Externally managed message queue"
  })
}
```

## Signature

```text
messaging_saas_unmanaged(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
