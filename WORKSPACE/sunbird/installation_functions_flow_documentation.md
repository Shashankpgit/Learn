# Sunbird-ED Installation Functions - Execution Flow Documentation

## Overview
This document explains the detailed execution flow and outcomes of the three critical functions that form the foundation of Sunbird-ED infrastructure setup:

1. `create_tf_backend`
2. `backup_configs` 
3. `create_tf_resources`

## Function 1: create_tf_backend

### What It Does
Creates Azure infrastructure required for Terraform remote state management.

### Execution Flow
```bash
create_tf_backend
```

### Step-by-Step Process
1. **Validates Prerequisites**
   - Checks `global-values.yaml` exists
   - Verifies `yq` and `jq` tools are installed
   - Confirms Azure CLI authentication

2. **Extracts Configuration**
   - Reads `building_block`, `environment_name`, `location` from global-values.yaml
   - Gets Azure tenant ID and subscription ID

3. **Creates Azure Resources**
   - **Resource Group**: `{building_block}-{environment}`
   - **Storage Account**: `{environment}tfstate{tenant_id_prefix}`
   - **Blob Container**: `{environment}tfstate`

4. **Generates Environment File**
   - Creates `tf.sh` with backend configuration variables

### What You'll See After Execution

#### Console Output
```bash
Extracted building_block: "sunbird"
Extracted environment_name: "demo"
RESOURCE_GROUP_NAME: sunbird-demo
STORAGE_ACCOUNT_NAME: demotfstate12345678
CONTAINER_NAME: demotfstate
SUBSCRIPTION_ID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx

Terraform backend setup complete!
Run the following command to set the environment variables:
source tf.sh
```

#### Created Azure Resources
- **Resource Group**: `sunbird-demo`
- **Storage Account**: `demotfstate12345678`
  - SKU: Standard_LRS
  - Encryption: Enabled
  - Location: As specified in global-values.yaml
- **Blob Container**: `demotfstate`

#### Generated Files
- **tf.sh**: Contains environment variables for Terraform backend
```bash
export AZURE_TERRAFORM_BACKEND_RG=sunbird-demo
export AZURE_TERRAFORM_BACKEND_STORAGE_ACCOUNT=demotfstate12345678
export AZURE_TERRAFORM_BACKEND_CONTAINER=demotfstate
export AZURE_SUBSCRIPTION_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

### Impact on System
- Azure subscription now has dedicated resources for Terraform state
- All subsequent Terraform operations will use this remote backend
- State files will be centrally stored and locked for team collaboration

---

## Function 2: backup_configs

### What It Does
Backs up existing configuration files that will be overwritten during installation.

### Execution Flow
```bash
backup_configs
```

### Step-by-Step Process
1. **Creates Timestamp**
   ```bash
   timestamp=$(date +%d%m%y_%H%M%S)  # Format: 151224_143022
   ```

2. **Backs Up Kubernetes Config**
   ```bash
   mkdir -p ~/.kube
   mv ~/.kube/config ~/.kube/config.$timestamp || true
   ```

3. **Backs Up Rclone Config**
   ```bash
   mkdir -p ~/.config/rclone
   mv ~/.config/rclone/rclone.conf ~/.config/rclone/rclone.conf.$timestamp || true
   ```

4. **Sets Kubernetes Environment**
   ```bash
   export KUBECONFIG=~/.kube/config
   ```

### What You'll See After Execution

#### Console Output
```bash
Backup existing config files if they exist
```

#### File System Changes
**Before Execution:**
```
~/.kube/config                    # Existing Kubernetes config
~/.config/rclone/rclone.conf      # Existing Rclone config
```

**After Execution:**
```
~/.kube/config.151224_143022      # Backed up Kubernetes config
~/.config/rclone/rclone.conf.151224_143022  # Backed up Rclone config
~/.kube/                          # Empty directory (ready for new config)
~/.config/rclone/                 # Empty directory (ready for new config)
```

#### Environment Variables Set
```bash
KUBECONFIG=~/.kube/config
```

### Impact on System
- Existing configurations are safely preserved with timestamp
- Directories are prepared for new configurations
- No data loss occurs during installation
- Previous configurations can be restored if needed

---

## Function 3: create_tf_resources

### What It Does
Creates all Azure infrastructure resources using Terraform and Terragrunt.

### Execution Flow
```bash
create_tf_resources
```

### Step-by-Step Process
1. **Loads Backend Configuration**
   ```bash
   source tf.sh  # Loads environment variables from previous step
   ```

2. **Initializes Terraform**
   ```bash
   terraform init -upgrade    # Downloads providers and modules
   terragrunt init -upgrade   # Initializes Terragrunt configuration
   ```

3. **Creates Infrastructure**
   ```bash
   terragrunt run-all apply --terragrunt-non-interactive
   ```

4. **Secures Kubernetes Config**
   ```bash
   chmod 600 ~/.kube/config
   ```

### What You'll See After Execution

#### Console Output
```bash
Creating resources on azure cloud

Initializing the backend...
Initializing provider plugins...
- Downloading plugin for provider "azurerm" (hashicorp/azurerm) 4.0.1...
- Downloading plugin for provider "azuread" (hashicorp/azuread) 2.x.x...

Terragrunt will perform the following actions:
  + Create multiple resources across modules

Apply complete! Resources: 45 added, 0 changed, 0 destroyed.
```

#### Created Azure Infrastructure

**1. Network Resources**
- **Virtual Network**: `sunbird-demo`
  - Address space: 10.0.0.0/16
  - Location: As specified in global-values.yaml
- **AKS Subnet**: `sunbird-demo-aks`
  - Address prefix: 10.0.1.0/24
  - Service endpoints: Storage, KeyVault, SQL

**2. AKS Cluster Resources**
- **Azure AD Application**: `sunbird-demo`
- **Service Principal**: For AKS authentication
- **Role Assignment**: Network Contributor on subnet
- **AKS Cluster**: `sunbird-demo`
  - Node pool: Default with specified VM size/count
  - Network plugin: Azure CNI
  - DNS prefix: `sunbird-demo`
  - Service CIDR: 10.1.0.0/16
  - DNS service IP: 10.1.0.10

**3. Storage Resources**
- **Primary Storage Account**: For application data
- **Backup Storage Account**: For backups and logs
- **Blob Containers**: 
  - Content storage
  - Telemetry data
  - Backup storage
  - Asset storage

**4. Security Resources**
- **Encryption Keys**: For data encryption
- **Certificates**: SSL/TLS certificates
- **Secrets**: Database passwords, API keys
- **Access Policies**: For secure resource access

**5. Random Passwords**
- Database passwords
- Service account passwords
- API tokens
- Encryption keys

#### Generated Files

**1. Kubernetes Configuration**
- **~/.kube/config**: Complete cluster access configuration
  - Cluster endpoint
  - Authentication certificates
  - User credentials
  - Context settings

**2. Terraform State Files** (in Azure Storage)
```
Container: demotfstate/
├── random_passwords/terraform.tfstate
├── network/terraform.tfstate
├── aks/terraform.tfstate
├── storage/terraform.tfstate
├── keys/terraform.tfstate
├── upload-files/terraform.tfstate
└── output-file/terraform.tfstate
```

**3. Configuration Files**
- **global-cloud-values.yaml**: Cloud-specific configurations
- **global-keys-values.yaml**: Generated keys and secrets

#### Resource Dependencies Created
```
Random Passwords → Network → AKS → Storage → Keys → Upload Files → Output File
```

### Impact on System

#### Local System Changes
- **Kubernetes Access**: `kubectl` commands now work with new cluster
- **Configuration Files**: New configs replace backed-up versions
- **Environment**: Ready for Helm chart deployments

#### Azure Cloud Changes
- **Complete Infrastructure**: All required resources for Sunbird-ED
- **Networking**: Secure virtual network with proper subnets
- **Compute**: Kubernetes cluster ready for application deployment
- **Storage**: Multiple storage accounts for different data types
- **Security**: Encryption, certificates, and access controls in place

#### Verification Commands
After execution, you can verify:
```bash
# Check Kubernetes cluster access
kubectl cluster-info

# List nodes
kubectl get nodes

# Check Azure resources
az resource list --resource-group sunbird-demo

# Verify storage accounts
az storage account list --resource-group sunbird-demo
```

---

## Complete Execution Sequence Impact

### Before Any Function Execution
- Clean system with potential existing configs
- Azure subscription with basic access
- No Sunbird-ED infrastructure

### After create_tf_backend
- Azure backend infrastructure ready
- Terraform state management configured
- Environment variables available

### After backup_configs
- Existing configurations safely backed up
- System prepared for new configurations
- No risk of data loss

### After create_tf_resources
- Complete Azure infrastructure deployed
- Kubernetes cluster operational
- Storage accounts and security configured
- System ready for application deployment

### Resource Summary
**Total Azure Resources Created**: ~45 resources including:
- 1 Resource Group
- 1 Virtual Network + Subnet
- 1 AKS Cluster with node pool
- 3-5 Storage Accounts
- Multiple Blob Containers
- Security keys and certificates
- Service principals and role assignments

### Next Steps After These Functions
The system is now ready for:
1. Helm chart installations (`install_helm_components`)
2. Application deployments
3. DNS configuration
4. SSL certificate setup
5. Post-installation validation

This completes the infrastructure foundation for Sunbird-ED platform deployment.