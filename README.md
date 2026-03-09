# Terraform Provider for Fractal Cloud

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

The Terraform Provider for [Fractal Cloud](https://fractal.cloud) enables platform and operations teams to manage Fractal Cloud governance resources through Terraform. It covers organizational structure, bounded contexts, environments, and fractal (blueprint) publication workflows.

Fractal Cloud is a platform engineering solution that delivers secure, compliant infrastructure across any cloud. It provides developers with ready-to-use building blocks and architecture templates while centralizing automation and governance for operations teams.

## Provider Scope

### Supported Resources

| Resource | Description |
|---|---|
| `fractalcloud_personal_bounded_context` | Personal bounded contexts (scoped to a user account) |
| `fractalcloud_organizational_bounded_context` | Organizational bounded contexts (scoped to an organization) |
| `fractalcloud_fractal` | Fractal definitions (blueprints) — publish, update, and manage component composition |
| `fractalcloud_management_environment` | Management environments for governance |
| `fractalcloud_operational_environment` | Operational environments for runtime governance |

### Supported Data Sources

| Data Source | Description |
|---|---|
| `fractalcloud_personal_bounded_context` | Look up an existing personal bounded context |
| `fractalcloud_organizational_bounded_context` | Look up an existing organizational bounded context |
| `fractalcloud_organization` | Look up an existing organization |
| `fractalcloud_fractal` | Look up an existing fractal definition |

### Not in Scope

This provider does **not** manage:

- **Live Systems** — created by development teams via the Fractal SDK or UI
- Application-level deployments or runtime operations
- CI/CD pipeline-triggered infrastructure instantiation

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.1
- [Go](https://golang.org/doc/install) >= 1.24 (to build the provider)
- A Fractal Cloud account with a service account

## Getting Started

### Provider Configuration

```hcl
terraform {
  required_providers {
    fractalcloud = {
      source  = "registry.terraform.io/fractalcloud/fc"
      version = "~> 0.1.0"
    }
  }
  required_version = ">= 1.1.0"
}

provider "fractalcloud" {
  service_account_id     = var.fc_service_account_id
  service_account_secret = var.fc_service_account_secret
}
```

### Authentication

The provider authenticates using a Fractal Cloud service account. Credentials can be set via:

- **HCL attributes**: `service_account_id` and `service_account_secret`
- **Environment variables**: `FRACTAL_CLOUD_SERVICE_ACCOUNT_ID` and `FRACTAL_CLOUD_SERVICE_ACCOUNT_SECRET`

Environment variables are used as defaults and can be overridden by HCL attributes.

| Attribute | Environment Variable | Required | Description |
|---|---|---|---|
| `service_account_id` | `FRACTAL_CLOUD_SERVICE_ACCOUNT_ID` | Yes | Service account identifier |
| `service_account_secret` | `FRACTAL_CLOUD_SERVICE_ACCOUNT_SECRET` | Yes | Service account secret |
| `host` | — | No | API endpoint (defaults to `https://api.fractal.cloud`) |

## Usage Examples

### Personal Bounded Context

```hcl
resource "fractalcloud_personal_bounded_context" "production" {
  short_name   = "production"
  display_name = "Production"
  description  = "Production bounded context"
}
```

### Organizational Bounded Context

```hcl
data "fractalcloud_organization" "my_org" {
  id = "a15f2627-c927-4125-a3a5-7141977135b1"
}

resource "fractalcloud_organizational_bounded_context" "platform" {
  short_name      = "platform"
  organization_id = data.fractalcloud_organization.my_org.id
  display_name    = "Platform"
  description     = "Shared platform bounded context"
}
```

### Fractal (Blueprint)

```hcl
resource "fractalcloud_fractal" "microservice" {
  bounded_context_id = fractalcloud_personal_bounded_context.production.id
  name        = "microservice-template"
  version     = "1.0"
  description = "Standard microservice architecture blueprint"
  components = [
    {
      id           = "container-platform-1"
      type         = "NetworkAndCompute.PaaS.ContainerPlatform"
      display_name = "Container Platform"
      description  = "Container Platform component"
      version      = "1.0"
      parameters   = {}
    },
    {
      id               = "search-cluster-1"
      type             = "Storage.CaaS.Search"
      display_name     = "Search Cluster"
      version          = "1.0"
      parameters       = {}
      dependencies_ids = ["container-platform-1"]
    }
  ]
}
```

### Data Source Lookup

```hcl
data "fractalcloud_personal_bounded_context" "existing" {
  short_name = "production"
}

data "fractalcloud_fractal" "existing" {
  bounded_context_id = data.fractalcloud_personal_bounded_context.existing.id
  name    = "microservice-template"
  version = "1.0"
}
```

See the [`examples/`](examples/) directory for complete working configurations.

## Building the Provider

```sh
git clone https://github.com/FractalCloud/terraform-provider.git
cd terraform-provider
make build
```

### Local Development

To use a locally built provider, add a dev override to your `~/.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/fractalcloud/fc" = "/path/to/your/go/bin"
  }
  direct {}
}
```

Then build and install:

```sh
make install
```

### Running Tests

```sh
# Unit tests
make test

# Acceptance tests (requires a Fractal Cloud account)
make testacc
```

## Contributing

Contributions are welcome. Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Make your changes
4. Run `make fmt` and `make lint` before committing
5. Write or update tests as appropriate
6. Open a pull request with a clear description of the change

Please open an issue before submitting non-trivial pull requests to discuss the approach.

### Development Requirements

- Go 1.24+
- [golangci-lint](https://golangci-lint.run/) for linting
- Terraform 1.1+ for acceptance testing

## License

This project is licensed under the GNU General Public License v3.0 — see the [LICENSE](LICENSE) file for details.
