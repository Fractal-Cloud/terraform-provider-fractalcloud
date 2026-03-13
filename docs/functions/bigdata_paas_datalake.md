---
page_title: "bigdata_paas_datalake Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a Data Lake component.
---

# function: bigdata_paas_datalake

Creates a Data Lake component. This represents a managed data lake storage layer for big data workloads.

## Example Usage

```terraform
locals {
  datalake = provider::fc::bigdata_paas_datalake({
    id           = "main-datalake"
    display_name = "Main Data Lake"
    description  = "Central data lake for analytics workloads"
  })
}
```

## Signature

```text
bigdata_paas_datalake(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
