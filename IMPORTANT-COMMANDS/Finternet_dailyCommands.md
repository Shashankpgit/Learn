# Terraform Notes
# Personal Application
Sqash
--------------------------------------------------------------------------------------
#linux commands:
1. awk: awk automatically splits each line of a file into fields (columns), usually separated by spaces or tabs.




sankethika@sankethika-Vostro-3405:~/obsrv-argocd/obsrv2/obsrv-automation/helmcharts/kitchen/argocd-apps$ k apply -f application.yaml 
Warning: metadata.finalizers: "resources-finalizer.argocd.argoproj.io": prefer a domain-qualified finalizer name to avoid accidental conflicts with other finalizer writers
application.argoproj.io/obsrv-platform created
sankethika@sankethika-Vostro-3405:~/obsrv-argocd/obsrv2/obsrv-automation/helmcharts/kitchen/argocd-apps$ 


-------------------------------------------------------------------
###--> Personal Application
* Flow Started with the terraform scripts to create the cluster and resources.
1. Created the project manually.
2. created the SA for that project to use it in the cluster creation.
3. order in which i need to create the resources:
    - 1. VPC
    - 2. Subnets
    - 3. Secondary IP ranges
    - 4. 
hii ravi i configured the argocd to obsrv automation,
so first i built the 

gcloud auth application-default login --account=shashank.paramesh@finternetlab.io



-----------------------------------------------------------------------------------------------------------------------------------------------------
## Terraform Notes
 
#--> Topics Learnt:
* State files
* remote backend
* How to store remote state files in the backend
    1. its advantages
* Modules:
    1. here we are we are using the modules just as the user defined functions, where we can call this modules from the root main.tf files from which we can easily give the required resoure details like e2-micro,10gd disk kind of.
    
#`Error: Error waiting for creating GKE cluster: 
│       - requested resource is exhausted: Not all instances running in IGM after 38.32514773s. Expected 1, running 0, transitioning 1. Current errors: [IP_SPACE_EXHAUSTED_WITH_DETAILS]: Instance 'gke-personal-app-gke-clu-default-pool-94017e78-3phn' creation failed: IP space of 'projects/personal-project-475209/regions/asia-south1/subnetworks/subnet-1' is exhausted. Insufficient free IP addresses in the IP range '10.0.2.0/24'. Consider expanding the current IP range or selecting an alternative IP range. If this is a secondary range, consider adding an additional secondary range
│       - requested resource is exhausted: Not all instances running in IGM after 40.301658933s. Expected 1, running 0, transitioning 1. Current errors: [IP_SPACE_EXHAUSTED_WITH_DETAILS]: Instance 'gke-personal-app-gke-clu-default-pool-61a2216f-m0kq' creation failed: IP space of 'projects/personal-project-475209/regions/asia-south1/subnetworks/subnet-1' is exhausted. Insufficient free IP addresses in the IP range '10.0.2.0/24'. Consider expanding the current IP range or selecting an alternative IP range. If this is a secondary range, consider adding an additional secondary range.
│ 
│   with module.gke.google_container_cluster.primary,
│   on modules/GKE/main.tf line 1, in resource "google_container_cluster" "primary":
│    1: resource "google_container_cluster" "primary" {`



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
8. gcloud container clusters get-credentials finternet-dev-cluster \
  --region asia-southeast1-a \
  --project finternet-sandbox  ` to add the kubeconfig`
9. gsutil rm gs://fintenet-foundry-argocd/terraform/state/default.tflock
10. `gcloud iam workload-identity-pools list  --project=finternet-sandbox  --location=global`
11. `gcloud iam workload-identity-pools providers list --project=finternet-sandbox --location=global --workload-identity-pool=finternet-sandbox-github-pool`
12. 


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


ClickStack:
1. `clickhouse-client`
2. General sql commands

export GOOGLE_PROJECT_ID="finternet-sandbox"
export GOOGLE_TERRAFORM_BACKEND_LOCATION="asia-southeast1"
export GOOGLE_TERRAFORM_BACKEND_BUCKET="fintenet-automation-11"

echo $GOOGLE_PROJECT_ID
echo $GOOGLE_TERRAFORM_BACKEND_LOCATION
echo $GOOGLE_TERRAFORM_BACKEND_BUCKET-S → tells Git to sign the commit with your GPG key-S → tells Git to sign the commit with your GPG key.

.



Error: Error deleting service account: googleapi: Error 403: Permission 'iam.serviceAccounts.delete' denied on resource (or it may not exist).
│ Details:
│ [
│   {
│     "@type": "type.googleapis.com/google.rpc.ErrorInfo",
│     "domain": "iam.googleapis.com",
│     "metadata": {
│       "permission": "iam.serviceAccounts.delete"
│     },
│     "reason": "IAM_PERMISSION_DENIED"
│   }
│ ]
│ , forbidden
│ 



PROJECT=argocd-demo-475109
gcloud config set project $PROJECT
REGION=asia-south1
ZONE=asia-south1-a

# create three clusters (standard, non-autopilot example)
gcloud container clusters create argocd-cluster \
  --zone $ZONE --num-nodes 2 --release-channel=regular

gcloud container clusters create cluster-1 \
  --zone $ZONE --num-nodes 1

gcloud container clusters create cluster-2 \
  --zone $ZONE --num-nodes 1
----------------------------------------



okay create me the mater dataset which is batch information(or batch metadata)
which contains
1. batch_number
2. batch_ID
3. approved_by(should be object)
       {(emp_name, emp_id, emp_email, emp_phone)
        }

where the primary key i will use as the bacthID 
4. storage_conditions
and lets create one normal dataset for the drug which will have 
1. drug name, drug
