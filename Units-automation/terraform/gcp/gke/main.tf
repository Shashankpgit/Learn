terraform {
  backend "gcs" { }

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 7.8.0"
    }

        # google-beta = {
    #   source  = "hashicorp/google-beta"
    #   version = "~> 7.8.0"
    # }
    # local = {
    #   source = "hashicorp/local"
    #   version  = "~> 2.5.1"
    # }
    # helm = {
    #   source = "hashicorp/helm"
    #   version  = "~> 2.13.2"
    # }

    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.30.0"
    }
    
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.13.2"
    }
  }
}

provider "google" {
  project = var.project
  region  = var.region

  scopes = [
    "https://www.googleapis.com/auth/compute",
    "https://www.googleapis.com/auth/cloud-platform",
    "https://www.googleapis.com/auth/ndev.clouddns.readwrite",
    "https://www.googleapis.com/auth/devstorage.full_control",

    "https://www.googleapis.com/auth/userinfo.email",
  ]
}

# provider "google-beta" {
#   project = var.project
#   region  = var.region

#   scopes = [
#     # Default scopes
#     "https://www.googleapis.com/auth/compute",
#     "https://www.googleapis.com/auth/cloud-platform",
#     "https://www.googleapis.com/auth/ndev.clouddns.readwrite",
#     "https://www.googleapis.com/auth/devstorage.full_control",

#     # Required for google_client_openid_userinfo
#     "https://www.googleapis.com/auth/userinfo.email",
#   ]
# }
provider "kubernetes" {
  host                   = "https://${module.gke_cluster.endpoint}"
  token                  = data.google_client_config.default.access_token
  cluster_ca_certificate = module.gke_cluster.cluster_ca_certificate
}
provider "helm" {
  kubernetes {
    host                   = "https://${module.gke_cluster.endpoint}"
    token                  = data.google_client_config.default.access_token
    cluster_ca_certificate = module.gke_cluster.cluster_ca_certificate
  }
}

data "google_client_config" "default" {}

module "networking" {
  count                 = var.create_network ? 1 : 0
  source                = "../../modules/gcp/vpc-network"

  name_prefix           = "${var.building_block}-${var.env}"
  project               = var.project
  region                = var.region

  cidr_block            = var.vpc_cidr_block
  secondary_cidr_block  = var.vpc_secondary_cidr_block

  public_subnetwork_secondary_range_name = var.public_subnetwork_secondary_range_name
  public_services_secondary_range_name   = var.public_services_secondary_range_name
  public_services_secondary_cidr_block   = var.public_services_secondary_cidr_block
  private_services_secondary_cidr_block  = var.private_services_secondary_cidr_block
  secondary_cidr_subnetwork_width_delta  = var.secondary_cidr_subnetwork_width_delta
  secondary_cidr_subnetwork_spacing      = var.secondary_cidr_subnetwork_spacing

  igw_cidr              = var.igw_cidr
}

module "cloud_storage" {
  source          = "../../modules/gcp/cloud-storage"
  building_block  = var.building_block
  env             = var.env
  project         = var.project
  region          = var.region
}

module "gke_cluster" {
  source = "../../modules/gcp/gke-cluster"

  building_block                = var.building_block
  env                           = var.env

  name                          = "${var.building_block}-${var.env}-cluster"
  project                       = var.project
  location                      = var.zone 
  zone                          = var.zone 
  network                       = var.create_network ? module.networking[0].network : var.network

  subnetwork                    = var.create_network ? module.networking[0].public_subnetwork : var.subnetwork
  cluster_secondary_range_name  = var.create_network ? module.networking[0].public_subnetwork_secondary_range_name : var.cluster_secondary_range_name
  services_secondary_range_name = var.create_network ? module.networking[0].public_services_secondary_range_name : var.services_secondary_range_name

  master_ipv4_cidr_block        = var.gke_master_ipv4_cidr_block

  enable_private_nodes          = true

  disable_public_endpoint       = false

  master_authorized_networks_config = [
    {
      cidr_blocks = [
        {
          cidr_block   = var.igw_cidr[0]
          display_name = "IGW"
        },
      ]
    },
  ]

  gke_node_pool_network_tags      = var.create_network ? [module.networking[0].public] : []
  gke_node_default_disk_size_gb   = var.gke_node_default_disk_size_gb

  gke_node_pool_instance_type     = var.gke_node_pool_instance_type
  gke_node_pool_scaling_config    = var.gke_node_pool_scaling_config

  enable_vertical_pod_autoscaling = var.enable_vertical_pod_autoscaling
  enable_autoscaling              = var.enable_autoscaling

  resource_labels = {
    environment = var.env
  }
}

resource "null_resource" "configure_kubectl" {
  provisioner "local-exec" {

    command = "gcloud container clusters get-credentials ${module.gke_cluster.name} --region ${var.zone} --project ${var.project}"

    environment = {
      KUBECONFIG = var.kubectl_config_path != "" ? var.kubectl_config_path : "credentials/config-${var.building_block}-${var.env}.yaml"
    }
  }

  depends_on = [ module.gke_cluster ]
}

module "argocd" {
  count  = var.enable_argocd ? 1 : 0
  source = "../../modules/gcp/argocd"

  argocd_namespace      = var.argocd_namespace
  argocd_chart_version  = var.argocd_chart_version

  depends_on = [
    module.gke_cluster,
    null_resource.configure_kubectl
  ]
}

module "applications" {
  count  = var.enable_argocd && var.enable_applications ? 1 : 0
  source = "../../modules/gcp/applications"

  argocd_namespace             = var.argocd_namespace
  applications_repo_url       = var.applications_repo_url
  applications_target_revision = var.applications_target_revision
  applications_path           = var.applications_path

  depends_on_argocd = [module.argocd]

  depends_on = [
    module.argocd
  ]
}

resource "local_file" "service_credentials" {
  count = var.enable_argocd ? 1 : 0

  filename = "${path.module}/service-credentials.json"
  content = <<EOF
argocd_username: admin
argocd_password: ${module.argocd[0].admin_password}
EOF

  depends_on = [module.argocd]
}

resource "null_resource" "installation_complete" {
  provisioner "local-exec" {
    command = "echo 'Installation completed successfully'"
  }

  depends_on = [
    module.gke_cluster,
    null_resource.configure_kubectl,
    module.argocd,
    local_file.service_credentials
  ]
}

