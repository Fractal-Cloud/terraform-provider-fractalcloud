---
page_title: "storage_paas_relational_dbms Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a Relational DBMS Platform blueprint component.
---

# function: storage_paas_relational_dbms

Creates a Relational DBMS (Database Management System) Platform blueprint component. This represents a managed database server instance (e.g. RDS, Azure Database, Cloud SQL) on which individual databases are hosted.

## Example Usage

```terraform
locals {
  dbms = provider::fc::storage_paas_relational_dbms({
    id             = "postgres-server"
    display_name   = "PostgreSQL Server"
    engine_version = "16"
  })
}
```

## Signature

```text
storage_paas_relational_dbms(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `engine_version` | String | No | Database engine version (e.g. `"16"`, `"8.0"`). |
