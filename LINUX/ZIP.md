Here’s a **quick reference guide** for `zip` and `unzip` commands — concise, practical, and copy-paste friendly.

---

# 📦 ZIP & UNZIP – Quick Reference Guide

## 🔹 `zip` — Create ZIP archives

### **Basic Syntax**

```bash
zip [options] archive.zip file_or_directory
```

---

### **Zip a single file**

```bash
zip app.zip app.yaml
```

---

### **Zip multiple files**

```bash
zip configs.zip values.yaml Chart.yaml README.md
```

---

### **Zip a directory (MOST COMMON)**

```bash
zip -r units-flyway-migration.zip units-flyway-migration/
```

📌 `-r` → recursive (required for directories)

---

### **Zip current directory contents**

```bash
zip -r project.zip .
```

⚠️ Includes **everything** in the directory (be careful with `.git`, `node_modules`, etc.)

---

### **Exclude files or directories**

```bash
zip -r app.zip . -x "*.log" "node_modules/*"
```

---

### **Quiet mode (less output)**

```bash
zip -rq app.zip app/
```

---

### **Update zip (add only changed/new files)**

```bash
zip -ur app.zip app/
```

---

### **Encrypt a zip**

```bash
zip -e secrets.zip secrets.txt
```

---

### **Compression levels**

```bash
zip -r -9 max-compress.zip app/   # best compression
zip -r -1 fast.zip app/           # fastest
```

---

### **Test a zip file**

```bash
zip -T app.zip
```

---

## 🔹 `unzip` — Extract ZIP archives

### **Basic Syntax**

```bash
unzip archive.zip
```

---

### **Extract to current directory**

```bash
unzip app.zip
```

---

### **Extract to a specific directory**

```bash
unzip app.zip -d /tmp/app
```

---

### **List contents (no extraction)**

```bash
unzip -l app.zip
```

---

### **Extract specific file**

```bash
unzip app.zip values.yaml
```

---

### **Extract without overwriting**

```bash
unzip -n app.zip
```

---

### **Overwrite without prompt**

```bash
unzip -o app.zip
```

---

### **Quiet extraction**

```bash
unzip -q app.zip
```

---

### **Extract password-protected zip**

```bash
unzip secrets.zip
```

(You’ll be prompted for the password)

---

## 🔹 Common Mistakes ❌

❌ Missing recursive flag

```bash
zip app.zip app/
# ❌ "Nothing to do"
```

✅ Correct:

```bash
zip -r app.zip app/
```

---

❌ Wrong order

```bash
zip folder/
```

✅ Correct order:

```bash
zip archive.zip folder/
```

---

## 🔹 Handy One-Liners

| Task             | Command                       |
| ---------------- | ----------------------------- |
| Zip folder       | `zip -r app.zip app/`         |
| Unzip            | `unzip app.zip`               |
| List zip         | `unzip -l app.zip`            |
| Exclude files    | `zip -r app.zip . -x "*.log"` |
| Fast zip         | `zip -r -1 fast.zip app/`     |
| Best compression | `zip -r -9 best.zip app/`     |

---

## 🔹 Helm / DevOps Tip 💡

For Helm charts and CI pipelines, `tar.gz` is more common:

```bash
tar -czvf chart.tgz units-flyway-migration/
```

---
Below is a **DevOps-focused, practical guide** covering both requests:
1️⃣ **ZIP vs TAR (with GZIP)** comparison
2️⃣ **ZIP usage inside Docker & Kubernetes pipelines**

No fluff — this is what you actually need in real infra work.

---

# 🧠 DevOps-Focused: ZIP vs TAR Comparison

## 🔹 High-level summary

| Aspect                | `zip`        | `tar.gz` (`tar + gzip`)         |
| --------------------- | ------------ | ------------------------------- |
| Native to Linux       | ❌ (optional) | ✅ (default)                     |
| Compression           | Built-in     | External (`gzip`, `xz`, `zstd`) |
| Preserves permissions | ❌ (limited)  | ✅                               |
| Preserves symlinks    | ❌            | ✅                               |
| Streaming support     | ❌            | ✅                               |
| CI/CD friendliness    | ⚠️           | ✅                               |
| Helm compatibility    | ❌            | ✅                               |
| Cross-platform        | ✅            | ⚠️ (mostly Linux)               |

---

## 🔹 Why `tar.gz` dominates DevOps & Cloud

### ✅ File permissions & ownership

Critical for:

* Docker images
* Kubernetes manifests
* Init containers
* Helm charts

```bash
tar -czvf app.tgz app/
```

ZIP may break executable bits:

```bash
-rwxr-xr-x  →  -rw-r--r--
```

---

### ✅ Streaming (important in pipelines)

```bash
tar -cz app/ | docker build -
```

ZIP **cannot stream**.

---

### ✅ Tooling compatibility

* Helm → **only `.tgz`**
* Kubernetes images → tar layers
* Linux distros → tarballs
* GitHub Actions → tar preferred

---

## 🔹 When ZIP still makes sense

Use `zip` when:

* Sharing artifacts with **Windows users**
* Uploading to **email / Slack**
* Packaging **logs, reports, exports**
* Human-downloaded artifacts

❌ Avoid zip for:

* Docker images
* Helm charts
* Init containers
* Linux binaries

---

## 🔹 Decision Rule (simple)

> **Infra / CI / K8s / Docker → `tar.gz`**
> **Human download / cross-platform → `zip`**

---

# 🐳 ZIP Usage Inside Docker Pipelines

## 🔹 Installing zip in Docker images

### Debian / Ubuntu

```dockerfile
RUN apt-get update && apt-get install -y zip unzip
```

### Alpine (MOST COMMON)

```dockerfile
RUN apk add --no-cache zip unzip
```

---

## 🔹 Example: Zip build artifacts inside Docker

```dockerfile
FROM alpine:3.19

RUN apk add --no-cache zip

WORKDIR /app
COPY build/ .

RUN zip -r build-artifacts.zip .
```

---

## 🔹 CI pipeline example (GitHub Actions style)

```bash
zip -r artifacts.zip \
  manifests/ \
  values.yaml \
  Chart.yaml
```

Upload as artifact:

```yaml
- uses: actions/upload-artifact@v4
  with:
    name: build-artifacts
    path: artifacts.zip
```

---

## 🔹 Multi-stage Docker build (recommended)

```dockerfile
FROM node:20 AS build
WORKDIR /app
RUN npm run build

FROM alpine
RUN apk add --no-cache zip
COPY --from=build /app/dist /dist
RUN zip -r dist.zip /dist
```

Keeps final image clean.

---

# ☸️ ZIP Usage in Kubernetes Pipelines

## 🔹 Init container zipping logs or configs

```yaml
initContainers:
- name: zip-configs
  image: alpine
  command: ["sh", "-c"]
  args:
    - apk add --no-cache zip &&
      zip -r /out/configs.zip /configs
  volumeMounts:
  - name: configs
    mountPath: /configs
  - name: output
    mountPath: /out
```

---

## 🔹 Debug pod extracting ZIP

```bash
kubectl exec -it pod -- unzip configs.zip -d /tmp/configs
```

---

## 🔹 Common K8s mistake ❌

Using zip for init containers with scripts:

```bash
#!/bin/sh
# loses executable bit → fails
```

✅ Use tar instead:

```bash
tar -czvf scripts.tgz scripts/
```

---

# 🧠 DevOps Best Practices (Important)

### ✅ Prefer `tar.gz` for:

* Helm charts
* Docker layers
* Linux tools
* Kubernetes init scripts

### ⚠️ Use `zip` only for:

* Artifacts
* Reports
* Logs
* Cross-platform sharing

---

# 🧪 Cheat Commands (copy-paste)

```bash
# DevOps default
tar -czvf app.tgz app/

# CI artifact
zip -r artifacts.zip reports/

# Helm
helm package my-chart/

# Docker stream
tar -cz app/ | docker build -
```

---
