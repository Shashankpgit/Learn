# Obsrv Repository Overview

## What is Obsrv?

**Obsrv** (Observability) is a comprehensive observability and data analytics platform designed to collect, process, store, and analyze telemetry and event data at scale. This repository (`obsrv-automation`) contains the **Kubernetes deployment automation** for the Obsrv platform using Helm charts and Terraform.

The platform is built on a modern data stack including:
- **Apache Kafka** for event streaming
- **Apache Flink** for real-time stream processing
- **Apache Druid** for time-series analytics and OLAP queries
- **Apache Spark** for batch processing
- **PostgreSQL** for metadata and configuration storage
- **Trino** for SQL querying across data sources
- **Superset** for data visualization
- Various monitoring and observability tools (Prometheus, Grafana, Loki)

---

## Repository Structure

### Key Directories

- **`helmcharts/`** - Helm charts for deploying all Obsrv services to Kubernetes
- **`Dockerfiles/`** - Container images for various components
- **`exporters/`** - Custom exporters for metrics and monitoring
- **`connectors/`** - Data connectors for various sources
- **`terraform/`** - Infrastructure as Code (AWS, Azure, GCP) - *ignored per your request*
- **`infra-setup/`** - Installation scripts and configuration

---

## Helmcharts Directory Structure

The `helmcharts/` directory is the core of this repository:

```
helmcharts/
├── services/          # Individual service Helm charts
├── obsrv/            # Umbrella chart that bundles all services
├── bootstrapper/     # Initial setup and CRDs
├── base/             # Base templates and common configurations
├── global-values.yaml          # Global configuration
├── global-cloud-values-*.yaml # Cloud-specific configurations
└── images.yaml                # Container image versions
```

---

## Services Overview

### 🎨 **Frontend & Backend APIs**

#### **web-console**
- **Purpose**: Frontend web application for Obsrv platform
- **Technology**: React/Node.js web console
- **Port**: 3000
- **Function**: User interface for managing datasets, viewing dashboards, configuring ingestion pipelines

#### **dataset-api**
- **Purpose**: REST API for dataset management
- **Function**: CRUD operations for datasets, schema management, dataset configuration
- **Port**: 8000
- **Database**: PostgreSQL

#### **config-api**
- **Purpose**: Configuration management API
- **Function**: Manages system configurations, ingestion rules, transformation rules
- **Port**: 8000
- **Database**: PostgreSQL

#### **command-api**
- **Purpose**: Command execution API for data pipeline operations
- **Function**: Submits Flink jobs, manages ingestion pipelines, orchestrates data processing tasks
- **Port**: 8000
- **Integration**: Flink, Kafka

---

### 🔄 **Data Processing Services**

#### **flink**
- **Purpose**: Real-time stream processing engine
- **Technology**: Apache Flink
- **Function**: Processes streaming data from Kafka, performs transformations, aggregations, and real-time analytics
- **Port**: 8081 (REST API)
- **Integration**: Kafka (consumes), Druid (writes processed data)

#### **spark**
- **Purpose**: Batch processing and ETL jobs
- **Technology**: Apache Spark with Apache Livy
- **Function**: Large-scale batch data processing, data transformations, machine learning workloads
- **Integration**: Kafka, S3/MinIO, Druid

#### **lakehouse-connector**
- **Purpose**: Connector service for data lakehouse operations
- **Function**: Bridges data between different storage systems, handles data lake ingestion
- **Technology**: Flink-based connector
- **Integration**: Flink, cloud storage (S3/Azure/GCS)

#### **submit-ingestion**
- **Purpose**: Service for submitting data ingestion jobs
- **Function**: Initiates ingestion pipelines, submits jobs to Flink/Spark, manages ingestion workflows
- **Technology**: CronJob/Job that triggers ingestion processes

---

### 💾 **Data Storage Services**

#### **druid-raw-cluster**
- **Purpose**: Time-series OLAP database for analytics
- **Technology**: Apache Druid
- **Function**: Stores and queries time-series data, supports fast aggregations and time-based queries
- **Components**: Coordinator, Overlord, Broker, Historical, MiddleManager
- **Storage**: Cloud storage (S3/Azure/GCS) for deep storage
- **Metadata**: PostgreSQL

#### **postgresql**
- **Purpose**: Relational database for metadata and configuration
- **Function**: Stores dataset schemas, ingestion configurations, user data, system metadata
- **Databases**: Multiple databases for different services (druid_raw, obsrv, superset, etc.)

#### **minio**
- **Purpose**: Object storage (S3-compatible)
- **Function**: Stores raw data files, checkpoints, backups, intermediate processing data
- **Use Cases**: Data lake storage, Flink checkpoints, backup storage

#### **hms** (Hive Metastore)
- **Purpose**: Metadata management for data lake
- **Function**: Manages table schemas, partitions, and metadata for data stored in object storage
- **Port**: 9083
- **Integration**: Trino, Spark

---

### 📨 **Message Queue & Streaming**

#### **kafka** / **kafka40**
- **Purpose**: Distributed event streaming platform
- **Function**: Message broker for event streaming, decouples producers and consumers
- **Topics**: Handles telemetry events, ingestion streams, processing queues
- **Versions**: Supports both Kafka and Kafka 4.0

#### **valkey-dedup** / **valkey-denorm**
- **Purpose**: In-memory data structures (Redis/Valkey)
- **Function**: 
  - **valkey-dedup**: Deduplication cache to prevent duplicate event processing
  - **valkey-denorm**: Denormalization cache for fast lookups during stream processing
- **Technology**: Valkey (Redis fork)

---

### 🔍 **Query & Analytics Services**

#### **trino**
- **Purpose**: Distributed SQL query engine
- **Function**: Query data across multiple sources (Druid, PostgreSQL, object storage) with SQL
- **Port**: 8080
- **Integration**: Connects to Druid, PostgreSQL, Hive Metastore, S3

#### **superset**
- **Purpose**: Data visualization and business intelligence platform
- **Function**: Create dashboards, charts, and reports from Obsrv data
- **Port**: 8088
- **Database**: PostgreSQL (for Superset metadata)
- **Integration**: Connects to Druid, Trino, PostgreSQL for querying

---

### 📊 **Monitoring & Observability Services**

#### **kube-prometheus-stack**
- **Purpose**: Complete Prometheus monitoring stack
- **Components**: Prometheus, Grafana, AlertManager, Node Exporter
- **Function**: Metrics collection, alerting, visualization dashboards

#### **prometheus-pushgateway**
- **Purpose**: Metrics gateway for short-lived jobs
- **Function**: Allows batch jobs and short-lived services to push metrics to Prometheus

#### **loki**
- **Purpose**: Log aggregation system
- **Function**: Collects, stores, and queries logs from all services
- **Integration**: Works with Grafana for log visualization

#### **promtail**
- **Purpose**: Log shipper for Loki
- **Function**: Collects logs from Kubernetes pods and ships them to Loki

#### **grafana-configs**
- **Purpose**: Pre-configured Grafana dashboards
- **Function**: Provides ready-made dashboards for monitoring Obsrv services

#### **druid-exporter**
- **Purpose**: Prometheus exporter for Druid metrics
- **Function**: Exposes Druid health, datasource, segment, and task metrics to Prometheus

#### **kafka-exporter**
- **Purpose**: Prometheus exporter for Kafka metrics
- **Function**: Exposes Kafka broker, topic, and consumer group metrics

#### **postgresql-exporter**
- **Purpose**: Prometheus exporter for PostgreSQL metrics
- **Function**: Exposes database performance and health metrics

#### **kafka-message-exporter**
- **Purpose**: Custom exporter for Kafka message statistics
- **Function**: Exports Kafka message counts, throughput, and processing stats to Prometheus
- **Technology**: Node.js service

#### **azure-exporter**
- **Purpose**: Exporter for Azure-specific metrics
- **Function**: Collects Azure resource metrics and cloud storage statistics

#### **s3-exporter**
- **Purpose**: Exporter for S3/object storage metrics
- **Function**: Monitors S3 bucket usage, object counts, storage metrics

#### **opentelemetry-collector**
- **Purpose**: Telemetry data collection and processing
- **Function**: Collects traces, metrics, and logs using OpenTelemetry standard
- **Integration**: Exports to Prometheus, Loki, and other backends

---

### 🔐 **Security & API Gateway Services**

#### **kong**
- **Purpose**: API Gateway and microservices management
- **Function**: Routes API requests, handles authentication, rate limiting, load balancing
- **Integration**: Routes to dataset-api, config-api, command-api, web-console

#### **keycloak**
- **Purpose**: Identity and access management
- **Function**: Single sign-on (SSO), user authentication, authorization, OAuth2/OIDC provider
- **Integration**: Used by web-console and APIs for authentication

#### **kong-ingress-routes**
- **Purpose**: Kong route configurations
- **Function**: Defines API routes, services, and plugins for Kong gateway

#### **cert-manager**
- **Purpose**: Certificate management for TLS/SSL
- **Function**: Automatically provisions and renews SSL certificates (Let's Encrypt)

#### **letsencrypt-ssl**
- **Purpose**: SSL certificate management
- **Function**: Manages Let's Encrypt certificates for HTTPS endpoints

---

### 🔌 **Data Connectors & Ingestion**

#### **system-rules-ingestor**
- **Purpose**: Ingests system rules and configurations
- **Function**: Processes system-level rules, validation rules, transformation rules into the system

#### **masterdata-indexer-cron**
- **Purpose**: Cron job for indexing master data
- **Function**: Periodically indexes reference/master data into search systems or databases

#### **secor**
- **Purpose**: Kafka to object storage connector
- **Function**: Consumes messages from Kafka and writes them to S3/Azure/GCS in partitioned format
- **Output**: Stores raw events in object storage for batch processing and archival

---

### 🛠️ **Backup & Maintenance Services**

#### **postgresql-backup**
- **Purpose**: Automated PostgreSQL backups
- **Function**: Scheduled backups of PostgreSQL databases to cloud storage
- **Schedule**: Cron-based backup jobs

#### **redis-backup**
- **Purpose**: Automated Redis/Valkey backups
- **Function**: Backs up Redis data to cloud storage

#### **velero**
- **Purpose**: Kubernetes backup and disaster recovery
- **Function**: Backs up entire Kubernetes cluster state, persistent volumes, and configurations

#### **postgresql-migration**
- **Purpose**: Database migration service
- **Function**: Runs Flyway/Liquibase migrations to update database schemas
- **Files**: SQL migration scripts in `configs/migrations/`

#### **volume-autoscaler**
- **Purpose**: Automatic volume scaling
- **Function**: Automatically increases persistent volume sizes when they fill up

---

### 🔧 **Infrastructure Services**

#### **kubernetes-reflector**
- **Purpose**: Kubernetes secret/configmap reflector
- **Function**: Replicates secrets and configmaps across namespaces

#### **druid-operator**
- **Purpose**: Kubernetes operator for Druid
- **Function**: Manages Druid cluster lifecycle, scaling, and configuration

---

### 📋 **Configuration Services**

#### **alert-rules**
- **Purpose**: Prometheus alerting rules
- **Function**: Defines alert rules for monitoring Obsrv services
- **Integration**: Works with AlertManager

---

## Data Flow Architecture

```
Data Sources → Kafka → Flink/Spark → Druid → Superset/Trino
                ↓
            Secor → Object Storage (S3/MinIO) → Trino
                ↓
            Valkey (dedup/denorm) → Flink
```

1. **Ingestion**: Data arrives via APIs or connectors → Kafka
2. **Processing**: Flink processes streams in real-time; Spark handles batch jobs
3. **Storage**: 
   - Processed data → Druid (for analytics)
   - Raw data → Object Storage via Secor
   - Metadata → PostgreSQL
4. **Query**: Trino and Superset query data from Druid and object storage
5. **Monitoring**: All services export metrics to Prometheus, logs to Loki

---

## Deployment Model

- **Helm Charts**: Each service has its own Helm chart in `services/`
- **Umbrella Chart**: `obsrv/` chart bundles all services together
- **Cloud Support**: AWS, Azure, GCP (via terraform and cloud-specific values)
- **Base Templates**: `base/` provides reusable Kubernetes templates

---

## Key Features

1. **Multi-cloud Support**: Works on AWS, Azure, and GCP
2. **Scalable**: Horizontal scaling for all services
3. **Observable**: Built-in monitoring with Prometheus, Grafana, Loki
4. **Real-time & Batch**: Supports both streaming (Flink) and batch (Spark) processing
5. **SQL Interface**: Query data using Trino SQL
6. **Visualization**: Superset for dashboards and reports
7. **API-driven**: REST APIs for all operations
8. **Secure**: Keycloak for auth, Kong for API gateway, TLS support

---

## Notes

- **Frontend/Backend APIs** (dataset-api, config-api, command-api, web-console): These are application services whose full source code is in separate repositories. This repo only contains their Helm deployment charts.
- **Infrastructure Code**: Terraform code is present but excluded from this explanation as requested.
- **Dockerfiles**: Custom container images are built in `Dockerfiles/` directory for services like Flink connectors, Spark with Livy, etc.

