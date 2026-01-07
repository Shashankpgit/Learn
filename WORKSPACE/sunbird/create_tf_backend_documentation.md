# Terraform Backend Setup Script - Knowledge Transfer Document

## Script Overview
**File**: `create_tf_backend.sh`  
**Purpose**: Creates Azure infrastructure for Terraform remote state management  
**Execution Context**: First step in Sunbird-ED installation process  

## What This Script Does

### Core Function
Sets up secure, centralized storage for Terraform state files in Azure, enabling:
- Team collaboration on infrastructure
- State locking to prevent conflicts
- Encrypted state storage
- Disaster recovery for state files

### Prerequisites Validation
1. **File Dependencies**:
   - `global-values.yaml` must exist in current directory
   - Contains environment configuration values

2. **Tool Dependencies**:
   - `yq` - YAML processor for parsing configuration
   - `jq` - JSON processor for Azure CLI output
   - `az` - Azure CLI (must be authenticated)

3. **Azure Authentication**:
   - User must be logged in: `az login --tenant TENANT_ID`
   - Requires permissions to create resource groups and storage accounts

## Step-by-Step Execution Flow

### Step 1: Environment Validation
```bash
set -euo pipefail  # Strict error handling
```
- **`-e`**: Exit on any command failure
- **`-u`**: Exit on undefined variables
- **`-o pipefail`**: Exit on pipe command failures

### Step 2: Configuration Extraction
```bash
building_block=$(yq '.global.building_block' global-values.yaml)
environment_name=$(yq '.global.environment' global-values.yaml)
location=$(yq '.global.cloud_storage_region' global-values.yaml)
```

**Expected Values Example**:
- `building_block`: "sunbird"
- `environment_name`: "demo"
- `location`: "East US"

### Step 3: Azure Context Retrieval
```bash
ID=$(az account show | jq -r .tenantId | cut -d '-' -f1)
SUBSCRIPTION_ID=$(az account show | jq -r .id)
```

**What Happens**:
- Extracts first 8 characters of Azure Tenant ID
- Gets full Azure Subscription ID
- Uses current authenticated Azure context

### Step 4: Resource Naming Convention
```bash
RESOURCE_GROUP_NAME="${building_block}-${environment_name}"
STORAGE_ACCOUNT_NAME="${environment_name}tfstate$ID"
CONTAINER_NAME="${environment_name}tfstate"
```

**Naming Pattern**:
- **Resource Group**: `{building_block}-{environment}`
- **Storage Account**: `{environment}tfstate{tenant_id_prefix}`
- **Container**: `{environment}tfstate`

**Example Output**:
- Resource Group: `sunbird-demo`
- Storage Account: `demotfstate12345678`
- Container: `demotfstate`

### Step 5: Azure Resource Creation

#### 5.1 Resource Group Creation
```bash
az group create --name "$RESOURCE_GROUP_NAME" --location "$location"
```
- Creates logical container for all Terraform backend resources
- Uses location from global-values.yaml

#### 5.2 Storage Account Creation
```bash
az storage account create --resource-group "$RESOURCE_GROUP_NAME" \
  --name "$STORAGE_ACCOUNT_NAME" --sku Standard_LRS --encryption-services blob
```

**Storage Account Specifications**:
- **SKU**: Standard_LRS (Locally Redundant Storage)
- **Replication**: 3 copies within same datacenter
- **Encryption**: Blob encryption enabled by default
- **Performance**: Standard tier (sufficient for state files)

#### 5.3 Blob Container Creation
```bash
az storage container create --name "$CONTAINER_NAME" --account-name "$STORAGE_ACCOUNT_NAME"
```
- Creates container within storage account
- Will hold individual Terraform state files for each module

### Step 6: Environment Variables Export
```bash
echo "export AZURE_TERRAFORM_BACKEND_RG=$RESOURCE_GROUP_NAME" > tf.sh
echo "export AZURE_TERRAFORM_BACKEND_STORAGE_ACCOUNT=$STORAGE_ACCOUNT_NAME" >> tf.sh
echo "export AZURE_TERRAFORM_BACKEND_CONTAINER=$CONTAINER_NAME" >> tf.sh
echo "export AZURE_SUBSCRIPTION_ID=$SUBSCRIPTION_ID" >> tf.sh
```

**Generated tf.sh File**:
```bash
export AZURE_TERRAFORM_BACKEND_RG=sunbird-demo
export AZURE_TERRAFORM_BACKEND_STORAGE_ACCOUNT=demotfstate12345678
export AZURE_TERRAFORM_BACKEND_CONTAINER=demotfstate
export AZURE_SUBSCRIPTION_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

## How Terragrunt Uses These Variables

### Backend Configuration Generation
Terragrunt uses these environment variables to generate `backend.tf`:
```hcl
terraform {
  backend "azurerm" {
    resource_group_name  = "${get_env("AZURE_TERRAFORM_BACKEND_RG")}"
    storage_account_name = "${get_env("AZURE_TERRAFORM_BACKEND_STORAGE_ACCOUNT")}"
    container_name       = "${get_env("AZURE_TERRAFORM_BACKEND_CONTAINER")}"
    key                  = "${path_relative_to_include()}/terraform.tfstate"
  }
}
```

### State File Organization
Each Terraform module gets its own state file:
```
Container: demotfstate/
├── network/terraform.tfstate
├── aks/terraform.tfstate
├── storage/terraform.tfstate
├── keys/terraform.tfstate
├── random_passwords/terraform.tfstate
├── upload-files/terraform.tfstate
└── output-file/terraform.tfstate
```

## Error Scenarios and Troubleshooting

### Common Error 1: Missing global-values.yaml
```bash
Error: global-values.yaml file does not exist!
```
**Solution**: Ensure you're in the correct directory with global-values.yaml

### Common Error 2: yq not installed
```bash
Error: yq is not installed. Please install yq to process YAML files.
```
**Solution**: Install yq:
```bash
sudo wget -qO /usr/local/bin/yq https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64
sudo chmod +x /usr/local/bin/yq
```

### Common Error 3: Azure CLI not authenticated
```bash
az: command not found
# or
ERROR: Please run 'az login' to setup account.
```
**Solution**: Authenticate with Azure:
```bash
az login --tenant YOUR_TENANT_ID
```

### Common Error 4: Storage account name conflict
```bash
The storage account name 'demotfstate12345678' is already taken.
```
**Solution**: Storage account names are globally unique. The script uses tenant ID to minimize conflicts, but if it occurs:
- Change environment name in global-values.yaml
- Or manually modify the generated name

### Common Error 5: Insufficient permissions
```bash
The client does not have authorization to perform action 'Microsoft.Resources/subscriptions/resourceGroups/write'
```
**Solution**: Ensure Azure account has required permissions:
- Resource Group Contributor
- Storage Account Contributor

## Security Considerations

### Data Protection
- **Encryption at Rest**: Storage account encrypts all blobs by default
- **Encryption in Transit**: HTTPS enforced for all operations
- **Access Control**: Azure RBAC controls who can access state files

### State Locking
- Azure Storage provides automatic state locking
- Prevents concurrent Terraform operations
- Uses blob leases for locking mechanism

### Audit Trail
- All operations logged in Azure Activity Log
- Resource creation/modification tracked
- Access patterns monitored

## Integration with Main Installation Process

### Call Sequence in install.sh
```bash
create_tf_backend     # This script
backup_configs        # Backup existing configs
create_tf_resources   # Create infrastructure
# ... rest of installation
```

### Environment Variable Usage
After script execution:
```bash
source tf.sh  # Load environment variables
terragrunt run-all plan   # Uses backend configuration
terragrunt run-all apply  # Creates infrastructure
```

## Best Practices for New Team Members

### Before Running Script
1. **Verify Azure Authentication**: `az account show`
2. **Check global-values.yaml**: Ensure correct values
3. **Validate Tools**: Confirm yq, jq, az are installed
4. **Review Permissions**: Ensure sufficient Azure access

### During Execution
1. **Monitor Output**: Watch for any error messages
2. **Verify Resource Creation**: Check Azure portal
3. **Validate tf.sh**: Ensure environment variables are correct

### After Execution
1. **Source Variables**: `source tf.sh`
2. **Test Backend**: Run `terragrunt validate` in any module
3. **Backup tf.sh**: Keep copy of environment variables
4. **Document Changes**: Note any customizations made

### Cleanup (if needed)
To remove backend resources:
```bash
az storage container delete --name $CONTAINER_NAME --account-name $STORAGE_ACCOUNT_NAME
az storage account delete --name $STORAGE_ACCOUNT_NAME --resource-group $RESOURCE_GROUP_NAME
az group delete --name $RESOURCE_GROUP_NAME
```

## Key Takeaways for New Joiners

1. **This script is mandatory** - Must run before any Terraform operations
2. **One-time setup** - Only needs to run once per environment
3. **Environment-specific** - Each environment gets its own backend
4. **Team shared** - All team members use same backend for collaboration
5. **State critical** - Backend stores all infrastructure state information

## Support and Escalation

### When to Escalate
- Persistent authentication issues
- Storage account naming conflicts
- Permission-related errors
- Azure service outages

### Debug Information to Collect
- Azure subscription ID
- Tenant ID
- Error messages (full output)
- global-values.yaml content (sanitized)
- Azure CLI version: `az --version`

This document provides comprehensive understanding of the Terraform backend setup process for successful onboarding of new team members.