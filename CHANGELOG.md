## 0.1.0 (Unreleased)

FEATURES:

* **New Resource:** `fc_personal_bounded_context` - Manage personal bounded contexts
* **New Resource:** `fc_organizational_bounded_context` - Manage organizational bounded contexts
* **New Resource:** `fc_fractal` - Manage fractal definitions (blueprints) with component composition
* **New Resource:** `fc_management_environment` - Manage governance environments
* **New Resource:** `fc_operational_environment` - Manage runtime environments

* **New Data Source:** `fc_personal_bounded_context` - Look up personal bounded contexts
* **New Data Source:** `fc_organizational_bounded_context` - Look up organizational bounded contexts
* **New Data Source:** `fc_organization` - Look up organizations
* **New Data Source:** `fc_fractal` - Look up fractal definitions

* **New Provider Functions:** 46 blueprint component builder functions across 8 infrastructure domains:
  * NetworkAndCompute (7): virtual network, subnet, load balancer, security group, virtual machine, container platform, unmanaged
  * CustomWorkloads (5): CaaS, IaaS, PaaS, FaaS workloads, unmanaged
  * Storage (14): files/blobs, relational/document/column-oriented/key-value/graph databases, search, unmanaged
  * Messaging (5): PaaS and CaaS brokers and entities, unmanaged
  * BigData (6): distributed data processing, compute cluster, data processing job, ML experiment, datalake, unmanaged
  * APIManagement (3): PaaS and CaaS API gateways, unmanaged
  * Observability (4): monitoring, tracing, logging, unmanaged
  * Security (2): service mesh security, unmanaged
