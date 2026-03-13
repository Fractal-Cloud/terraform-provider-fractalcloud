---
page_title: "messaging_caas_entity Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a CaaS Message Entity component.
---

# function: messaging_caas_entity

Creates a CaaS Message Entity component. If `broker` is provided, it is added as a dependency to ensure the broker is provisioned before the entity.

## Example Usage

```terraform
locals {
  broker = provider::fc::messaging_caas_broker({
    id = "caas-broker"
  })

  entity = provider::fc::messaging_caas_entity({
    id           = "caas-topic"
    display_name = "Order Events"
    broker       = local.broker
  })
}
```

## Signature

```text
messaging_caas_entity(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `broker` | Component Object | No | The CaaS Message Broker component to depend on. If provided, added as a dependency. |
