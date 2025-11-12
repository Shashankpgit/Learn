variable "building_block" {
  type        = string
  description = "Building block name. All resources will be prefixed with this value."
  default     = "finternet"
  validation {
    condition     = length(var.building_block) > 0
    error_message = "The building block name must not be empty."
  }
  validation {
    condition     = can(regex("[a-zA-Z0-9\\-]+", var.building_block))
    error_message = "The building block name can only contain alphanumeric characters and dashes."
  }
}

variable "env" {
  type        = string
  description = "Environment name. All resources will be prefixed with this value."
  default     = "dev"
}

## Google Cloud Platform

variable "project" {
  description = "The project ID where all resources will be launched."
  type        = string
  default     = "finternet-gcp"
}

variable "region" {
  description = "The region for the network. If the cluster is regional, this must be the same region. Otherwise, it should be the region of the zone."
  type        = string
  default     = "asia-south1"
}

variable "zone" {
  description = "The zone for the cluster. If the cluster is regional, this should be one of the zones in the region. Otherwise, this should be the same zone as the region."
  type        = string
  default     = "asia-south1-a"
}

variable "gke_node_default_disk_size_gb" {
  description = "Default disk size for GKE nodes"
  type        = number
  default     = 30
}

## VPC and Subnetworks
variable "create_network" {
  description = "Create a new VPC network."
  type        = bool
  default     = true
}

variable "network" {
  description = "The VPC network to use. If create_network is true, this will be the name of the new network."
  type        = string
  default     = ""
}

variable "subnetwork" {
  description = "The subnetwork to use. If create_network is true, this will be the name of the new subnetwork."
  type        = string
  default     = ""
}

variable "cluster_secondary_range_name" {
  description = "The name associated with the pod subnetwork secondary range, used when adding an alias IP range to a VM instance. The name must be 1-63 characters long, and comply with RFC1035. The name must be unique within the subnetwork."
  type        = string
  default     = ""
}

variable "services_secondary_range_name" {
  description = "The name associated with the services subnetwork secondary range, used when adding an alias IP range to a VM instance. The name must be 1-63 characters long, and comply with RFC1035. The name must be unique within the subnetwork."
  type        = string
  default     = ""
}


variable "cidr_block" {
  type        = string
  description = "VPC CIDR range"
  default     = "10.0.0.0/16"
}

variable "secondary_cidr_block" {
  description = "The IP address range of the VPC's secondary address range in CIDR notation. A prefix of /16 is recommended. Do not use a prefix higher than /27."
  type        = string
  default     = "10.1.0.0/20"
}

variable "auto_assign_public_ip" {
  type        = bool
  description = "Auto assign public ip's to instances in this subnet"
  default     = true
}

variable "public_subnetwork_secondary_range_name" {
  description = "The name associated with the pod subnetwork secondary range, used when adding an alias IP range to a VM instance. The name must be 1-63 characters long, and comply with RFC1035. The name must be unique within the subnetwork."
  type        = string
  default     = "public-cluster"
}

variable "public_services_secondary_range_name" {
  description = "The name associated with the services subnetwork secondary range, used when adding an alias IP range to a VM instance. The name must be 1-63 characters long, and comply with RFC1035. The name must be unique within the subnetwork."
  type        = string
  default     = "public-services"
}

variable "public_services_secondary_cidr_block" {
  description = "The IP address range of the VPC's public services secondary address range in CIDR notation. A prefix of /16 is recommended. Do not use a prefix higher than /27. Note: this variable is optional and is used primarily for backwards compatibility, if not specified a range will be calculated using var.secondary_cidr_block, var.secondary_cidr_subnetwork_width_delta and var.secondary_cidr_subnetwork_spacing."
  type        = string
  default     = null
}

variable "private_services_secondary_cidr_block" {
  description = "The IP address range of the VPC's private services secondary address range in CIDR notation. A prefix of /16 is recommended. Do not use a prefix higher than /27. Note: this variable is optional and is used primarily for backwards compatibility, if not specified a range will be calculated using var.secondary_cidr_block, var.secondary_cidr_subnetwork_width_delta and var.secondary_cidr_subnetwork_spacing."
  type        = string
  default     = null
}

variable "secondary_cidr_subnetwork_width_delta" {
  description = "The difference between your network and subnetwork's secondary range netmask; an /16 network and a /20 subnetwork would be 4."
  type        = number
  default     = 4
}

variable "secondary_cidr_subnetwork_spacing" {
  description = "How many subnetwork-mask sized spaces to leave between each subnetwork type's secondary ranges."
  type        = number
  default     = 0
}

variable "igw_cidr" {
  type        = list(string)
  description = "Internet gateway CIDR range."
  default     = ["0.0.0.0/0"]
}

variable "gcs_service_account_name" {
  description = "The name of the custom service account used for GCS. This parameter is limited to a maximum of 28 characters."
  type        = string
  default     = "gcs-object-admin"
}

variable "gcs_service_account_description" {
  description = "A description of the custom service account used for the GKE cluster."
  type        = string
  default     = "GCS Service Account managed by Terraform"
}

#GCE VM Variables

variable "vm_name" {
  description = "Name of the VM instance. This is used to generate the actual name of the VM."
  type        = string
  default     = "finternet-vm"
}

variable "vm_zone" {
  description = "The zone where the VM instance will be created."
  type        = string
  default     = "asia-south1-a"
  
}

variable "vm_machine_type" {
  description = "Machine type of the VM instance"
  type        = string
  default     = "e2-medium"
}

variable "vm_image" {
  description = "The image to use for the VM instance"
  type        = string
  default     = "projects/debian-cloud/global/images/family/debian-12"
}

variable "vm_disk_size" {
  description = "Disk size in GB"
  type        = number
  default     = 10
}

variable "vm_disk_type" {
  description = "Disk type (pd-standard, pd-ssd, pd-balanced)"
  type        = string
  default     = "pd-standard"
}

variable "vm_auto_delete_boot_disk" {
  description = "Whether to auto-delete the boot disk when the VM is deleted."
  type        = bool
  default     = true
  
}

variable "ssh_keys" {
  description = "SSH keys to be added to the VM instance"
  type        = string
  default     = ""
}

variable "preemptible" {
  description = "Whether the VM instance is preemptible"
  type        = bool
  default     = false
  
}

variable "automatic_restart" {
  description = "Whether the VM instance should automatically restart if it is terminated by Compute Engine."
  type        = bool
  default     = true
}

variable "on_host_maintenance" {
  description = "Specifies the maintenance behavior for the VM instance. Valid values are MIGRATE and TERMINATE."
  type        = string
  default     = "MIGRATE"
}

variable "create_service_account" {
  description = "Whether to create a new service account for the VM instance"
  type        = bool
  default     = false
}
variable "vm_enable_public_ip" {
  description = "Whether to assign a public IP to the VM instance"
  type        = bool
  default     = true  
}

variable "vm_network_tag" {
  description = "Network tag to be applied to the VM instance for firewall rules"
  type        = string
  default     = "public"
}

variable "vm_prevent_destroy" {
  description = "Whether to prevent the VM from being destroyed"
  type        = bool
  default     = false
}

variable "vm_service_account_description" {
  description = "Description for the VM service account"
  type        = string
  default     = "VM Service Account managed by Terraform"
  
}

variable "vm_service_account_name" {
  description = "Name for the VM service account"
  type        = string
  default     = "vm-service-account"
  
}
