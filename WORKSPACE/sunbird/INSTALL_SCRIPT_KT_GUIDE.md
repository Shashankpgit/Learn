# Sunbird-ED install.sh Script - Knowledge Transfer Guide

## Overview for New Joiners

Hey there! Welcome to the team. This document will help you understand what the `install.sh` script does in the Sunbird-ED deployment process. Think of this script as the **master orchestrator** that takes you from zero to a fully functional Sunbird-ED platform.

**Location**: `terraform/azure/<environment>/install.sh`

---

## What Does install.sh Actually Do?

The script is like a **recipe book** that follows a specific sequence to deploy Sunbird-ED. When you run `./install.sh`, it executes **11 sequential steps** automatically.

### The Complete Journey (Default Execution)

```bash
./install.sh  # No arguments = Full deployment
```

Here's what happens step by step:

---

## Phase 1: Infrastructure Setup (Steps 1-3)

### Step 1: `create_tf_backend`
**What it does**: Creates the "storage locker" for Terraform state files
```bash
# Runs: create_tf_backend.sh
```
**Real Impact**: 
- Creates Azure Storage Account to store Terraform state
- Sets up blob container for state management
- **Why needed**: Terraform needs to remember what it created

### Step 2: `backup_configs` 
**What it does**: Backs up your existing configuration files
```bash
# Backs up:
# ~/.kube/config → ~/.kube/config.TIMESTAMP
# ~/.config/rclone/rclone.conf → ~/.config/rclone/rclone.conf.TIMESTAMP
```
**Real Impact**: 
- Protects your existing Kubernetes and rclone configs
- Creates timestamped backups before overwriting

### Step 3: `create_tf_resources`
**What it does**: **CREATES THE ENTIRE AZURE INFRASTRUCTURE**
```bash
# Runs: terraform init, terragrunt init, terragrunt run-all apply
```
**Real Impact**: 
- Creates AKS (Azure Kubernetes Service) cluster
- Sets up Virtual Networks, subnets, security groups
- Creates storage accounts, load balancers
- Provisions all Azure resources needed
- **Duration**: 15-20 minutes typically

---

## Phase 2: Application Deployment (Step 4)

### Step 4: `install_helm_components`
**What it does**: **DEPLOYS ALL SUNBIRD SERVICES TO KUBERNETES**

**The Magic Sequence**:
```bash
components=("monitoring" "edbb" "learnbb" "knowledgebb" "obsrvbb" "inquirybb" "additional")
```

**What each component contains**:

#### 4a. **monitoring** (Infrastructure monitoring)
- Prometheus (metrics collection)
- Grafana (dashboards) 
- Loki (log aggregation)
- AlertManager (notifications)

#### 4b. **edbb** (Core Education Building Block)
- nginx-public-ingress (external load balancer)
- nginx-private-ingress (internal routing)
- kong (API gateway)
- player (content player)
- nodebb (discussion forums)
- bot (chatbot service)

#### 4c. **learnbb** (Learning Management)
- lms (Learning Management System)
- userorg (user/organization management)
- groups (group management)
- notification (notification service)
- cert (certificate generation)
- registry (registry service)

#### 4d. **knowledgebb** (Content Management)
- content (content service)
- learning (learning service) 
- search (search functionality)
- taxonomy (content categorization)
- dial (DIAL code services - optional)

#### 4e. **obsrvbb** (Observability)
- analytics (analytics service)
- telemetry (telemetry collection)
- superset (data visualization)

#### 4f. **inquirybb** (Assessment)
- assessment (assessment service)
- flink (stream processing jobs)

#### 4g. **additional** (Supporting services)
- Additional utility services

**Real Impact**: 
- Deploys 50+ microservices to Kubernetes
- Sets up databases (PostgreSQL, Cassandra, Redis, Neo4j)
- Configures message queues (Kafka)
- **Duration**: 30-45 minutes typically

---

## Phase 3: Post-Deployment Configuration (Steps 5-11)

### Step 5: `post_install_nodebb_plugins`
**What it does**: Configures the discussion forum
```bash
# Activates NodeBB plugins:
# - nodebb-plugin-create-forum
# - nodebb-plugin-sunbird-oidc  
# - nodebb-plugin-write-api
```
**Real Impact**: Enables forum functionality with SSO integration

### Step 6: `restart_workloads_using_keys`
**What it does**: Restarts services to pick up security keys
```bash
# Restarts: neo4j, knowledge-mw, player, report, content, adminutil, 
#          cert-registry, groups, userorg, lms, notification, registry, analytics
```
**Real Impact**: Ensures all services have proper authentication keys

### Step 7: `certificate_config`
**What it does**: Configures certificate signing keys
```bash
# Generates RSA key pairs for certificate signing
# Injects public keys into registry service
```
**Real Impact**: Enables certificate generation for users

### Step 8: `dns_mapping`
**What it does**: **WAITS FOR DNS PROPAGATION**
```bash
# Gets load balancer IP from: kubectl get svc nginx-public-ingress
# Waits up to 20 minutes for DNS to propagate
```
**Real Impact**: 
- Shows you the IP address to configure in your DNS
- Waits for your domain to resolve to the correct IP
- **Action Required**: You need to update your DNS A record

### Step 9: `generate_postman_env`
**What it does**: Creates API testing environment
```bash
# Creates: env.json file with all API credentials
```
**Real Impact**: 
- Extracts API keys, passwords, domain info from Kubernetes
- Creates Postman environment file for API testing

### Step 10: `run_post_install`
**What it does**: **CONFIGURES THE PLATFORM VIA APIs**
```bash
# Runs: postman collection run collection${RELEASE}.json
```
**Real Impact**: 
- Creates default users (admin, content creators, etc.)
- Sets up initial platform configuration
- Configures system settings via API calls

### Step 11: `create_client_forms`
**What it does**: Creates client-side forms and configurations
```bash
# Runs multiple Postman collections for form creation
```
**Real Impact**: 
- Sets up editor forms
- Configures mobile app settings
- Creates question set editor configurations

---

## Individual Function Usage

You can also run individual functions for troubleshooting:

```bash
# Infrastructure only
./install.sh create_tf_resources

# Applications only  
./install.sh install_helm_components

# Specific component only
./install.sh install_component monitoring

# Post-install only
./install.sh run_post_install

# DNS check only
./install.sh dns_mapping

# Generate API environment
./install.sh generate_postman_env

# Destroy everything
./install.sh destroy_tf_resources
```

---

## Key Files the Script Uses

### Configuration Files:
- `global-values.yaml` - Your environment configuration
- `global-cloud-values.yaml` - Cloud-specific settings (auto-generated)
- `postman.env.json` - API testing template

### Generated Files:
- `env.json` - Postman environment with real credentials
- `certkey.pem` / `certpubkey.pem` - Certificate signing keys
- `~/.kube/config` - Kubernetes cluster access

### External Dependencies:
- `../../../helmcharts/` - All Helm charts
- `../../../postman-collection/` - API configuration collections
- `create_tf_backend.sh` - Terraform backend setup
- `tf.sh` - Terraform environment variables

---

## What to Expect During Execution

### Timeline:
- **Total Duration**: 60-90 minutes for full deployment
- **Infrastructure**: 15-20 minutes
- **Applications**: 30-45 minutes  
- **Configuration**: 15-30 minutes

### Success Indicators:
```bash
# You'll see these success messages:
"All pods are running successfully"
"DNS mapping has propagated successfully"
"NodeBB plugins are activated, built, and NodeBB has been restarted"
```

### Default Users Created:
| Role | Email | Password |
|------|-------|----------|
| Admin | admin@yopmail.com | Admin@123 |
| Content Creator | contentcreator@yopmail.com | Creator@123 |
| Content Reviewer | contentreviewer@yopmail.com | Reviewer@123 |

---

## Common Scenarios for New Joiners

### Scenario 1: Fresh Deployment
```bash
# You have: domain, SSL cert, Azure subscription
cd terraform/azure/demo
./install.sh  # Sit back and wait 60-90 minutes
```

### Scenario 2: Something Failed Mid-Way
```bash
# Check what failed, then resume from specific step
./install.sh dns_mapping          # If DNS failed
./install.sh run_post_install     # If API config failed
```

### Scenario 3: Need to Redeploy Just Applications
```bash
./install.sh install_helm_components  # Skips infrastructure
```

### Scenario 4: Complete Cleanup
```bash
./install.sh destroy_tf_resources  # Destroys everything in Azure
```

---

## Pro Tips for New Joiners

1. **Always check prerequisites first** - Domain, SSL cert, Azure access
2. **Monitor the logs** - The script is verbose, read the output
3. **DNS is critical** - Step 8 will fail if DNS isn't configured
4. **Be patient** - Full deployment takes 60-90 minutes
5. **Use individual functions** for troubleshooting
6. **Keep backups** - Script backs up configs automatically

---

## What Happens Behind the Scenes

The script is essentially running these tools in sequence:
1. **Terraform/Terragrunt** → Creates Azure infrastructure
2. **Helm** → Deploys Kubernetes applications  
3. **kubectl** → Manages Kubernetes resources
4. **Postman CLI** → Configures platform via APIs

Think of it as an **automated DevOps engineer** that knows exactly how to deploy Sunbird-ED from scratch!

---

## Quick Reference Commands

```bash
# Full deployment
./install.sh

# Check what's running
kubectl get pods -n sunbird

# Get load balancer IP  
kubectl get svc -n sunbird nginx-public-ingress

# Check logs of a service
kubectl logs -n sunbird deployment/lms

# Access your platform
https://your-domain.com
```

That's it! You now understand what the install.sh script does. It's your one-stop-shop for deploying the entire Sunbird-ED platform. Any questions? Feel free to ask the team!