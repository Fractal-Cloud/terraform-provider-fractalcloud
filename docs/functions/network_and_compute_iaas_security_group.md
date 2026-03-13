---
page_title: "network_and_compute_iaas_security_group Function - Fractal Cloud"
subcategory: ""
description: |-
  Creates a SecurityGroup blueprint component.
---

# function: network_and_compute_iaas_security_group

Creates a SecurityGroup blueprint component. If `vpc` is provided, it is validated to ensure it is a VirtualNetwork component and added as a dependency. Ingress rules are serialized into parameters for the agent to reconcile.

## Example Usage

```terraform
locals {
  vpc = provider::fc::network_and_compute_iaas_virtual_network({
    id         = "main-vpc"
    cidr_block = "10.0.0.0/16"
  })

  sg = provider::fc::network_and_compute_iaas_security_group({
    id           = "web-sg"
    display_name = "Web Security Group"
    vpc          = local.vpc
    ingress_rules = [
      {
        from_port   = 443
        to_port     = 443
        protocol    = "tcp"
        source_cidr = "0.0.0.0/0"
      },
      {
        from_port              = 8080
        to_port                = 8080
        protocol               = "tcp"
        source_component_id    = "backend-vm"
      }
    ]
  })
}
```

## Signature

```text
network_and_compute_iaas_security_group(config object) object
```

## Arguments

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | String | Yes | Unique identifier for the component within the blueprint. |
| `display_name` | String | No | Human-readable name for the component. |
| `description` | String | No | Description of the component's purpose. |
| `vpc` | Component Object | No | A VirtualNetwork component to add as a dependency. Must be a component returned by `network_and_compute_iaas_virtual_network`. |
| `ingress_rules` | List of Object | No | List of ingress rule objects. Each rule supports the fields described below. |

### Ingress Rule Object

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `from_port` | Number | Yes | Start of the port range. |
| `to_port` | Number | No | End of the port range. Defaults to `from_port`. |
| `protocol` | String | No | Protocol (`"tcp"`, `"udp"`, `"icmp"`). Defaults to `"tcp"`. |
| `source_cidr` | String | No | Source CIDR block (e.g. `"0.0.0.0/0"`). Mutually exclusive with `source_component_id`. |
| `source_component_id` | String | No | ID of a source component. Mutually exclusive with `source_cidr`. |
