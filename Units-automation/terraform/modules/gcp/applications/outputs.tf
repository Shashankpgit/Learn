output "applications_release_name" {
  description = "Name of the applications Helm release"
  value       = helm_release.applications.name
}