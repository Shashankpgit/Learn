# Kafka Version Compatibility Issue - Knowledge Transfer Document

## 🚨 Problem Statement
Flink jobs are failing to connect to Kafka with `UnsupportedVersionException` errors, causing missing consumer groups and failed data processing.

## 🔍 Root Cause Analysis

### The Issue
- **Kafka Server**: Running version 4.0.0 (upgraded)
- **Flink Applications**: Using Kafka client library 2.4.1 (outdated)
- **Result**: API version mismatch causing connection failures

### Error Symptoms
```bash
# Only one consumer group visible instead of multiple
kafka-consumer-groups.sh --bootstrap-server kafka:9092 --list
create-entity-consumer-group

# Kafka logs showing version errors
UnsupportedVersionException: Received request for api with key 0 (Produce) and unsupported version 2
```

## 📁 Affected Components

### 1. Infrastructure (Helm Charts)
```yaml
# Location: helmcharts/edbb/values.yaml
kafka:
  image:
    registry: "docker.io"
    repository: "bitnami/kafka"
    tag: "4.0.0"  # ← Kafka SERVER version
```

### 2. Application Code (Maven Projects)
```xml
<!-- Location: Both repositories' pom.xml -->
<kafka.version>2.4.1</kafka.version>  <!-- ← Kafka CLIENT LIBRARY version -->
```

## 🎯 Affected Repositories

### Repository 1: data-pipeline
- **Path**: `/home/sanketika7420/workspace/sunbird/data-pipeline`
- **Purpose**: LERN Flink jobs (Learning & Earning)
- **Jobs**: 
  - assessment-aggregator
  - activity-aggregate-updater
  - certificate-generator
  - notification-job
  - user-cache-updater

### Repository 2: knowledge-platform-jobs
- **Path**: `/home/sanketika7420/workspace/sunbird/knowledge-platform-jobs`
- **Purpose**: Knowledge Platform Flink jobs
- **Jobs**:
  - content-publish
  - search-indexer
  - qrcode-image-generator
  - post-publish-processor

## 💡 Solution Options

### Option A: Downgrade Kafka Server (Quick Fix - Recommended)
```yaml
# File: helmcharts/edbb/values.yaml
kafka:
  image:
    tag: "3.6.1"  # Change from 4.0.0 to 3.6.1
```

**Pros**: 
- ✅ No code changes required
- ✅ Immediate fix
- ✅ All existing jobs work

**Cons**: 
- ❌ Missing latest Kafka 4.0 features

### Option B: Upgrade Client Libraries (Proper Fix)
```xml
<!-- File: data-pipeline/pom.xml -->
<!-- File: knowledge-platform-jobs/pom.xml -->
<kafka.version>3.6.1</kafka.version>  <!-- Change from 2.4.1 to 3.6.1 -->
```

**Pros**: 
- ✅ Future-proof solution
- ✅ Uses latest Kafka features

**Cons**: 
- ❌ Requires code changes
- ❌ Need to rebuild all job images
- ❌ Testing required

## 🛠️ Implementation Steps

### For Option A (Recommended)
1. Edit `helmcharts/edbb/values.yaml`
2. Change Kafka image tag to `3.6.1`
3. Redeploy Kafka
4. Verify consumer groups appear

### For Option B (If needed)
1. Update `data-pipeline/pom.xml`
2. Update `knowledge-platform-jobs/pom.xml`
3. Rebuild all Flink job images
4. Update image tags in Helm charts
5. Redeploy all jobs

## 🔍 Verification Steps

### Check Consumer Groups
```bash
# Should show multiple groups after fix
kubectl exec -it kafka-client -- kafka-consumer-groups.sh --bootstrap-server kafka:9092 --list

# Expected output:
# dev-assessment-aggregator-group
# dev-activity-aggregate-group
# dev-certificate-generator-group
# dev-notification-group
# ... (more groups)
```

### Check Kafka Logs
```bash
# Should not show version errors
kubectl logs kafka-controller-0 -n sunbird
```

### Check Flink Job Status
```bash
# All jobs should be running
kubectl get pods -n sunbird | grep flink
```

## 📊 Impact Assessment

### Before Fix
- ❌ Missing consumer groups
- ❌ Failed telemetry processing
- ❌ Broken assessment aggregation
- ❌ Certificate generation issues
- ❌ Notification failures

### After Fix
- ✅ All consumer groups visible
- ✅ Data processing restored
- ✅ Assessment scores calculated
- ✅ Certificates generated
- ✅ Notifications sent

## 🚀 Recommended Action Plan

1. **Immediate**: Use Option A (downgrade Kafka to 3.6.1)
2. **Short-term**: Monitor system stability
3. **Long-term**: Plan Option B implementation if Kafka 4.0 features needed

## 📞 Key Contacts

- **Infrastructure Team**: For Helm chart changes
- **Development Team**: For client library upgrades
- **DevOps Team**: For deployment coordination

## 📚 Additional Resources

- [Kafka Compatibility Matrix](https://kafka.apache.org/protocol.html)
- [Flink Kafka Connector Docs](https://nightlies.apache.org/flink/flink-docs-stable/docs/connectors/datastream/kafka/)
- [Sunbird Architecture Docs](https://sunbird.org)

---
**Document Version**: 1.0  
**Last Updated**: January 2026  
**Author**: Infrastructure Team