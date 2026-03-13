---
page_title: "fc_organization Data Source - Fractal Cloud"
subcategory: ""
description: |-
  Looks up an existing organization in Fractal Cloud.
---

# fc_organization (Data Source)

Use this data source to look up an existing organization by ID.

## Example Usage

```terraform
data "fc_organization" "my_org" {
  id = "your-organization-id"
}
```

## Schema

### Required

- `id` (String) Organization ID.

### Read-Only

- `display_name` (String) Display name.
- `description` (String) Description.
- `icon` (String) Icon identifier.
- `tags` (List of String) Tags.
- `social_links` (List of String) Social links.
- `admins` (List of String) Admin user IDs.
- `members` (List of String) Member user IDs.
- `teams` (List of String) Team IDs.
- `bounded_contexts` (List of String) Bounded context IDs.
- `status` (String) Current status.
- `subscription_id` (String) Subscription ID.
- `created_at` (String) Creation timestamp.
- `created_by` (String) Creator.
- `updated_at` (String) Last update timestamp.
- `updated_by` (String) Last updater.
