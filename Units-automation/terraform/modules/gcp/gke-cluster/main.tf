locals {
  workload_identity_config = !var.enable_workload_identity ? [] : var.identity_namespace == null ? [{
    identity_namespace = "${var.project}.svc.id.goog" }] : [{ identity_namespace = var.identity_namespace
  }]
}

resource "google_container_cluster" "cluster" {
  provider    = google

  name        = var.name
  description = var.description

  project    = var.project
  location   = var.location
  network    = var.network
  subnetwork = var.subnetwork

  logging_service    = "none"
  monitoring_service = "none"
  min_master_version = local.kubernetes_version
  deletion_protection = var.deletion_protection

  enable_legacy_abac = var.enable_legacy_abac


  remove_default_node_pool = true


  node_pool {
    name = "default-pool"
    initial_node_count = 2

    node_config {
      image_type   = "COS_CONTAINERD"
      machine_type = var.gke_node_pool_instance_type

      tags = var.gke_node_pool_network_tags

      disk_size_gb = var.gke_node_default_disk_size_gb
      disk_type    = var.kubernetes_storage_class
      preemptible  = var.gke_node_pool_preemptible

      service_account = var.alternative_default_service_account

      oauth_scopes = [
        "https://www.googleapis.com/auth/cloud-platform",
      ]

      confidential_nodes {
        enabled = true
      }

      shielded_instance_config {
        enable_secure_boot          = true
        enable_integrity_monitoring = true
      }
    }
  }

  ip_allocation_policy {
    cluster_secondary_range_name  = var.cluster_secondary_range_name
    services_secondary_range_name = var.services_secondary_range_name
  }

  private_cluster_config {
    enable_private_endpoint = var.disable_public_endpoint
    enable_private_nodes    = var.enable_private_nodes
    master_ipv4_cidr_block  = var.master_ipv4_cidr_block
  }

  addons_config {
    http_load_balancing {
      disabled = !var.http_load_balancing
    }

    horizontal_pod_autoscaling {
      disabled = !var.horizontal_pod_autoscaling
    }

    network_policy_config {
      disabled = !var.enable_network_policy
    }

    gcp_filestore_csi_driver_config {
      enabled = var.gcp_filestore_csi_driver
    }

    gce_persistent_disk_csi_driver_config {
      enabled = var.gce_persistent_disk_csi_driver
    }
  }

  network_policy {
    enabled = var.enable_network_policy

    provider = var.enable_network_policy ? "CALICO" : "PROVIDER_UNSPECIFIED"
  }

  vertical_pod_autoscaling {
    enabled = var.enable_vertical_pod_autoscaling
  }



  dynamic "master_authorized_networks_config" {
    for_each = var.master_authorized_networks_config
    content {
      dynamic "cidr_blocks" {
        for_each = lookup(master_authorized_networks_config.value, "cidr_blocks", [])
        content {
          cidr_block   = cidr_blocks.value.cidr_block
          display_name = lookup(cidr_blocks.value, "display_name", null)
        }
      }
    }
  }

  maintenance_policy {
    daily_maintenance_window {
      start_time = var.maintenance_start_time
    }
  }

  lifecycle {
    ignore_changes = [

      node_config, node_pool
    ]
  }

  dynamic "authenticator_groups_config" {
    for_each = [
      for x in [var.gsuite_domain_name] : x if var.gsuite_domain_name != null
    ]

    content {
      security_group = "gke-security-groups@${authenticator_groups_config.value}"
    }
  }

  dynamic "database_encryption" {
    for_each = [
      for x in [var.secrets_encryption_kms_key] : x if var.secrets_encryption_kms_key != null
    ]

    content {
      state    = "ENCRYPTED"
      key_name = database_encryption.value
    }
  }

  dynamic "workload_identity_config" {
    for_each = local.workload_identity_config

    content {
      workload_pool = workload_identity_config.value.identity_namespace
    }
  }

  resource_labels = var.resource_labels
}

resource "google_container_node_pool" "node_pool" {
  provider = google

  name     = "${var.building_block}-${var.env}-pool"
  project  = var.project
  location = var.zone
  cluster  = google_container_cluster.cluster.name

  initial_node_count = var.enable_autoscaling ? var.gke_node_pool_scaling_config["desired_size"] : null
  node_count = var.gke_node_pool_scaling_config["desired_size"]

dynamic "autoscaling" {
  for_each = var.enable_autoscaling ? [1] : []
  content {
    min_node_count = var.gke_node_pool_scaling_config["min_size"]
    max_node_count = var.gke_node_pool_scaling_config["max_size"]
  }
}

  management {
    auto_repair  = "true"
    auto_upgrade = "true"
  }

  node_config {
    image_type   = "COS_CONTAINERD"
    machine_type = var.gke_node_pool_instance_type

    tags = var.gke_node_pool_network_tags

    disk_size_gb = var.gke_node_default_disk_size_gb
    disk_type    = var.kubernetes_storage_class
    preemptible  = var.gke_node_pool_preemptible

    service_account = var.alternative_default_service_account

    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]

    confidential_nodes {
      enabled = true
    }

    shielded_instance_config {
      enable_secure_boot          = true
      enable_integrity_monitoring = true
    }
  }

  lifecycle {
    ignore_changes = [initial_node_count]
  }

  timeouts {
    create = "60m"
    update = "60m"
    delete = "60m"
  }
}


locals {
  latest_version     = data.google_container_engine_versions.location.latest_master_version
  kubernetes_version = var.kubernetes_version != "latest" ? var.kubernetes_version : local.latest_version
  network_project    = var.network_project != "" ? var.network_project : var.project
}

data "google_container_engine_versions" "location" {
  location = var.location
  project  = var.project
}

