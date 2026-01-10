Nice, this is exactly what you should know when working with Kafka on K8s 👍
Since you’re using **Kafka in KRaft mode inside AKS and running commands from `kafka-client` pod**, here’s your **complete practical command cheat-sheet** 🔥

> 📍 First go inside client pod:

```bash
kubectl exec -it -n sunbird kafka-client -- bash
cd /opt/kafka/bin
```

---

# ✅ TOPICS

### 🔹 List topics

```bash
./kafka-topics.sh --bootstrap-server kafka:9092 --list
```

### 🔹 Describe a topic (partitions, replicas)

```bash
./kafka-topics.sh --bootstrap-server kafka:9092 --describe --topic <topic-name>
```

### 🔹 Create a topic

```bash
./kafka-topics.sh --bootstrap-server kafka:9092 \
  --create --topic test-topic --partitions 3 --replication-factor 3
```

### 🔹 Delete a topic

```bash
./kafka-topics.sh --bootstrap-server kafka:9092 --delete --topic <topic-name>
```

⚠️ Topic deletion must be enabled in broker config (`delete.topic.enable=true`)

---

# ✅ PRODUCER

### 🔹 Produce messages manually

```bash
./kafka-console-producer.sh --bootstrap-server kafka:9092 --topic <topic-name>
```

Then type:

```
hello
world
```

---

# ✅ CONSUMER

### 🔹 Consume from latest

```bash
./kafka-console-consumer.sh --bootstrap-server kafka:9092 --topic <topic-name>
```

### 🔹 Consume from beginning

```bash
./kafka-console-consumer.sh --bootstrap-server kafka:9092 \
  --topic <topic-name> --from-beginning
```

### 🔹 Consume with group (important for lag testing)

```bash
./kafka-console-consumer.sh --bootstrap-server kafka:9092 \
  --topic <topic-name> --group test-group
```

---

# ✅ CONSUMER GROUPS

### 🔹 List consumer groups

```bash
./kafka-consumer-groups.sh --bootstrap-server kafka:9092 --list
```

### 🔹 Describe group (lag, members, partitions)

```bash
./kafka-consumer-groups.sh --bootstrap-server kafka:9092 \
  --describe --group <group-name>
```

### 🔹 Reset offsets (DANGEROUS – use carefully)

To earliest:

```bash
./kafka-consumer-groups.sh --bootstrap-server kafka:9092 \
  --group <group-name> --reset-offsets --to-earliest --execute --all-topics
```

To latest:

```bash
./kafka-consumer-groups.sh --bootstrap-server kafka:9092 \
  --group <group-name> --reset-offsets --to-latest --execute --all-topics
```

Dry-run:

```bash
./kafka-consumer-groups.sh --bootstrap-server kafka:9092 \
  --group <group-name> --reset-offsets --to-earliest --dry-run --all-topics
```

---

# ✅ BROKER / CLUSTER INFO

### 🔹 Check broker API versions

```bash
./kafka-broker-api-versions.sh --bootstrap-server kafka:9092
```

### 🔹 Check metadata quorum (KRaft)

```bash
./kafka-metadata-quorum.sh --bootstrap-server kafka:9092 describe --status
```

---

# ✅ ACLs (if security enabled)

### 🔹 List ACLs

```bash
./kafka-acls.sh --bootstrap-server kafka:9092 --list
```

---

# ✅ CONFIGS

### 🔹 Describe topic config

```bash
./kafka-configs.sh --bootstrap-server kafka:9092 \
  --entity-type topics --entity-name <topic-name> --describe
```

### 🔹 Describe broker config

```bash
./kafka-configs.sh --bootstrap-server kafka:9092 \
  --entity-type brokers --entity-name 0 --describe
```

---

# ⚠️ IMPORTANT for your setup (Sunbird + Flink)

* Kafka is **KRaft**, not Zookeeper
* Flink jobs already producing/consuming fine
* Always prefer **kafka-client pod** over local CLI
* Service to use:
  👉 `kafka:9092`

---

If you want, I can also give you:

* ✅ how to check lag for Flink jobs
* ✅ how to verify producer is writing
* ✅ how to trace message flow end-to-end

Just tell me what you want to debug next. 💪🔥
