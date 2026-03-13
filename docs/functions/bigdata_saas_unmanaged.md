---
page_title: "bigdata_saas_unmanaged Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates an unmanaged BigData component for external resources.
---

# function: bigdata_saas_unmanaged

Creates an unmanaged BigData component for external resources. Use this to reference big data infrastructure that is managed outside of Fractal Cloud, allowing other components to declare dependencies or links to it.

## Example Usage

```terraform
locals {
  external_spark = provider::fc::bigdata_saas_unmanaged({
    id           = "external-spark"
    display_name = "External Spark Cluster"
    description  = "Spark cluster managed outside Fractal Cloud"
  })
}
```

## Signature

```text
bigdata_saas_unmanaged(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
