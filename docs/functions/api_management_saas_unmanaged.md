---
page_title: "api_management_saas_unmanaged Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates an unmanaged API Management component for external resources.
---

# function: api_management_saas_unmanaged

Creates an unmanaged API Management component for external resources. Use this to reference API management infrastructure that is managed outside of Fractal Cloud, allowing other components to declare dependencies or links to it.

## Example Usage

```terraform
locals {
  external_apigw = provider::fc::api_management_saas_unmanaged({
    id           = "external-api-gateway"
    display_name = "External API Gateway"
    description  = "API gateway managed outside Fractal Cloud"
  })
}
```

## Signature

```text
api_management_saas_unmanaged(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
