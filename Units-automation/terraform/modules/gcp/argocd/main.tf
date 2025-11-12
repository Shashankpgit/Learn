terraform {
  required_providers {
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.13.2"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.30.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6.0"
    }
  }
}

# Generate a random password for ArgoCD admin user
resource "random_password" "argocd_admin_password" {
  length  = 16
  special = true
  upper   = true
  lower   = true
  numeric = true
}

# Deploy ArgoCD using Helm chart
resource "helm_release" "argocd" {
  name             = var.argocd_release_name
  repository       = "https://argoproj.github.io/argo-helm"
  chart            = "argo-cd"
  namespace        = var.argocd_namespace
  create_namespace = true

  # Basic configuration values
  values = [
    yamlencode({
      configs = {
        secret = {
          # bcrypt-hashed value of the auto-generated password
          argocdServerAdminPassword = bcrypt(random_password.argocd_admin_password.result)
        }
      }
    })
  ]

  timeout = var.helm_timeout
  wait    = var.helm_wait

  depends_on = [random_password.argocd_admin_password]
}
