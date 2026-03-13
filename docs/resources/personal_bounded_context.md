---
page_title: "fc_personal_bounded_context Resource - Fractal Cloud"
subcategory: ""
description: |-
  Manages a personal bounded context in Fractal Cloud.
---

# fc_personal_bounded_context (Resource)

Manages a personal bounded context scoped to a user account. A bounded context is a logical segmentation unit that groups fractals, live systems, and environments under a single governance scope.

## Example Usage

```terraform
resource "fc_personal_bounded_context" "production" {
  short_name   = "production"
  display_name = "Production"
  description  = "Production bounded context"
}
```

## Schema

### Required

- `short_name` (String) Short name identifier for the bounded context.

### Optional

- `display_name` (String) Human-readable display name.
- `description` (String) Description of the bounded context.
- `icon` (String) Icon identifier.

### Read-Only

- `id` (Object) Composite identifier containing `type`, `owner_id`, and `short_name`.
- `status` (String) Current status of the bounded context.
- `live_systems_ids` (List of String) IDs of live systems in this bounded context.
- `fractals_ids` (List of String) IDs of fractals in this bounded context.
- `created_at` (String) Creation timestamp.
- `updated_at` (String) Last update timestamp.
