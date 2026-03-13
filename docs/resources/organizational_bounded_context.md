---
page_title: "fc_organizational_bounded_context Resource - Fractal Cloud"
subcategory: ""
description: |-
  Manages an organizational bounded context in Fractal Cloud.
---

# fc_organizational_bounded_context (Resource)

Manages an organizational bounded context scoped to an organization. Supports member, team, and manager assignments.

## Example Usage

```terraform
data "fc_organization" "my_org" {
  id = "your-organization-id"
}

resource "fc_organizational_bounded_context" "platform" {
  short_name      = "platform"
  organization_id = data.fc_organization.my_org.id
  display_name    = "Platform"
  description     = "Shared platform bounded context"
}
```

## Schema

### Required

- `short_name` (String) Short name identifier for the bounded context.
- `organization_id` (String) ID of the organization that owns this bounded context.

### Optional

- `display_name` (String) Human-readable display name.
- `description` (String) Description of the bounded context.
- `icon` (String) Icon identifier.
- `members_ids` (List of String) IDs of members assigned to this bounded context.
- `teams_ids` (List of String) IDs of teams assigned to this bounded context.
- `managers_ids` (List of String) IDs of managers assigned to this bounded context.

### Read-Only

- `id` (Object) Composite identifier containing `type`, `owner_id`, and `short_name`.
- `status` (String) Current status of the bounded context.
- `live_systems_ids` (List of String) IDs of live systems in this bounded context.
- `fractals_ids` (List of String) IDs of fractals in this bounded context.
- `created_at` (String) Creation timestamp.
- `created_by` (String) User who created this bounded context.
- `updated_at` (String) Last update timestamp.
- `updated_by` (String) User who last updated this bounded context.
