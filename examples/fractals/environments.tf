resource "fractalcloud_management_environment" "production" {
  type = "Personal"
  owner_id = "xxx"
  display_name = "Production"
  azure_agent = {
    region = "westeurope"
    tenant_id = "xxx"
    subscription_id = "xxx"
  }
  gcp_agent = {
    region = "europe-west3"
    organization_id = "xxx"
    project_id = "xxx"
  }
  resource_groups = [
    fractalcloud_resource_group.production.id
  ]
}

resource "fractalcloud_operational_environment" "toyota_production" {
  management_environment_id = fractalcloud_management_environment.production.id
  display_name = "Toyota Production"
  agents = ["Gcp"]
  resource_groups = [
    fractalcloud_resource_group.production_toyota.id
  ]
}

resource "fractalcloud_operational_environment" "audi_production" {
  management_environment_id = fractalcloud_management_environment.production.id
  display_name = "Audi Production"
  agents = ["Azure"]
  resource_groups = [
    fractalcloud_resource_group.production_audi.id
  ]
}


