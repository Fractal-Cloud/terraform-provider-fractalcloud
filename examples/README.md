# Examples

This directory contains examples that are mostly used for documentation, but can also be run/tested manually via the Terraform CLI.

The document generation tool looks for files in the following locations by default. All other *.tf files besides the ones mentioned below are ignored by the documentation tool. This is useful for creating examples that can run and/or are testable even if some parts are not relevant for the documentation.

* **provider/provider.tf** example file for the provider index page

## Example Directories

### `personal/`

Demonstrates Fractals within a **Personal Bounded Context**:

- **IaaS Fractal** — VPC, subnets, security groups, and VMs with port-based links for managed security group traffic rules.
- **Container Fractal** — Managed Kubernetes cluster with CaaS workloads linked for inter-service traffic (e.g. API service linking to a database on port 5432).

Shows how to use `provider::fc::network_and_compute_iaas_*` and `provider::fc::custom_workloads_caas_*` functions, including dependency wiring via direct object references and port-based links between compute components.

### `organizational/`

Demonstrates a multi-tier IaaS architecture within an **Organizational Bounded Context**:

- Separate web and app subnets with dedicated security groups.
- Web server links to app server on port 8080, generating managed security group rules automatically.

### `big-data/`

Demonstrates a complete end-to-end data platform Fractal:

```
Producer (PaaS Workload on K8s) → Ingest Stream (Messaging Entity) → ETL Cluster (Spark) → Data Lake
                                                                          ↑
                                                          Legacy Hadoop (SaaS, mounted as archive)
```

- **Container Platform** + **PaaS Workload** — a data-producer microservice running on managed Kubernetes.
- **Messaging Broker** + **Entity** — an event stream that carries raw events from the producer to the Spark cluster.
- **DistributedDataProcessing** platform (e.g. Databricks) linked to a **Datalake** and a **Legacy Hadoop** cluster via mount settings.
- **Compute Cluster** — autoscaling Spark cluster that consumes from the ingest stream.
- **Data Processing Job** — scheduled ETL that transforms raw events into curated datasets.
- **ML Experiment** — MLflow tracking for model development.

Shows generic links with settings (e.g. `mountName`, `consumerGroup`) to wire runtime relationships between components across different infrastructure domains.
