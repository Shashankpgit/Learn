#############################
# Outputs
#############################

output "instances_self_links" {
  description = "List of self-links for compute instances"
  value       = google_compute_instance.this[*].self_link
}

output "instances_details" {
  description = "List of all details for compute instances"
  sensitive   = true
  value       = google_compute_instance.this[*]
}

output "instance_names" {
  description = "List of instance names"
  value       = google_compute_instance.this[*].name
}

output "instance_zones" {
  description = "List of instance zones"
  value       = google_compute_instance.this[*].zone
}

output "instance_internal_ip" {
  description = "Internal IP address of the VM instance"
  value       = google_compute_instance.this.network_interface[0].network_ip
}

output "instance_external_ip" {
  description = "External IP address of the VM instance"
  value       = google_compute_instance.this.network_interface[0].access_config[0].nat_ip
}