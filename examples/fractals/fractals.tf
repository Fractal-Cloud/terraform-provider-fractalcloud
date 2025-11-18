resource "fractalcloud_fractal" "containerized" {
  resource_group_id = fractalcloud_resource_group.production_audi.id
  version = "1.0"
  description = "Microservice Fractal for our Organization"
  blueprint_components = [{
    id = "k8s-containerized"
    display_name = "Kubernetes cluster for workloads"
    type = "NetworkAndCompute.PaaS.ContainerPlatform"
    version = "1.2.12"
    locked = false
    parameters = {
      "apiGateway" = {
        "type" = "APIManagement.CaaS.Ambassador"
        "displayName" = "ambassador"
        "hostOwnerEmail" = "automotive@fractal.cloud"
        "acmeProviderAuthority" = "https://acme-v02.api.letsencrypt.org/directory"
        "namespace" = "ambassador-01"
      }
      "monitoring" = {
        "type" = "Observability.CaaS.Prometheus"
        "displayName" = "prometheus"
        "namespace" = "monitoring"
      }
      "logging"    = {
        "type" = "Observability.CaaS.Elastic"
        "displayName" = "elastic-logging"
        "namespace" = "elastic"
        "version" = "2.5.0"
        "withKibana" = true
        "withAPM" = true
        "withElasticVersion" = "8.5.0"
        "withStorage" = "10Gi"
      }
    }
    links = []
  }]
  interface_operations = [{
    name = "withExposedWorkload"
    parameters = ["workload, externalDns"]
    instructions = [{
      component_id = "k8s-containerized"
      operation = "withWorkload"
      input_parameters = ["workload"]
    },
    {
      component_id = "k8s-containerized"
      operation = "exposeService"
      input_parameters = ["workload", "externalDns"]
    }]
  }]
}

