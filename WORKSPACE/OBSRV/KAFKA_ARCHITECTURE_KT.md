# Kafka Architecture Knowledge Transfer - Obsrv Platform

## Overview
Obsrv uses Apache Kafka as the central message streaming backbone for data ingestion, processing, and distribution across all system components.

## Kafka Infrastructure

### Core Setup
- **Bootstrap Server**: `kafka-headless.kafka.svc.cluster.local:9092`
- **Zookeeper**: `kafka-zookeeper-headless.kafka.svc.cluster.local:2181`
- **Namespace**: `kafka`
- **Partitions**: 4 per topic (default)
- **Replication Factor**: 1 (single node)
- **Version**: Kafka 3.3.1 (Bitnami)

## Kafka Topics

### Data Processing Pipeline Topics

#### 1. `ingest` Topic
- **Purpose**: Entry point for all raw data
- **Producers**: Dataset API, external data sources
- **Consumers**: Flink validation jobs
- **Flow**: Raw events → Validation

#### 2. `raw` Topic  
- **Purpose**: Validated raw events
- **Producers**: Flink validation jobs
- **Consumers**: Secor backup, deduplication jobs
- **Flow**: Validated events → Deduplication

#### 3. `unique` Topic
- **Purpose**: Deduplicated events
- **Producers**: Flink deduplication jobs
- **Consumers**: Secor backup, denormalization jobs
- **Flow**: Unique events → Denormalization

#### 4. `denorm` Topic
- **Purpose**: Denormalized/enriched events
- **Producers**: Flink denormalization jobs
- **Consumers**: Secor backup, transformation jobs
- **Flow**: Enriched events → Transformation

#### 5. `transform` Topic
- **Purpose**: Final transformed events
- **Producers**: Flink transformation jobs
- **Consumers**: Secor backup, storage connectors
- **Flow**: Final events → Storage

#### 6. `failed` Topic
- **Purpose**: Processing failures at any stage
- **Producers**: All Flink/Spark jobs (error handling)
- **Consumers**: Secor backup, monitoring systems
- **Flow**: Failed events → Analysis/Reprocessing

### System & Monitoring Topics

#### 7. `system.events` Topic
- **Purpose**: System operational events
- **Producers**: All system components
- **Consumers**: Druid ingestion, Secor backup
- **Flow**: System events → Druid → Dashboards

#### 8. `system.telemetry.events` Topic
- **Purpose**: Component telemetry data
- **Producers**: Flink jobs, Spark jobs, system services
- **Consumers**: Druid ingestion, Secor backup
- **Flow**: Telemetry → Druid → Monitoring

#### 9. `stats` Topic
- **Purpose**: Processing job statistics
- **Producers**: Spark processing jobs
- **Consumers**: Kafka Message Exporter → Prometheus
- **Flow**: Job stats → Metrics

### Specialized Topics

#### 10. `masterdata.ingest` Topic
- **Purpose**: Master data ingestion
- **Producers**: Master data APIs, batch jobs
- **Consumers**: Master data processors, Secor backup

#### 11. `masterdata.stats` Topic
- **Purpose**: Master data processing metrics
- **Producers**: Master data processing jobs
- **Consumers**: Monitoring systems

#### 12. `hudi.connector.in` Topic
- **Purpose**: Lakehouse connector input
- **Producers**: Data processing pipelines
- **Consumers**: Lakehouse Connector (Flink)
- **Flow**: Processed data → Hudi tables → Data lake

#### 13. `obsrv-connectors-metrics` Topic
- **Purpose**: Connector performance metrics
- **Producers**: Flink/Spark connectors
- **Consumers**: Kafka Message Exporter → Prometheus

#### 14. `connectors.failed` Topic
- **Purpose**: Connector operation failures
- **Producers**: All connector services
- **Consumers**: Monitoring/alerting systems

## Kafka Producers

### 1. Dataset API Service
- **Location**: `dataset-api.dataset-api.svc.cluster.local:3000`
- **Topics**: `ingest`, `system.events`
- **Role**: Primary data ingestion endpoint
- **Config**: Snappy compression, 5 retries, 3s delay

### 2. Flink Processing Jobs
- **Location**: Flink namespace
- **Topics**: `raw`, `unique`, `denorm`, `transform`, `failed`, `system.events`
- **Role**: Stream processing pipeline stages
- **Config**: 10MB max request, 98KB batch, 10ms linger, Snappy compression

### 3. Spark Processing Jobs
- **Location**: Spark namespace
- **Topics**: `stats`, `system.events`, `obsrv-connectors-metrics`
- **Role**: Batch processing and connectors
- **Config**: 1MB max request, Snappy compression

### 4. Lakehouse Connector
- **Location**: Flink-based Hudi connector
- **Topics Produced**: `system.events`, `failed`
- **Topics Consumed**: `hudi.connector.in`
- **Role**: Data lake integration

### 5. System Components
- **Topics**: `system.telemetry.events`, `system.events`
- **Role**: System monitoring and telemetry

## Kafka Consumers

### 1. Flink Stream Processing
- **Topics**: `ingest`, `raw`, `unique`, `denorm`, `hudi.connector.in`
- **Role**: Real-time stream processing pipeline
- **Parallelism**: Configurable per job (default: 1)

### 2. Secor Backup Service
- **Topics**: All data topics for archival
- **Role**: Data backup to cloud storage (S3/Azure/GCS)
- **Config**: 100MB max file, 15min max age, timestamp partitioning

### 3. Druid Ingestion
- **Topics**: `system.events`, `system.telemetry.events`
- **Role**: Real-time analytics ingestion
- **Config**: Hourly segments, 1-hour tasks

### 4. Kafka Message Exporter
- **Topics**: `obsrv-connectors-metrics`, `stats`
- **Role**: Export messages to Prometheus metrics
- **Consumer Group**: `kafka-message-exporter`

### 5. Monitoring Systems
- **Topics**: `failed`, `system.events`, `connectors.failed`
- **Role**: System health monitoring and alerting

## Data Flow Patterns

### Primary Processing Pipeline
```
External Data → Dataset API → ingest → Flink Validation → raw → 
Flink Deduplication → unique → Flink Denormalization → denorm → 
Flink Transformation → transform → Storage Systems
```

### Error Handling Flow
```
Processing Error → failed → Secor Backup + Monitoring → Alerts
```

### System Monitoring Flow
```
Components → system.events/system.telemetry.events → Druid → Dashboards
```

### Metrics Flow
```
Connectors → obsrv-connectors-metrics → Kafka Message Exporter → Prometheus
```

## Configuration Details

### Global Settings
- **Partitions**: 4 per topic
- **Replication**: 1 (single node)
- **Compression**: Snappy (all producers)
- **Retention**: 7 days default
- **Max Message Size**: 5MB

### Producer Settings
- **Batch Size**: 98KB (Flink), 1MB (Spark)
- **Linger Time**: 10ms
- **Retries**: 5 with 3s initial delay
- **Idempotence**: Enabled

### Consumer Settings
- **Auto Offset Reset**: Latest
- **Session Timeout**: 30s
- **Heartbeat Interval**: 3s
- **Enable Auto Commit**: True

## Monitoring & Observability

### Metrics Available
- **JMX Exporter**: Kafka broker metrics
- **Kafka Exporter**: Topic-specific metrics
- **Consumer Lag**: Per consumer group
- **Throughput**: Messages/second per topic

### Alerting Thresholds
- **Consumer Lag**: >50k (warning), >100k (critical)
- **Topic Unavailability**: Broker down
- **High Error Rate**: Excessive failed topic messages

## Backup & Recovery

### Secor Backup
- **Storage**: Cloud storage (S3/Azure/GCS)
- **Format**: JSON with Snappy compression
- **Partitioning**: Hourly timestamp-based
- **Retention**: Configurable per topic

### Recovery Process
1. Topic recreation via Helm charts
2. Data recovery from cloud backups
3. Consumer group offset reset if needed

## Zookeeper Dependencies

### Why Flink Uses Zookeeper
Despite Zookeeper being primarily for Kafka coordination, Flink references Zookeeper for:
- **Legacy Kafka Consumer Integration**: Older Kafka client libraries that store offsets in Zookeeper
- **Topic Metadata Discovery**: Some Kafka clients query Zookeeper for topic/partition information
- **Consumer Group Coordination**: Legacy consumer group management approach

**Flink Configuration:**
```yaml
kafka {
  zookeeper = "kafka-zookeeper-headless.kafka.svc.cluster.local:2181"
}
```

### Why Secor Uses Zookeeper
Secor uses Zookeeper for:
- **Consumer Offset Storage**: Stores consumer offsets in Zookeeper (legacy approach)
- **Consumer Group Management**: Coordinating which Secor instances consume from which partitions
- **Topic Discovery**: Getting topic and partition metadata

**Secor Configuration:**
```properties
zookeeper.session.timeout.ms=3000
zookeeper.sync.time.ms=200
secor.zookeeper.path=/
kafka.zookeeper.path=/
```

### Modern vs Legacy Approach
**Modern Kafka (0.9+):**
- Consumer offsets stored in Kafka's `__consumer_offsets` topic
- No Zookeeper dependency for consumers
- Direct broker communication for metadata

**Legacy Kafka / Legacy Clients:**
- Consumer offsets stored in Zookeeper
- Topic metadata from Zookeeper
- Consumer group coordination via Zookeeper

### Recommendations
- **Update Client Libraries**: Use newer Kafka clients that don't require Zookeeper
- **Remove Legacy References**: Test removing Zookeeper configs if using modern clients
- **Gradual Migration**: Move to Kafka-native offset management

## Security

### Current Setup
- **Authentication**: PLAINTEXT (internal cluster)
- **Authorization**: Open access within cluster
- **Network**: Kubernetes network policies for isolation
- **TLS**: Disabled for internal communication

## Performance Tuning

### Broker Optimization
- **Log Segment Size**: 1GB
- **Flush Interval**: 10k messages or 1s
- **Socket Buffers**: 102KB
- **Compression**: Snappy

### Producer Optimization
- **Batching**: Enabled with optimal sizes
- **Compression**: Snappy for best performance
- **Acknowledgments**: All replicas

### Consumer Optimization
- **Parallelism**: Matches partition count
- **Fetch Size**: Optimized per consumer
- **Processing**: At-least-once delivery

## Troubleshooting

### Common Issues
1. **Consumer Lag**: Check job health and scaling
2. **Broker Unavailable**: Verify Kafka pod status
3. **Message Loss**: Check producer ack settings
4. **Duplicates**: Verify consumer offset management

### Debug Tools
- **Kafka Console**: Available in Kafka pods
- **Metrics**: Prometheus/Grafana dashboards
- **Logs**: Centralized via Loki
- **CLI Tools**: kafka-topics, kafka-console-consumer/producer

### Key Commands
```bash
# List topics
kubectl exec -it kafka-0 -n kafka -- kafka-topics.sh --bootstrap-server localhost:9092 --list

# Check consumer groups
kubectl exec -it kafka-0 -n kafka -- kafka-consumer-groups.sh --bootstrap-server localhost:9092 --list

# Monitor topic
kubectl exec -it kafka-0 -n kafka -- kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic ingest --from-beginning
```

## Service Dependencies

### Critical Dependencies
- **Zookeeper**: Required for Kafka cluster coordination
- **PostgreSQL**: Metadata storage for various services
- **Valkey/Redis**: Caching for deduplication and denormalization
- **Druid**: Real-time analytics ingestion

### Network Dependencies
- All services communicate via Kubernetes service discovery
- DNS resolution within cluster for service-to-service communication
- Load balancing handled by Kubernetes services

## Operational Procedures

### Scaling Kafka
1. Update replica count in Helm values
2. Redeploy Kafka chart
3. Verify all brokers are healthy
4. Update producer/consumer configurations if needed

### Adding New Topics
1. Add topic definition to global-values.yaml
2. Update Kafka provisioning configuration
3. Redeploy Kafka chart
4. Verify topic creation and partition assignment

### Monitoring Health
1. Check Kafka broker metrics in Grafana
2. Monitor consumer lag across all groups
3. Verify Secor backup operations
4. Check failed topic for processing errors

This document serves as the complete reference for understanding Kafka's role in the Obsrv platform. New team members should use this as their primary guide for Kafka operations and troubleshooting.