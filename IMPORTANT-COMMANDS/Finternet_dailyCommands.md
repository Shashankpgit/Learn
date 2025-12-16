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
8. gcloud container clusters get-credentials finternet-pw-lab-dev-cluster \
  --region asia-southeast1-a \
  --project finternet-sandbox  ` to add the kubeconfig`
9. gsutil rm gs://fintenet-foundry-argocd/terraform/state/default.tflock
10. `gcloud iam workload-identity-pools list  --project=finternet-sandbox  --location=global`
11. `gcloud iam workload-identity-pools providers list --project=finternet-sandbox --location=global --workload-identity-pool=finternet-sandbox-github-pool`
12. `./install.sh install --provider gcp --config ./env.conf --install_dependencies install --resource gke`


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





anti-Gravity
