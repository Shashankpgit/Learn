variable "argocd_namespace" {
  description = "The namespace to deploy ArgoCD into"
  type        = string
  default     = "argocd"
}

variable "argocd_release_name" {
  description = "The Helm release name for ArgoCD"
  type        = string
  default     = "argocd"
}

variable "argocd_chart_version" {
  description = "The version of the ArgoCD Helm chart to deploy. Leave empty for latest version."
  type        = string
  default     = "" # Latest version
}

variable "argocd_service_type" {
  description = "The Kubernetes service type for ArgoCD server (LoadBalancer, NodePort, or ClusterIP)"
  type        = string
  default     = "LoadBalancer"
  validation {
    condition     = contains(["LoadBalancer", "NodePort", "ClusterIP"], var.argocd_service_type)
    error_message = "Service type must be one of: LoadBalancer, NodePort, ClusterIP"
  }
}

variable "argocd_server_insecure" {
  description = "Whether to run ArgoCD server in insecure mode (no TLS)"
  type        = bool
  default     = false
}

variable "argocd_custom_values" {
  description = "Map of custom values to pass to the ArgoCD Helm chart"
  type        = map(string)
  default     = {}
}

variable "helm_timeout" {
  description = "Timeout for Helm operations in seconds"
  type        = number
  default     = 600
}

variable "helm_wait" {
  description = "Whether to wait for all resources to be ready before marking the release as successful"
  type        = bool
  default     = true
}