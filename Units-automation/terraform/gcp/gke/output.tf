output "cluster_name" {
  description = "The name of the GKE cluster"
  value       = module.gke_cluster.name
}

# ArgoCD outputs (only when enabled)
output "argocd_enabled" {
  description = "Whether ArgoCD is enabled"
  value       = var.enable_argocd
}

output "argocd_admin_username" {
  description = "ArgoCD admin username"
  value       = var.enable_argocd ? "admin" : null
}

output "argocd_admin_password" {
  description = "ArgoCD admin password (auto-generated)"
  value       = var.enable_argocd ? nonsensitive(module.argocd[0].admin_password) : null
}

output "service_credentials_file" {
  description = "Path to the service credentials JSON file"
  value       = var.enable_argocd ? local_file.service_credentials[0].filename : null
}

output "applications_release_name" {
  description = "Name of the applications Helm release"
  value       = var.enable_argocd && var.enable_applications ? module.applications[0].applications_release_name : null
}