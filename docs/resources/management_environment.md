---
page_title: "fc_management_environment Resource - Fractal Cloud"
subcategory: ""
description: |-
  Manages a management environment in Fractal Cloud.
---

# fc_management_environment (Resource)

Manages a management environment for governance. A management environment defines the cloud provider agent configuration and bounded context assignments for infrastructure management.

~> **Note** This resource is currently a schema placeholder. CRUD operations are not yet implemented.

## Example Usage

```terraform
resource "fc_management_environment" "production" {
  type         = "Production"
  owner_id     = "your-owner-id"
  display_name = "Production Environment"

  bounded_contexts = [
    fc_personal_bounded_context.production.id,
  ]

  aws_agent = {
    region          = "eu-central-1"
    organization_id = "123456789012"
    account_id      = "123456789012"
  }
}
```

## Schema

### Required

- `type` (String) Environment type.
- `owner_id` (String) Owner identifier.
- `display_name` (String) Human-readable display name.
- `bounded_contexts` (List of Object) Bounded contexts assigned to this environment.

### Optional

- `parameters` (Map of Object) Environment parameters.
- `aws_agent` (Object) AWS agent configuration with `region`, `organization_id`, `account_id`.
- `azure_agent` (Object) Azure agent configuration with `region`, `tenant_id`, `subscription_id`.
- `gcp_agent` (Object) GCP agent configuration with `region`, `organization_id`, `project_id`.
- `oci_agent` (Object) OCI agent configuration with `region`, `tenancy_id`, `compartment_id`.
- `hetzner_agent` (Object) Hetzner agent configuration with `region`, `project_id`.
- `default_cicd_profile_short_name` (String) Default CI/CD profile short name.

### Read-Only

- `id` (Object) Composite identifier.
- `status` (String) Current status.
- `created_at` (String) Creation timestamp.
- `created_by` (String) Creator.
- `updated_at` (String) Last update timestamp.
- `updated_by` (String) Last updater.
