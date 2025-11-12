# ArgoCD Terraform Module for GKE

This module deploys ArgoCD to a Google Kubernetes Engine (GKE) cluster using Helm.

## Usage

This module is integrated into the GKE cluster deployment and can be enabled by setting the `enable_argocd` variable to `true`.

### Example in terraform.tfvars

```hcl
# Enable ArgoCD deployment
enable_argocd        = true
argocd_namespace     = "argocd"
argocd_chart_version = ""  # Leave empty for latest version
```

## Requirements

- Terraform >= 1.0
- helm provider ~> 2.13.2
- kubernetes provider ~> 2.30.0
- A running GKE cluster with kubectl access configured

## Notes

- The module depends on the GKE cluster being fully created and accessible
- The Helm and Kubernetes providers must be configured with valid GKE cluster credentials
- ArgoCD deployment is conditional based on the `enable_argocd` variable
