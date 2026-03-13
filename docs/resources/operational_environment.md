---
page_title: "fc_operational_environment Resource - Fractal Cloud"
subcategory: ""
description: |-
  Manages an operational environment in Fractal Cloud.
---

# fc_operational_environment (Resource)

Manages an operational environment for runtime governance. An operational environment is linked to a management environment and defines the runtime scope for live system deployments.

~> **Note** This resource is currently a schema placeholder. CRUD operations are not yet implemented.

## Example Usage

```terraform
resource "fc_operational_environment" "staging" {
  management_environment_id = fc_management_environment.production.id
  display_name              = "Staging"

  bounded_contexts = [
    fc_personal_bounded_context.production.id,
  ]
}
```

## Schema

### Required

- `management_environment_id` (Object) Composite identifier of the parent management environment.
- `display_name` (String) Human-readable display name.
- `bounded_contexts` (List of Object) Bounded contexts assigned to this environment.

### Optional

- `parameters` (Map of Object) Environment parameters.
- `agents` (Set of String) Agent identifiers.
- `default_cicd_profile_short_name` (String) Default CI/CD profile short name.

### Read-Only

- `id` (Object) Composite identifier.
- `status` (String) Current status.
- `created_at` (String) Creation timestamp.
- `created_by` (String) Creator.
- `updated_at` (String) Last update timestamp.
- `updated_by` (String) Last updater.
