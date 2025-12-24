gcloud auth application-default login --account=shashank.paramesh@finternetlab.io

GCP:
# issue here i faced was i was trying to initialize the create the new cluster with the old terraform state.
# This command removed everything related to state files.
1. `rm -rf .terraform terraform.tfstate terraform.tfstate.backup`
2. `gcloud auth application-default login --account=shashank.paramesh@finternetlab.io`
3. `gsutil ls` and eneter the password
4. `gcloud container clusters list --project finternet-sandbox` use this command to know about the cluster in the project
5. `kubectl patch storageclass <storage-class-name> \
  -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}' to set the default storage class
6. `gcloud container node-pools list --cluster=finternet-dev-cluster --zone=asia-southeast1-a` list nodes
7. `gcloud container node-pools describe finternet-dev-pool --cluster=finternet-dev-cluster --zone=asia-southeast1-a`
8. gcloud container clusters get-credentials finternet-pwd-dev-cluster \
  --region asia-southeast1-a \
  --project finternet-sandbox  ` to add the kubeconfig`LTNv~XPQQViZ0~
9. gsutil rm gs://fintenet-foundry-argocd/terraform/state/default.tflock
10. `gcloud iam workload-identity-pools list  --project=finternet-sandbox  --location=global`
11. `gcloud iam workload-identity-pools providers list --project=finternet-sandbox --location=global --workload-identity-pool=finternet-sandbox-github-pool`
12. `./install.sh install --provider gcp --config ./env.conf --install_dependencies true --resource gke`


kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

kubectl create secret generic repo-units-automation \
  --from-file=sshPrivateKey=/home/sankethika/.ssh/id_ed25519 \
  --from-literal=url=git@github.com:finternet-io/units-automation.git \
  -n argocd \
  --type=Opaque

kubectl label secret repo-units-automation \
  -n argocd argocd.argoproj.io/secret-type=repository

kubectl create secret generic repo-helmcharts \
  --from-file=sshPrivateKey=/home/sankethika/.ssh/id_ed25519 \
  --from-literal=url=git@github.com:finternet-io/helmcharts.git \
  -n argocd \
  --type=Opaque

kubectl label secret repo-helmcharts \
  -n argocd argocd.argoproj.io/secret-type=repository

for ns in $(kubectl get ns --no-headers | awk '$2=="Terminating" {print $1}'); do
  echo "🧨 Force deleting terminating namespace: $ns"
  kubectl get namespace "$ns" -o json | \
    jq 'del(.spec.finalizers)' | \
    kubectl replace --raw "/api/v1/namespaces/$ns/finalize" -f -
done

export GOOGLE_PROJECT_ID="finternet-sandbox"
export GOOGLE_TERRAFORM_BACKEND_LOCATION="asia-southeast1"
export GOOGLE_TERRAFORM_BACKEND_BUCKET="fintenet-automation-11"

echo $GOOGLE_PROJECT_ID
echo $GOOGLE_TERRAFORM_BACKEND_LOCATION
echo $GOOGLE_TERRAFORM_BACKEND_BUCKET-S → tells Git to sign the commit with your GPG key-S → tells Git to sign the commit with your GPG key.
---------------------------------------------------------------------------------------------------------------------------------------------


Steps while creating the cluster:
1. Update the tfvars file
2. Update the env.conf file
3. 

---------------------------------------------------------

`googleapi: Error 403: Insufficient regional quota to satisfy request: resource "SSD_TOTAL_GB": request requires '60.0' and is short '14.0'. project has a quota of '500.0' with '46.0' available. View and manage quotas at https://console.cloud.google.com/iam-admin/quotas?usage=USED&project=finternet-sandbox.`



cluster creation started at 11:16:50
support@finternetlab.io
Finternet@1234

anti-Gravity
Use wisperflow to modif the and give cursor the proper context


  lifecycle {
    ignore_changes = [
      initial_node_count,
      node_config.oauth_scopes
    ]
  }
  
  
  URL: https://cloud.google.com/kubernetes-engine/docs/how-to/access-scopes
  
  
  
  
  
  `Terraform used the selected providers to generate the following execution
plan. Resource actions are indicated with the following symbols:
  + create
  ~ update in-place
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # module.argocd[0].helm_release.argocd will be updated in-place
  ~ resource "helm_release" "argocd" {
        id                         = "argocd"
      ~ metadata                   = [
          - {
              - app_version = "v3.2.2"
              - chart       = "argo-cd"
              - name        = "argocd"
              - namespace   = "argocd"
              - revision    = 1
              - values      = jsonencode(
                    {
                      - configs = {
                          - secret = {
                              - argocdServerAdminPassword = "$2a$10$SRmeoCf4Yz2f0K9VZPvEMeMd7/i6PT4XR1KUdadngCrNR0AGvnus."
                            }
                        }
                    }
                )
              - version     = "9.2.0"
            },
        ] -> (known after apply)
        name                       = "argocd"
      ~ values                     = [
          - (sensitive value),
        ] -> (known after apply)
        # (26 unchanged attributes hidden)
    }

  # module.gke_cluster.google_container_node_pool.node_pool must be replaced
-/+ resource "google_container_node_pool" "node_pool" {
      ~ id                          = "projects/finternet-sandbox/locations/asia-southeast1-a/clusters/finternet-release31-dev-cluster/nodePools/finternet-release31-dev-pool" -> (known after apply)
      ~ initial_node_count          = 2 -> (known after apply)
      ~ instance_group_urls         = [
          - "https://www.googleapis.com/compute/v1/projects/finternet-sandbox/zones/asia-southeast1-a/instanceGroupManagers/gke-finternet-releas-finternet-releas-03f02da9-grp",
        ] -> (known after apply)
      ~ managed_instance_group_urls = [
          - "https://www.googleapis.com/compute/v1/projects/finternet-sandbox/zones/asia-southeast1-a/instanceGroups/gke-finternet-releas-finternet-releas-03f02da9-grp",
        ] -> (known after apply)
      ~ max_pods_per_node           = 110 -> (known after apply)
        name                        = "finternet-release31-dev-pool"
      + name_prefix                 = (known after apply)
      ~ node_locations              = [
          - "asia-southeast1-a",
        ] -> (known after apply)
      + operation                   = (known after apply)
      ~ version                     = "1.34.1-gke.3759000" -> (known after apply)
        # (4 unchanged attributes hidden)

      ~ network_config (known after apply)
      - network_config {
          - create_pod_range     = false -> null
          - enable_private_nodes = true -> null
          - pod_ipv4_cidr_block  = "10.1.0.0/20" -> null
          - pod_range            = "public-cluster" -> null
          - subnetwork           = "projects/finternet-sandbox/regions/asia-southeast1/subnetworks/finternet-release31-dev-subnetwork-public" -> null
        }

      ~ node_config {
          ~ effective_taints            = [] -> (known after apply)
          - enable_confidential_storage = false -> null
          - flex_start                  = false -> null
          ~ labels                      = {} -> (known after apply)
          ~ local_ssd_count             = 0 -> (known after apply)
          ~ logging_variant             = "DEFAULT" -> (known after apply)
          ~ metadata                    = {
              - "disable-legacy-endpoints" = "true"
            } -> (known after apply)
          + min_cpu_platform            = (known after apply)
          ~ oauth_scopes                = [ # forces replacement
              - "https://www.googleapis.com/auth/monitoring.write",
                # (3 unchanged elements hidden)
            ]
          - resource_labels             = {
              - "goog-gke-node-pool-provisioning-model" = "on-demand"
            } -> null
          - resource_manager_tags       = {} -> null
          - storage_pools               = [] -> null
            tags                        = [
                "public",
            ]
            # (11 unchanged attributes hidden)

          ~ boot_disk (known after apply)
          - boot_disk {
              - disk_type              = "pd-ssd" -> null
              - provisioned_iops       = 0 -> null
              - provisioned_throughput = 0 -> null
              - size_gb                = 30 -> null
            }

          ~ confidential_nodes {
              - confidential_instance_type = "SEV" -> null
                # (1 unchanged attribute hidden)
            }

          ~ gcfs_config (known after apply)

          ~ guest_accelerator (known after apply)

          ~ kubelet_config (known after apply)
          - kubelet_config {
              - allowed_unsafe_sysctls                 = [] -> null
              - container_log_max_files                = 0 -> null
              - cpu_cfs_quota                          = false -> null
              - eviction_max_pod_grace_period_seconds  = 0 -> null
              - image_gc_high_threshold_percent        = 0 -> null
              - image_gc_low_threshold_percent         = 0 -> null
              - insecure_kubelet_readonly_port_enabled = "FALSE" -> null
              - max_parallel_image_pulls               = 3 -> null
              - pod_pids_limit                         = 0 -> null
              - single_process_oom_kill                = false -> null
                # (5 unchanged attributes hidden)
            }

          ~ linux_node_config (known after apply)

          ~ windows_node_config (known after apply)
          - windows_node_config {
                # (1 unchanged attribute hidden)
            }

          ~ workload_metadata_config (known after apply)
          - workload_metadata_config {
              - mode = "GKE_METADATA" -> null
            }

            # (1 unchanged block hidden)
        }

      ~ upgrade_settings (known after apply)
      - upgrade_settings {
          - max_surge       = 1 -> null
          - max_unavailable = 0 -> null
          - strategy        = "SURGE" -> null
        }

        # (2 unchanged blocks hidden)
    }

  # module.password_generator.null_resource.global_overrides_yaml[0] must be replaced
-/+ resource "null_resource" "global_overrides_yaml" {
      ~ id       = "865659589390999324" -> (known after apply)
      ~ triggers = { # forces replacement
          # Warning: this attribute value will be marked as sensitive and will not
          # display in UI output after applying this change.
          ~ "auth_passwords" = (sensitive value)
          ~ "users_hash"     = "4f53cda18c2baa0c0354bb5f9a3ecbe5ed12ab4d8e11ba873c2f11161202b945" -> "9886686f3220422236b4c47dfcd4c96cac36da1242ef0618fe277f68af69093e"
            # (2 unchanged elements hidden)
        }
    }

  # module.password_generator.random_password.basic_auth_passwords["user1"] will be created
  + resource "random_password" "basic_auth_passwords" {
      + bcrypt_hash      = (sensitive value)
      + id               = (known after apply)
      + length           = 14
      + lower            = true
      + min_lower        = 0
      + min_numeric      = 0
      + min_special      = 0
      + min_upper        = 0
      + number           = true
      + numeric          = true
      + override_special = "!@%^*_-+={}[]~"
      + result           = (sensitive value)
      + special          = true
      + upper            = true
    }

  # module.password_generator.random_password.basic_auth_passwords["user2"] will be created
  + resource "random_password" "basic_auth_passwords" {
      + bcrypt_hash      = (sensitive value)
      + id               = (known after apply)
      + length           = 14
      + lower            = true
      + min_lower        = 0
      + min_numeric      = 0
      + min_special      = 0
      + min_upper        = 0
      + number           = true
      + numeric          = true
      + override_special = "!@%^*_-+={}[]~"
      + result           = (sensitive value)
      + special          = true
      + upper            = true
    }

Plan: 4 to add, 1 to change, 2 to destroy.`
here first i created the cluster usin this command
`./install.sh install --provider gcp --config ./env.conf --install_dependencies true --resource gke`
then i made a small changes in the   default     = ["user1", "user2"] in the password generator module then i again triggerd the same command
but here why am does the cluster getting destroying the exixting nodepools and creating the new one and also i didn't added this oauth permission
`oauth_scopes                = [ # forces replacement
              - "https://www.googleapis.com/auth/monitoring.write",
                # (3 unchanged elements hidden)
            ]` why it is showing force replacemets
