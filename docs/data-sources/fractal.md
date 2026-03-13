---
page_title: "fc_fractal Data Source - Fractal Cloud"
subcategory: ""
description: |-
  Looks up an existing fractal definition in Fractal Cloud.
---

# fc_fractal (Data Source)

Use this data source to look up an existing fractal definition by bounded context, name, and version.

## Example Usage

```terraform
data "fc_fractal" "existing" {
  bounded_context_id = data.fc_personal_bounded_context.existing.id
  name    = "microservice-template"
  version = "1.0"
}
```

## Schema

### Required

- `bounded_context_id` (Object) Composite bounded context identifier containing `type`, `owner_id`, and `short_name`.
- `name` (String) Fractal name.
- `version` (String) Fractal version.

### Read-Only

- `description` (String) Description.
- `is_private` (Boolean) Whether the fractal is private.
- `components` (List of Object) Blueprint components, each containing `id`, `type`, `display_name`, `description`, `parameters`, `dependencies_ids`, and `links`.
- `created_at` (String) Creation timestamp.
