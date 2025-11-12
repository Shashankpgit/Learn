variable "applications_release_name" {
  description = "Helm release name for the applications chart"
  type        = string
  default     = "applications"
}

variable "applications_chart_path" {
  description = "Local path to the applications Helm chart"
  type        = string
  default     = "../../../manifests/application"
}

variable "argocd_namespace" {
  description = "Namespace where ArgoCD is deployed"
  type        = string
  default     = "argocd"
}

variable "applications_repo_url" {
  description = "Git repo URL for the applications"
  type        = string
  default     = "git@github.com:finternet-io/units-automation.git"
}

variable "applications_target_revision" {
  description = "Branch/tag for the applications repo"
  type        = string
  default     = "argocd-setup"
}

variable "applications_path" {
  description = "Path for the root-application"
  type        = string
  default     = "."
}

variable "helm_timeout" {
  description = "Timeout for Helm operations"
  type        = number
  default     = 600
}

variable "helm_wait" {
  description = "Whether to wait for resources"
  type        = bool
  default     = true
}

variable "depends_on_argocd" {
  description = "Dependency on ArgoCD"
  type        = any
  default     = []
}