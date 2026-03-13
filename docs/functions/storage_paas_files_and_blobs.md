---
page_title: "storage_paas_files_and_blobs Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a File & Blob Storage blueprint component.
---

# function: storage_paas_files_and_blobs

Creates a File & Blob Storage blueprint component for managed object storage services such as S3, Azure Blob Storage, or Google Cloud Storage.

## Example Usage

```terraform
locals {
  storage = provider::fc::storage_paas_files_and_blobs({
    id           = "app-storage"
    display_name = "Application Storage"
    description  = "Object storage for application assets"
  })
}
```

## Signature

```text
storage_paas_files_and_blobs(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
