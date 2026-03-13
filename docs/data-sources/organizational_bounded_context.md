---
page_title: "fc_organizational_bounded_context Data Source - Fractal Cloud"
subcategory: ""
description: |-
  Looks up an existing organizational bounded context in Fractal Cloud.
---

# fc_organizational_bounded_context (Data Source)

Use this data source to look up an existing organizational bounded context by its short name and organization.

## Example Usage

```terraform
data "fc_organizational_bounded_context" "existing" {
  short_name      = "platform"
  organization_id = data.fc_organization.my_org.id
}
```

## Schema

### Required

- `short_name` (String) Short name of the bounded context to look up.
- `organization_id` (String) ID of the organization.

### Read-Only

- `id` (Object) Composite identifier containing `type`, `owner_id`, and `short_name`.
- `display_name` (String) Human-readable display name.
- `description` (String) Description.
- `status` (String) Current status.
- `icon` (String) Icon identifier.
- `members_ids` (List of String) IDs of members.
- `teams_ids` (List of String) IDs of teams.
- `managers_ids` (List of String) IDs of managers.
- `live_systems_ids` (List of String) IDs of live systems.
- `fractals_ids` (List of String) IDs of fractals.
- `created_at` (String) Creation timestamp.
- `created_by` (String) Creator.
- `updated_at` (String) Last update timestamp.
- `updated_by` (String) Last updater.
