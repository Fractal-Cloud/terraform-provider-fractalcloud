---
page_title: "fc_personal_bounded_context Data Source - Fractal Cloud"
subcategory: ""
description: |-
  Looks up an existing personal bounded context in Fractal Cloud.
---

# fc_personal_bounded_context (Data Source)

Use this data source to look up an existing personal bounded context by its short name.

## Example Usage

```terraform
data "fc_personal_bounded_context" "existing" {
  short_name = "production"
}
```

## Schema

### Required

- `short_name` (String) Short name of the bounded context to look up.

### Read-Only

- `id` (Object) Composite identifier containing `type`, `owner_id`, and `short_name`.
- `display_name` (String) Human-readable display name.
- `description` (String) Description.
- `status` (String) Current status.
- `icon` (String) Icon identifier.
- `live_systems_ids` (List of String) IDs of live systems.
- `fractals_ids` (List of String) IDs of fractals.
- `created_at` (String) Creation timestamp.
- `updated_at` (String) Last update timestamp.
