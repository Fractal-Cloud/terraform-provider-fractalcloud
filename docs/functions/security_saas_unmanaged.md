---
page_title: "security_saas_unmanaged Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates an unmanaged Security component for external resources.
---

# function: security_saas_unmanaged

Creates an unmanaged Security component for external resources. Use this to reference security infrastructure that is managed outside of Fractal Cloud, allowing other components to declare dependencies or links to it.

## Example Usage

```terraform
locals {
  external_vault = provider::fc::security_saas_unmanaged({
    id           = "external-vault"
    display_name = "External HashiCorp Vault"
    description  = "Vault instance managed outside Fractal Cloud"
  })
}
```

## Signature

```text
security_saas_unmanaged(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
