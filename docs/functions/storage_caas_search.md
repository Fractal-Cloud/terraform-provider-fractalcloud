---
page_title: "storage_caas_search Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a containerized Search Platform component.
---

# function: storage_caas_search

Creates a containerized Search Platform component. If `container_platform` is provided, it is added as a dependency to ensure the container platform is provisioned before the search component.

## Example Usage

```terraform
locals {
  k8s = provider::fc::network_and_compute_paas_kubernetes({
    id = "k8s-cluster"
  })

  search = provider::fc::storage_caas_search({
    id                 = "search-platform"
    display_name       = "Search Engine"
    container_platform = local.k8s
  })
}
```

## Signature

```text
storage_caas_search(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `container_platform` | Component Object | No | The container platform component to depend on. If provided, added as a dependency. |
