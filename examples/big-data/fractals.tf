# Big Data Fractal with DistributedDataProcessing platform, compute clusters,
# data lake, processing jobs, ML experiment tracking, and an external resource.
#
# All components are defined as locals so that dependencies between them
# are expressed via direct object references (type-checked) rather than
# copy-pasted string IDs.

locals {
  # ── Core platform ──────────────────────────────────────────────────────
  spark_platform = provider::fractalcloud::bigdata_paas_distributed_data_processing({
    id           = "spark-platform"
    display_name = "Spark Processing Platform"
    description  = "Central Databricks workspace for all data workloads"
    pricing_tier = "standard"
  })

  # ── Data lake ──────────────────────────────────────────────────────────
  data_lake = provider::fractalcloud::bigdata_paas_datalake({
    id           = "data-lake"
    display_name = "Data Lake"
    description  = "Cloud object storage for raw ingestion and curated datasets"
  })

  # ── Compute clusters ───────────────────────────────────────────────────

  # Interactive cluster for ad-hoc analytics and development
  analytics_cluster = provider::fractalcloud::bigdata_paas_compute_cluster({
    id                       = "analytics-cluster"
    display_name             = "Analytics Cluster"
    description              = "Interactive cluster for ad-hoc queries and notebook development"
    platform                 = local.spark_platform
    cluster_name             = "analytics"
    spark_version            = "14.3.x-scala2.12"
    node_type_id             = "i3.xlarge"
    num_workers              = 2
    min_workers              = 1
    max_workers              = 8
    auto_termination_minutes = 30
    spark_conf = {
      "spark.databricks.delta.preview.enabled" = "true"
      "spark.sql.shuffle.partitions"           = "auto"
      "spark.databricks.io.cache.enabled"      = "true"
    }
    pypi_libraries = [
      "pandas>=2.0",
      "pyarrow>=14.0",
      "great-expectations>=0.18"
    ]
    maven_libraries = []
  })

  # Dedicated ETL cluster with autoscaling for batch workloads
  etl_cluster = provider::fractalcloud::bigdata_paas_compute_cluster({
    id                       = "etl-cluster"
    display_name             = "ETL Cluster"
    description              = "Autoscaling cluster for scheduled ETL pipelines"
    platform                 = local.spark_platform
    cluster_name             = "etl-processing"
    spark_version            = "14.3.x-scala2.12"
    node_type_id             = "m5.2xlarge"
    num_workers              = 4
    min_workers              = 2
    max_workers              = 20
    auto_termination_minutes = 15
    spark_conf = {
      "spark.dynamicAllocation.enabled"              = "true"
      "spark.shuffle.service.enabled"                = "true"
      "spark.sql.adaptive.enabled"                   = "true"
      "spark.databricks.delta.optimizeWrite.enabled" = "true"
    }
    maven_libraries = [
      "org.apache.hadoop:hadoop-aws:3.3.4",
      "io.delta:delta-core_2.12:2.4.0"
    ]
    pypi_libraries = []
  })

  # ── Data processing jobs ───────────────────────────────────────────────

  # Daily ingestion job: loads raw data from external sources into the bronze layer
  daily_ingestion = provider::fractalcloud::bigdata_paas_data_processing_job({
    id               = "daily-ingestion"
    display_name     = "Daily Data Ingestion"
    description      = "Ingests raw data from upstream sources into the bronze layer"
    platform         = local.spark_platform
    job_name         = "daily-data-ingestion"
    task_type        = "notebook"
    notebook_path    = "/Repos/data-team/pipelines/ingestion/daily_load"
    python_file      = null
    main_class_name  = null
    jar_uri          = null
    cron_schedule    = "0 2 * * *"
    max_retries      = 3
    existing_cluster = false
    parameters       = ["--env=prod", "--layer=bronze"]
  })

  # Hourly transformation job: curates bronze data into silver/gold layers
  hourly_transform = provider::fractalcloud::bigdata_paas_data_processing_job({
    id               = "hourly-transform"
    display_name     = "Hourly Data Transformation"
    description      = "Transforms and curates data from bronze to silver and gold layers"
    platform         = local.spark_platform
    job_name         = "hourly-data-transform"
    task_type        = "python"
    notebook_path    = null
    python_file      = "dbfs:/pipelines/transform/medallion_etl.py"
    main_class_name  = null
    jar_uri          = null
    cron_schedule    = "0 * * * *"
    max_retries      = 2
    existing_cluster = true
    parameters       = ["--env=prod", "--layers=silver,gold"]
  })

  # Weekly model retraining job
  weekly_model_training = provider::fractalcloud::bigdata_paas_data_processing_job({
    id               = "weekly-model-training"
    display_name     = "Weekly Model Training"
    description      = "Retrains ML models on latest curated data"
    platform         = local.spark_platform
    job_name         = "weekly-model-retrain"
    task_type        = "notebook"
    notebook_path    = "/Repos/ml-team/training/retrain_models"
    python_file      = null
    main_class_name  = null
    jar_uri          = null
    cron_schedule    = "0 6 * * 0"
    max_retries      = 1
    existing_cluster = false
    parameters       = ["--env=prod", "--experiment=churn-prediction"]
  })

  # Spark JAR batch job for high-performance data processing
  batch_aggregation = provider::fractalcloud::bigdata_paas_data_processing_job({
    id               = "batch-aggregation"
    display_name     = "Batch Aggregation Job"
    description      = "JVM-based batch aggregation for reporting datasets"
    platform         = local.spark_platform
    job_name         = "batch-aggregate-reports"
    task_type        = "jar"
    notebook_path    = null
    python_file      = null
    main_class_name  = "com.example.pipelines.AggregateReports"
    jar_uri          = "dbfs:/jars/aggregate-reports-1.0.jar"
    cron_schedule    = "0 4 * * *"
    max_retries      = 2
    existing_cluster = false
    parameters       = []
  })

  # ── ML experiment tracking ─────────────────────────────────────────────

  churn_experiment = provider::fractalcloud::bigdata_paas_ml_experiment({
    id                = "churn-experiment"
    display_name      = "Churn Prediction Experiment"
    description       = "MLflow experiment tracking for customer churn prediction models"
    platform          = local.spark_platform
    experiment_name   = "/Shared/ml-experiments/churn-prediction"
    artifact_location = "s3://data-platform-artifacts/ml/churn-prediction"
  })

  forecast_experiment = provider::fractalcloud::bigdata_paas_ml_experiment({
    id                = "forecast-experiment"
    display_name      = "Demand Forecasting Experiment"
    description       = "MLflow experiment tracking for demand forecasting models"
    platform          = local.spark_platform
    experiment_name   = "/Shared/ml-experiments/demand-forecasting"
    artifact_location = "s3://data-platform-artifacts/ml/demand-forecasting"
  })

  # ── External (unmanaged) resource ──────────────────────────────────────
  legacy_hadoop = provider::fractalcloud::bigdata_saas_unmanaged({
    id           = "legacy-hadoop"
    display_name = "Legacy Hadoop Cluster"
    description  = "On-premises Hadoop cluster for legacy batch workloads pending migration"
  })
}

resource "fractalcloud_fractal" "big_data_fractal" {
  bounded_context_id = data.fractalcloud_personal_bounded_context.existing_bounded_context.id
  name        = "data-platform"
  version     = "1.0"
  description = "Big Data Fractal with Spark processing, ETL jobs, ML experiments, and data lake storage"

  components = [
    local.spark_platform,
    local.data_lake,
    local.analytics_cluster,
    local.etl_cluster,
    local.daily_ingestion,
    local.hourly_transform,
    local.weekly_model_training,
    local.batch_aggregation,
    local.churn_experiment,
    local.forecast_experiment,
    local.legacy_hadoop,
  ]
}

output "big_data_fractal" {
  value = fractalcloud_fractal.big_data_fractal
}
