# project config
project                       = "finternet-sandbox"
building_block                = "argocdhelm"
env                           = "dev"
region                        = "asia-southeast1"
gke_cluster_location          = "asia-southeast1"

# cluster sizing
gke_node_pool_instance_type   = "c2d-standard-4"
gke_node_pool_scaling_config = {
  desired_size = 2
  max_size = 3
  min_size = 0
}

# cluster networking
create_network                = true
network                       = "finternet-test"
subnetwork                    = "finternet-argocd-subnetwork"
cluster_secondary_range_name  = "finternet-argocd-cluster-secondary"
services_secondary_range_name = "finternet-argocd-services-secondary"

# cluster node pool configuration
gke_node_default_disk_size_gb = 30
enable_autoscaling            = false

# ArgoCD configuration
enable_argocd                 = false
argocd_namespace              = "argocd"
argocd_chart_version          = ""

# Applications configuration
enable_applications           = true
applications_repo_url         = "git@github.com:finternet-io/units-automation.git"
applications_target_revision  = "arocd-helmchart"
applications_path             = "manifests/application"