terraform {
  backend "gcs" { }

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 7.8.0"
    }

    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 7.8.0"
    }
  }
}

provider "google" {
  project = var.project
  region  = var.region

  # scopes = [
  #   # Default scopes
  #   "https://www.googleapis.com/auth/compute",
  #   "https://www.googleapis.com/auth/cloud-platform",
  #   "https://www.googleapis.com/auth/ndev.clouddns.readwrite",
  #   "https://www.googleapis.com/auth/devstorage.full_control",

  #   # Required for google_client_openid_userinfo
  #   "https://www.googleapis.com/auth/userinfo.email",
  # ]
}

# provider "google-beta" {
#   project = var.project
#   region  = var.region

#   # scopes = [
#   #   # Default scopes
#   #   "https://www.googleapis.com/auth/compute",
#   #   "https://www.googleapis.com/auth/cloud-platform",
#   #   "https://www.googleapis.com/auth/ndev.clouddns.readwrite",
#   #   "https://www.googleapis.com/auth/devstorage.full_control",

#   #   # Required for google_client_openid_userinfo
#   #   "https://www.googleapis.com/auth/userinfo.email",
#   # ]
# }

# Second VPC and its resources
module "networking" {
  source                = "../../modules/gcp/vpc-network"
  name_prefix           = "${var.building_block}-${var.env}-vpc2"
  project               = var.project
  region                = var.region

  cidr_block            = var.cidr_block
  secondary_cidr_block  = var.secondary_cidr_block

  public_subnetwork_secondary_range_name = var.public_subnetwork_secondary_range_name
  public_services_secondary_range_name   = var.public_services_secondary_range_name
  public_services_secondary_cidr_block   = var.public_services_secondary_cidr_block
  secondary_cidr_subnetwork_width_delta  = var.secondary_cidr_subnetwork_width_delta
  secondary_cidr_subnetwork_spacing      = var.secondary_cidr_subnetwork_spacing

  igw_cidr              = var.igw_cidr
}

# GCE VM in new VPC
module "gce_vm" {
  source                        = "../../modules/gcp/gce-vm"
  project_id                    = var.project
  vm_name                       = var.vm_name
  vm_zone                       = var.zone
  vm_machine_type               = var.vm_machine_type
  vm_image                      = var.vm_image
  vm_disk_size                  = var.vm_disk_size
  vm_disk_type                  = var.vm_disk_type
  auto_delete_boot_disk         = var.vm_auto_delete_boot_disk

  vm_network_tag                = var.vm_network_tag
  vm_network_self_link          = module.networking.network
  vm_subnetwork_self_link       = module.networking.public_subnetwork

  enable_public_ip              = var.vm_enable_public_ip
  ssh_keys                      = var.ssh_keys
  preemptible                   = var.preemptible
  automatic_restart             = var.automatic_restart
  on_host_maintenance           = var.on_host_maintenance
}

