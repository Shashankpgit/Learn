output "namespace" {
  description = "The namespace where ArgoCD is deployed"
  value       = helm_release.argocd.namespace
}

output "release_name" {
  description = "The Helm release name for ArgoCD"
  value       = helm_release.argocd.name
}

output "release_status" {
  description = "The status of the ArgoCD Helm release"
  value       = helm_release.argocd.status
}

output "chart_version" {
  description = "The version of the ArgoCD chart deployed"
  value       = helm_release.argocd.version
}

output "argocd_server_service_name" {
  description = "The name of the ArgoCD server service"
  value       = "${helm_release.argocd.name}-server"
}

output "admin_username" {
  description = "ArgoCD admin username"
  value       = "admin"
}

output "admin_password" {
  description = "ArgoCD admin password (auto-generated)"
  value       = random_password.argocd_admin_password.result
}
