#############################
# Project & Basic Configuration
#############################

variable "project_id" {
  description = "The GCP project ID"
  type        = string
}

variable "vm_name" {
  description = "Name for the VM instance (used as hostname base)"
  type        = string
  default     = "demo-vm"
}

#############################
# Compute Configuration
#############################

variable "vm_zone" {
  description = "GCP zone for the VM"
  type        = string
  default     = "asia-south1-a"
}

variable "vm_machine_type" {
  description = "GCE machine type"
  type        = string
  default     = "e2-medium"
}

variable "vm_image" {
  description = "Boot image for VM (image family or full path)"
  type        = string
  default     = "projects/ubuntu-os-cloud/global/images/family/ubuntu-2404-lts"
}

variable "vm_disk_size" {
  description = "Boot disk size in GB"
  type        = number
  default     = 10
}

variable "vm_disk_type" {
  description = "Boot disk type (pd-standard, pd-ssd, pd-balanced)"
  type        = string
  default     = "pd-standard"
}

variable "auto_delete_boot_disk" {
  description = "Whether to auto-delete the boot disk when the VM is deleted"
  type        = bool
  default     = true
}

#############################
# Network Configuration
#############################

variable "vm_network_self_link" {
  description = "Self link of the VPC network for the VM"
  type        = string
}

variable "vm_subnetwork_self_link" {
  description = "Self link of the subnetwork for the VM"
  type        = string
}
variable "vm_network_tag" {
  description = "Network tag for VM instance (for firewall rules)"
  type        = string
}

variable "enable_public_ip" {
  description = "Whether to assign a public IP to the VM"
  type        = bool
  default     = true
}

#############################
# Scheduling & Lifecycle
#############################

variable "preemptible" {
  description = "Whether the VM is preemptible"
  type        = bool
  default     = false
}

variable "automatic_restart" {
  description = "Whether the VM should automatically restart on failure"
  type        = bool
  default     = true
}

variable "on_host_maintenance" {
  description = "Maintenance behavior for the VM (MIGRATE or TERMINATE)"
  type        = string
  default     = "MIGRATE"
}

variable "deletion_protection" {
  description = "Enable deletion protection (must be disabled before removing the resource)"
  type        = bool
  default     = false
}
