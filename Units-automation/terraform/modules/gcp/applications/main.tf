terraform {
  required_providers {
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.13.2"
    }
  }
}

resource "helm_release" "applications" {
  name       = var.applications_release_name
  chart      = var.applications_chart_path
  namespace  = var.argocd_namespace
  create_namespace = true

  values = [
    yamlencode({
      config = {
        spec = {
          source = {
            repoURL        = var.applications_repo_url
            targetRevision = var.applications_target_revision
            path           = var.applications_path
          }
        }
      }
      applications = [
        {
          name      = "root-application"
          namespace = "argocd"
          path      = var.applications_path
          tool = {
            helm = {
              releaseName = "root-application"
            }
          }
        }
      ]
    })
  ]

  timeout = var.helm_timeout
  wait    = var.helm_wait

  depends_on = [var.depends_on_argocd]
}