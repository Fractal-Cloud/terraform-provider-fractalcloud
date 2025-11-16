# Terraform Provider for Fractal Cloud

The Terraform Provider for Fractal Cloud enables Operations and Platform teams to manage Fractal Cloud configuration and governance resources through Terraform. It focuses on the operational aspects of the platform and is not intended for application teams or for managing Live Systems.

Fractal Cloud is a platform engineering solution that delivers secure and compliant infrastructure across any cloud. It provides developers with ready-to-use building blocks and architecture templates while centralizing automation and governance for operations teams.

This provider is designed for infrastructure and operations teams responsible for maintaining organizational structure, governance, and Fractal publication workflows.

---

## Scope of the Provider

### Supported

The provider manages operational and structural entities within Fractal Cloud:

- Resource Groups  
- Management Environments  
- Operational Environments  
- Fractals (publishing, updating metadata, associating to environments)

These resources represent the governance layer maintained by Cloud Centers of Excellence or platform engineering teams.

### Not Included

The provider does **not** manage:

- Live Systems  
- Application-level deployments  
- Runtime operations  
- CI or pipeline-triggered infrastructure instantiation

Live Systems are created by development teams using the Fractal SDK or Fractal Cloud UI.

---

## Getting Started

### Requirements

- Terraform v1.4+  
- A Fractal Cloud account  
- An API token with platform-level permissions

### Installation

```hcl
terraform {
  required_providers {
    fractalcloud = {
      source  = "fractalcloud/fractalcloud"
      version = "~> 0.1.0"
    }
  }
}

provider "fractalcloud" {
  endpoint = var.fractal_api_endpoint
  token    = var.fractal_api_token
}
```

### Managing Resource Groups

Resource Groups align cloud resources to business capabilities or domains.

```hcl
resource "fractalcloud_resource_group" "payments" {
  name        = "payments"
  description = "Resource group for payment-related capabilities"
}
```

### Managing Environments
Management and Operational Environments

```hcl
resource "fractalcloud_environment" "prod_management" {
  name        = "management-prod"
  type        = "management"
  description = "Management environment for production governance"
}

resource "fractalcloud_environment" "prod_operational" {
  name              = "operational-prod"
  type              = "operational"
  resource_group_id = fractalcloud_resource_group.payments.id
}
```

Environments define controlled spaces where Fractals are published and governed.

### Publishing and Managing Fractals

Terraform can maintain Fractal metadata and link them to environments.

```hcl
resource "fractalcloud_fractal" "microservice" {
  name           = "microservice-template"
  version        = "1.0.0"
  resource_group = fractalcloud_resource_group.payments.id

  description = "Standard microservice architecture Blueprint"
}
```

This does not instantiate Live Systems.
It manages the catalog used by development teams.

### Data Sources
#### Retrieve a Fractal

```hcl
data "fractalcloud_fractal" "shared" {
  name = "microservice-template"
}
```

#### Retrieve Environment Information
```hcl
data "fractalcloud_environment" "prod" {
  name = "operational-prod"
}
```

### Authentication
```hcl
provider "fractalcloud" {
  token = var.fractal_api_token
}
```

### Versioning

The provider follows semantic versioning.
Compatibility is based on Fractal Cloud’s public API.

### Contributing

Contributions are welcome.
Please open an issue before submitting any non-trivial pull requests.