#############################
# VPC Network Outputs
#############################

output "vpc_network_name" {
  description = "Name of the VPC network created"
  value       = "${var.building_block}-${var.env}-vpc2-network"
}

output "vpc_network_self_link" {
  description = "Self link of the VPC network"
  value       = module.networking.network
}

output "public_subnet_name" {
  description = "Name of the public subnet"
  value       = module.networking.public_subnetwork_name
}

output "public_subnet_self_link" {
  description = "Self link of the public subnet"
  value       = module.networking.public_subnetwork
}

output "public_subnet_cidr" {
  description = "CIDR range of the public subnet"
  value       = module.networking.public_subnetwork_cidr_block
}

output "public_subnet_gateway" {
  description = "Gateway address of the public subnet"
  value       = module.networking.public_subnetwork_gateway
}

output "private_subnet_name" {
  description = "Name of the private subnet"
  value       = module.networking.private_subnetwork_name
}

output "private_subnet_self_link" {
  description = "Self link of the private subnet"
  value       = module.networking.private_subnetwork
}

output "private_subnet_cidr" {
  description = "CIDR range of the private subnet"
  value       = module.networking.private_subnetwork_cidr_block
}

output "private_subnet_gateway" {
  description = "Gateway address of the private subnet"
  value       = module.networking.private_subnetwork_gateway
}

output "cloud_nat_name" {
  description = "Name of the Cloud NAT created"
  value       = "${var.building_block}-${var.env}-vpc2-nat"
}

output "cloud_router_name" {
  description = "Name of the Cloud Router created"
  value       = "${var.building_block}-${var.env}-vpc2-router"
}

#############################
# Firewall Rules Outputs
#############################

output "firewall_rules_created" {
  description = "List of firewall rules created by the network module"
  value = [
    "allow-internal-public",
    "allow-internal-private",
    "allow-restricted-api-egress",
    "allow-http-https"
  ]
}

output "network_tag_public" {
  description = "Network tag used for public firewall rules"
  value       = "public"
}

#############################
# VM Instance Outputs
#############################

output "vm_instance_name" {
  description = "Name of the VM instance created"
  value       = module.gce_vm.instance_names
}

output "vm_instance_zone" {
  description = "Zone where the VM instance is deployed"
  value       = module.gce_vm.instance_zones
}

output "vm_internal_ip" {
  description = "Internal IP address of the VM instance"
  value       = module.gce_vm.instance_internal_ip
}

output "vm_external_ip" {
  description = "External (public) IP address of the VM instance"
  value       = module.gce_vm.instance_external_ip
}

output "vm_self_link" {
  description = "Self link of the VM instance"
  value       = module.gce_vm.instances_self_links
}

output "vm_network_tag" {
  description = "Network tag applied to the VM (for firewall rules)"
  value       = var.vm_network_tag
}

#############################
# Summary Output
#############################

output "deployment_summary" {
  description = "Summary of all resources created"
  value = {
    project        = var.project
    region         = var.region
    zone           = var.zone
    # vpc_network    = module.networking.network_name
    subnets        = {
      public  = module.networking.public_subnetwork_name
      private = module.networking.private_subnetwork_name
    }
    vm_instance    = {
      name        = module.gce_vm.instance_names
      internal_ip = module.gce_vm.instance_internal_ip
      external_ip = module.gce_vm.instance_external_ip
      zone        = module.gce_vm.instance_zones
    }
    firewall_tag   = var.vm_network_tag
  }
}

#############################
# Connection Information
#############################

output "ssh_command" {
  description = "SSH command to connect to the VM instance"
  value       = "gcloud compute ssh ${var.vm_name} --project=${var.project} --zone=${var.zone}"
}

output "vm_console_url" {
  description = "GCP Console URL for the VM instance"
  value       = "https://console.cloud.google.com/compute/instancesDetail/zones/${var.zone}/instances/${var.vm_name}?project=${var.project}"
}
