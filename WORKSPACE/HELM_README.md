# Finternet App Helm Chart

This Helm chart simplifies the deployment of the Finternet application (React frontend + Node.js backend) on Kubernetes. It manages the lifecycle of the application, including networking, environment configuration, and scaling.

## What this Helm Chart does:

1.  **Deployment Management**: Creates a Kubernetes `Deployment` that runs the Docker container for the Finternet application. It supports scaling via `replicaCount`.
2.  **Service Discovery**: Sets up a `Service` (default: `ClusterIP`) to allow internal or external traffic to reach the application pods.
3.  **Configurable Environment**: Uses a `ConfigMap` to manage environment variables such as `TARGET_URL`, `PORT`, and `HOST`. These are easily customizable through `values.yaml`.
4.  **Ingress Support**: Optionally creates an `Ingress` resource to expose the application to the internet via a domain name.
5.  **Health Monitoring**: Implements `livenessProbe` and `readinessProbe` using the application's `/health` endpoint to ensure high availability and automatic recovery.
6.  **Resource Control**: Defines CPU and memory limits/requests to ensure stable cluster performance.

## Prerequisites

-   A Kubernetes cluster.
-   Helm 3.x installed.
-   The Docker image for the application pushed to a registry (or available locally in the cluster).

## Installation

1.  **Navigate to the chart directory**:
    ```bash
    cd helmchart
    ```

2.  **Install the chart**:
    ```bash
    helm install my-finternet-app . -n <namespace>
    ```

## Configuration

You can override any value in `values.yaml` during installation:

```bash
helm install my-finternet-app . \
  --set env.TARGET_URL="https://api.prod-finternet.com" \
  --set replicaCount=3
```

### Key Parameters in `values.yaml`:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of pod replicas | `1` |
| `image.repository` | Docker image name | `finternet-app` |
| `image.tag` | Image tag | `latest` |
| `service.port` | Port exposed by the service | `80` |
| `service.targetPort` | Port the container listens on | `3000` |
| `env.TARGET_URL` | **(Required)** Target API URL for proxying | `https://api.example.com` |
| `ingress.enabled` | Whether to create an Ingress resource | `false` |
| `ingress.hosts` | List of hosts for Ingress | `finternet.local` |

## Health Checks

The deployment includes health checks that monitor the `/health` endpoint. If the application becomes unresponsive, Kubernetes will automatically restart the pod.

