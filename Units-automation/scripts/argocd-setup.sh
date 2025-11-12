#!/bin/bash
# Update and install dependencies
sudo apt update -y
sudo apt install -y curl

# Install Docker
curl -fsSL https://get.docker.com | sh

# Install k3s (single-node Kubernetes)
curl -sfL https://get.k3s.io | sh -

# Set kubectl context
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

# Install ArgoCD (latest)
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
