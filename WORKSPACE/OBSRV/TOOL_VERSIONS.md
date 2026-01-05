# Tool Versions Used in Obsrv

This document lists the versions of tools used in the Obsrv platform and where they are specified in the codebase.

## Required Tool Versions

### Helm
- **Version**: `>=3.10.2`
- **Location**: 
  - `infra-setup/obsrv.md` (line 30)
  - `infra-setup/obsrv.sh` (line 96, 104, 127)
- **Installation**: The script installs Helm v3.10.2 if not present
- **Download URL**: `https://get.helm.sh/helm-v3.10.2-linux-amd64.tar.gz`

### kubectl
- **Version**: `1.32.0-r1` (Alpine package version)
- **Locations**:
  - **Dockerfile**: `Dockerfiles/kubectl/Dockerfile` (line 2)
    ```dockerfile
    RUN apk add --no-cache curl kubectl>=1.32.0-r1 jq openssl postgresql-client
    ```
  - **Service Values Files**:
    - `helmcharts/services/dataset-api/values.yaml` (line 20)
    - `helmcharts/services/config-api/values.yaml` (line 20)
    - `helmcharts/services/web-console/values.yaml` (line 19)
    - `helmcharts/services/keycloak/values.yaml` (line 1378)
  - **Bootstrapper**: `helmcharts/bootstrapper/templates/secrets.yaml` (line 61)
    ```yaml
    image: sanketikahub/kubectl:1.32.0-r1
    ```
- **Note**: kubectl is used as a container image (`sanketikahub/kubectl:1.32.0-r1`) in init jobs and bootstrapper tasks

### AWS CLI
- **Version**: `>=2.13.8`
- **Locations**:
  - `infra-setup/obsrv.md` (line 29)
  - `infra-setup/obsrv.sh` (line 95, 118-119)
- **Installation**: The script installs AWS CLI v2.13.8 if not present
- **Download URL**: `https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip`

### Terraform
- **Version**: `>=1.5.7`
- **Locations**:
  - `infra-setup/obsrv.md` (line 31)
  - `infra-setup/obsrv.sh` (line 97, 105, 130-131)
- **Installation**: The script installs Terraform v1.5.7 if not present
- **Download URL**: `https://releases.hashicorp.com/terraform/1.5.7/terraform_1.5.7_linux_amd64.zip`

### Terragrunt
- **Version**: `>=0.45.6` (specified), but installs `0.45.8` (latest)
- **Locations**:
  - `infra-setup/obsrv.md` (line 33)
  - `infra-setup/obsrv.sh` (line 99, 107, 138-139)
- **Installation**: The script installs Terragrunt v0.45.8 if not present
- **Download URL**: `https://github.com/gruntwork-io/terragrunt/releases/download/v0.45.8/terragrunt_linux_amd64`

### Terrahelp
- **Version**: `>=0.7.5` (specified), but installs `0.4.3` (available)
- **Locations**:
  - `infra-setup/obsrv.md` (line 32)
  - `infra-setup/obsrv.sh` (line 98, 106, 134-135)
- **Installation**: The script installs Terrahelp v0.4.3 if not present
- **Download URL**: `https://github.com/opencredo/terrahelp/releases/download/v0.4.3/terrahelp-linux-amd64`

### GCloud CLI (for GCP)
- **Version**: `>=474.0.0`
- **Locations**:
  - `infra-setup/obsrv.sh` (line 103, 122-123)
- **Installation**: The script installs GCloud CLI v474.0.0 if not present
- **Download URL**: `https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-474.0.0-linux-x86_64.tar.gz`

---

## Summary Table

| Tool | Minimum Version | Installed Version | Location |
|------|----------------|-------------------|----------|
| **helm** | >=3.10.2 | 3.10.2 | `infra-setup/obsrv.md`, `infra-setup/obsrv.sh` |
| **kubectl** | 1.32.0-r1 | 1.32.0-r1 | `Dockerfiles/kubectl/Dockerfile`, service values.yaml files |
| **aws** | >=2.13.8 | 2.13.8 | `infra-setup/obsrv.md`, `infra-setup/obsrv.sh` |
| **terraform** | >=1.5.7 | 1.5.7 | `infra-setup/obsrv.md`, `infra-setup/obsrv.sh` |
| **terragrunt** | >=0.45.6 | 0.45.8 | `infra-setup/obsrv.md`, `infra-setup/obsrv.sh` |
| **terrahelp** | >=0.7.5 | 0.4.3 | `infra-setup/obsrv.md`, `infra-setup/obsrv.sh` |
| **gcloud** | >=474.0.0 | 474.0.0 | `infra-setup/obsrv.sh` (GCP only) |

---

## Version Check Commands

To check your installed versions:

```bash
# Helm
helm version --short

# kubectl
kubectl version --client

# AWS CLI
aws --version

# Terraform
terraform version

# Terragrunt
terragrunt --version

# GCloud (for GCP)
gcloud version
```

---

## Notes

1. **kubectl**: The version `1.32.0-r1` refers to the Alpine Linux package version. This corresponds to Kubernetes client version 1.32.0.

2. **Version Mismatch**: There's a discrepancy for `terrahelp`:
   - Required: `>=0.7.5`
   - Installed: `0.4.3` (which is actually older)
   - This might need to be addressed if terrahelp 0.7.5+ is actually required

3. **Terragrunt**: The script installs v0.45.8 even though the requirement is >=0.45.6, which is fine as it meets the minimum requirement.

4. **kubectl as Container Image**: kubectl is primarily used as a container image (`sanketikahub/kubectl:1.32.0-r1`) in Kubernetes init jobs rather than as a CLI tool on the host system.

5. **Installation Script**: The `infra-setup/obsrv.sh` script can automatically install these tools if `--install_dependencies true` is passed.

