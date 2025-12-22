# Certificate Management Flow in Finternet Infrastructure

## Table of Contents
1. [Overview](#overview)
2. [Components Involved](#components-involved)
3. [Certificate Creation Flow](#certificate-creation-flow)
4. [Certificate Usage](#certificate-usage)
5. [Detailed Component Analysis](#detailed-component-analysis)
6. [Configuration Details](#configuration-details)
7. [Troubleshooting](#troubleshooting)

---

## Overview

This document explains how TLS/SSL certificates are automatically created, managed, and used in the Finternet Kubernetes infrastructure. The system uses **cert-manager** with **Let's Encrypt** as the Certificate Authority (CA) to automatically provision and renew certificates for services exposed through the **Kong Ingress Controller**.

### Key Points
- **Automated Certificate Management**: Certificates are automatically requested, issued, and renewed
- **Let's Encrypt Integration**: Free, automated certificates from Let's Encrypt CA
- **Ingress-Driven**: Certificates are created when Ingress resources are deployed
- **Kong Integration**: Kong Ingress Controller handles TLS termination

---

## Components Involved

### 1. **cert-manager** (Namespace: cert-manager)
- **Purpose**: Kubernetes add-on to automate certificate management
- **Version**: v1.19.1
- **Location**: `charts/cert-manager/`
- **Role**: Core certificate lifecycle manager

### 2. **Let's Encrypt ClusterIssuer** (Namespace: cert-manager)
- **Purpose**: Defines how to obtain certificates from Let's Encrypt
- **Location**: `charts/letsencrypt-issuer/`
- **Role**: Certificate issuer configuration

### 3. **Kong Ingress Controller** (Namespace: ingress-controller)
- **Purpose**: API Gateway and Ingress Controller
- **Location**: `charts/ingress-controller/`
- **Role**: Routes traffic and terminates TLS

### 4. **Keycloak** (Namespace: keycloak)
- **Purpose**: Identity and Access Management
- **Location**: `charts/keycloak/`
- **Role**: Service that uses certificates via Ingress

---

## Certificate Creation Flow

### Phase 1: Infrastructure Setup

```
┌─────────────────────────────────────────────────────────────┐
│ Step 1: Deploy cert-manager                                 │
│ - Installs CRDs (Certificate, Issuer, ClusterIssuer, etc.) │
│ - Deploys cert-manager controller                          │
│ - Deploys cert-manager webhook                             │
│ - Deploys cert-manager cainjector                          │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Step 2: Deploy Let's Encrypt ClusterIssuer                  │
│ - Creates ClusterIssuer resource                            │
│ - Configures ACME protocol settings                         │
│ - Sets up HTTP-01 challenge solver                          │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Step 3: Deploy Kong Ingress Controller                      │
│ - Creates Kong Gateway pods                                 │
│ - Sets up IngressClass "kong"                               │
│ - Ready to handle Ingress resources                         │
└─────────────────────────────────────────────────────────────┘
```

### Phase 2: Certificate Request (Triggered by Ingress Creation)

```
┌─────────────────────────────────────────────────────────────┐
│ Step 1: Deploy Keycloak with Ingress                        │
│                                                              │
│ Ingress Resource Created:                                   │
│   apiVersion: networking.k8s.io/v1                          │
│   kind: Ingress                                              │
│   metadata:                                                  │
│     name: keycloak                                           │
│     namespace: keycloak                                      │
│     annotations:                                             │
│       cert-manager.io/cluster-issuer:                        │
│         "letsencrypt-issuer-letsencrypt-clusterissuer"      │
│       konghq.com/strip-path: "true"                         │
│       konghq.com/preserve-host: "true"                      │
│   spec:                                                      │
│     ingressClassName: kong                                   │
│     rules:                                                   │
│       - host: <your-domain>                                  │
│         http:                                                │
│           paths:                                             │
│             - path: /auth                                    │
│               pathType: Prefix                               │
│               backend:                                       │
│                 service:                                     │
│                   name: keycloak                             │
│                   port:                                      │
│                     number: 80                               │
│     tls:                                                     │
│       - hosts:                                               │
│           - <your-domain>                                    │
│         secretName: <your-domain>-tls                        │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Step 2: cert-manager Detects Ingress                        │
│                                                              │
│ cert-manager controller watches for Ingress resources with: │
│ - cert-manager.io/cluster-issuer annotation                 │
│ - tls section defined                                        │
│                                                              │
│ Action: Automatically creates Certificate resource          │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Step 3: Certificate Resource Created                        │
│                                                              │
│ apiVersion: cert-manager.io/v1                              │
│ kind: Certificate                                            │
│ metadata:                                                    │
│   name: <your-domain>-tls                                    │
│   namespace: keycloak                                        │
│ spec:                                                        │
│   secretName: <your-domain>-tls                              │
│   issuerRef:                                                 │
│     name: letsencrypt-issuer-letsencrypt-clusterissuer      │
│     kind: ClusterIssuer                                      │
│   dnsNames:                                                  │
│     - <your-domain>                                          │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Step 4: CertificateRequest Created                          │
│                                                              │
│ cert-manager creates CertificateRequest resource            │
│ This triggers the ACME protocol flow                        │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Step 5: ACME Order Created                                  │
│                                                              │
│ cert-manager contacts Let's Encrypt ACME server:            │
│ - URL: https://acme-v02.api.letsencrypt.org/directory       │
│ - Creates account (if first time)                           │
│ - Requests certificate for domain                           │
│                                                              │
│ Let's Encrypt responds with challenges                      │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Step 6: ACME Challenge Created (HTTP-01)                    │
│                                                              │
│ cert-manager creates:                                        │
│ 1. Challenge resource                                        │
│ 2. Temporary Pod (acme-solver)                              │
│ 3. Temporary Service                                         │
│ 4. Temporary Ingress for /.well-known/acme-challenge/       │
│                                                              │
│ Configuration from letsencrypt-issuer:                       │
│   solvers:                                                   │
│     - http01:                                                │
│         ingress:                                             │
│           ingressClassName: kong                             │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Step 7: Domain Validation                                   │
│                                                              │
│ Let's Encrypt validates domain ownership:                   │
│ 1. Makes HTTP request to:                                   │
│    http://<your-domain>/.well-known/acme-challenge/<token>  │
│ 2. Request goes through Kong Ingress Controller             │
│ 3. Kong routes to acme-solver pod                           │
│ 4. acme-solver responds with correct challenge response     │
│ 5. Let's Encrypt verifies response                          │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Step 8: Certificate Issuance                                │
│                                                              │
│ If validation successful:                                   │
│ 1. Let's Encrypt issues certificate                         │
│ 2. cert-manager receives certificate                        │
│ 3. cert-manager creates/updates Kubernetes Secret           │
│                                                              │
│ Secret created:                                              │
│   name: <your-domain>-tls                                    │
│   namespace: keycloak                                        │
│   type: kubernetes.io/tls                                    │
│   data:                                                      │
│     tls.crt: <base64-encoded-certificate>                   │
│     tls.key: <base64-encoded-private-key>                   │
│     ca.crt: <base64-encoded-ca-certificate>                 │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Step 9: Cleanup                                              │
│                                                              │
│ cert-manager removes temporary resources:                   │
│ - acme-solver pod                                            │
│ - Temporary service                                          │
│ - Temporary ingress                                          │
│ - Challenge resource                                         │
└─────────────────────────────────────────────────────────────┘
```

### Phase 3: Certificate Usage

```
┌─────────────────────────────────────────────────────────────┐
│ Kong Ingress Controller Uses Certificate                    │
│                                                              │
│ 1. Kong watches for Ingress resources                       │
│ 2. Detects tls section with secretName                      │
│ 3. Reads certificate from Secret                            │
│ 4. Configures TLS termination                               │
│ 5. Serves HTTPS traffic                                     │
└─────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────┐
│ Traffic Flow                                                 │
│                                                              │
│ Client (Browser)                                             │
│       ↓ HTTPS Request                                        │
│ Kong Ingress Controller                                      │
│       ↓ TLS Termination (using certificate)                 │
│       ↓ HTTP Request                                         │
│ Keycloak Service                                             │
│       ↓                                                      │
│ Keycloak Pod                                                 │
└─────────────────────────────────────────────────────────────┘
```

---

## Certificate Usage

### How Kong Uses the Certificate

1. **Secret Reference**: Kong reads the Ingress resource and finds the `tls.secretName`
2. **Certificate Loading**: Kong loads the certificate and private key from the Secret
3. **TLS Configuration**: Kong configures its TLS listener with the certificate
4. **Traffic Handling**: 
   - Client connects via HTTPS
   - Kong presents the certificate
   - TLS handshake completes
   - Kong decrypts traffic
   - Kong forwards HTTP traffic to backend service (Keycloak)

### Certificate Storage

Certificates are stored as Kubernetes Secrets:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: <your-domain>-tls
  namespace: keycloak
type: kubernetes.io/tls
data:
  tls.crt: <base64-encoded-certificate>
  tls.key: <base64-encoded-private-key>
  ca.crt: <base64-encoded-ca-certificate>
```

---

## Detailed Component Analysis

### 1. cert-manager Components

#### cert-manager Controller
- **Responsibility**: Main orchestrator
- **Functions**:
  - Watches Certificate resources
  - Creates CertificateRequests
  - Manages ACME protocol flow
  - Updates Secrets with issued certificates
  - Handles certificate renewal (60 days before expiry)

#### cert-manager Webhook
- **Responsibility**: Validation and mutation
- **Functions**:
  - Validates Certificate resources
  - Validates Issuer/ClusterIssuer resources
  - Ensures configuration correctness
  - Prevents invalid certificate requests

#### cert-manager CAInjector
- **Responsibility**: CA certificate injection
- **Functions**:
  - Injects CA certificates into webhooks
  - Updates ValidatingWebhookConfiguration
  - Updates MutatingWebhookConfiguration
  - Ensures webhook trust chains

### 2. Let's Encrypt ClusterIssuer Configuration

**File**: `charts/letsencrypt-issuer/templates/clusterissuer.yaml`

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-issuer-letsencrypt-clusterissuer
spec:
  acme:
    # Let's Encrypt production server
    server: https://acme-v02.api.letsencrypt.org/directory
    
    # Email for certificate expiry notifications
    email: dummysender111@gmail.com
    
    # Secret to store ACME account private key
    privateKeySecretRef:
      name: letsencrypt-production-private-key
    
    # Challenge solver configuration
    solvers:
      - http01:
          ingress:
            ingressClassName: kong
```

**Key Configuration**:
- **Environment**: Production (not staging)
- **ACME Server**: Let's Encrypt production API
- **Challenge Type**: HTTP-01 (domain validation via HTTP)
- **Ingress Class**: Kong (uses Kong for challenge validation)

### 3. Kong Ingress Controller

**Role in Certificate Flow**:
1. **Challenge Validation**: Routes ACME challenge requests to acme-solver pods
2. **Certificate Usage**: Reads certificates from Secrets and configures TLS
3. **Traffic Routing**: Routes incoming HTTPS traffic to backend services

**Configuration**:
- **IngressClass**: `kong`
- **TLS Termination**: Enabled
- **Certificate Source**: Kubernetes Secrets

### 4. Keycloak Certificate Configuration

**Ingress Configuration** (`charts/keycloak/values.yaml`):

```yaml
ingress:
  enabled: true
  ingressClassName: "kong"
  hostname: ""  # Set during deployment
  path: "/auth"
  annotations:
    konghq.com/strip-path: "true"
    konghq.com/preserve-host: "true"
    cert-manager.io/cluster-issuer: "letsencrypt-issuer-letsencrypt-clusterissuer"
  tls: false  # Set to true to enable TLS
```

**Important Annotations**:
- `cert-manager.io/cluster-issuer`: Tells cert-manager which issuer to use
- `konghq.com/strip-path`: Kong strips `/auth` before forwarding
- `konghq.com/preserve-host`: Kong preserves original Host header

---

## Configuration Details

### cert-manager Installation

**Chart Location**: `charts/cert-manager/`

**Key Settings**:
```yaml
installCRDs: true  # Install Custom Resource Definitions
crds:
  enabled: true
  keep: true  # Keep CRDs when uninstalling

global:
  logLevel: 2  # Verbosity level (0-6)
  leaderElection:
    namespace: "kube-system"

replicaCount: 1  # Number of controller replicas

prometheus:
  enabled: true  # Enable metrics
```

### Let's Encrypt Issuer Configuration

**Chart Location**: `charts/letsencrypt-issuer/`

**Configuration Options**:
```yaml
environment: production  # or "staging" for testing
email: dummysender111@gmail.com  # Certificate notifications
issuerUrl: "https://acme-v02.api.letsencrypt.org/directory"

solvers:
  - http01:
      ingress:
        ingressClassName: kong
```

**Staging vs Production**:
- **Staging**: Use for testing (higher rate limits, untrusted certificates)
  - URL: `https://acme-staging-v02.api.letsencrypt.org/directory`
- **Production**: Use for real certificates (lower rate limits, trusted certificates)
  - URL: `https://acme-v02.api.letsencrypt.org/directory`

### Certificate Lifecycle

**Certificate Validity**:
- **Duration**: 90 days (Let's Encrypt standard)
- **Renewal**: Automatic, 30 days before expiry
- **Renewal Process**: Same as initial issuance

**Renewal Flow**:
```
Day 0: Certificate issued
Day 60: cert-manager starts renewal process
Day 61-89: Retry if renewal fails
Day 90: Certificate expires (if renewal failed)
```

---

## Troubleshooting

### Common Issues and Solutions

#### 1. Certificate Not Created

**Symptoms**:
- Ingress created but no certificate
- No Certificate resource in namespace

**Checks**:
```bash
# Check if cert-manager is running
kubectl get pods -n cert-manager

# Check ClusterIssuer status
kubectl get clusterissuer
kubectl describe clusterissuer letsencrypt-issuer-letsencrypt-clusterissuer

# Check Ingress annotations
kubectl get ingress -n keycloak -o yaml
```

**Common Causes**:
- Missing `cert-manager.io/cluster-issuer` annotation
- Missing `tls` section in Ingress
- cert-manager not running

#### 2. Certificate Pending

**Symptoms**:
- Certificate resource exists but status is "Pending"
- No certificate issued

**Checks**:
```bash
# Check Certificate status
kubectl get certificate -n keycloak
kubectl describe certificate <cert-name> -n keycloak

# Check CertificateRequest
kubectl get certificaterequest -n keycloak
kubectl describe certificaterequest <request-name> -n keycloak

# Check Order
kubectl get order -n keycloak
kubectl describe order <order-name> -n keycloak

# Check Challenge
kubectl get challenge -n keycloak
kubectl describe challenge <challenge-name> -n keycloak
```

**Common Causes**:
- Domain not pointing to cluster
- Firewall blocking port 80
- Kong not routing challenge requests
- Rate limit exceeded

#### 3. Challenge Failed

**Symptoms**:
- Challenge resource shows "Failed" status
- Let's Encrypt cannot validate domain

**Checks**:
```bash
# Check challenge details
kubectl describe challenge <challenge-name> -n keycloak

# Check acme-solver pod
kubectl get pods -n keycloak | grep acme-solver
kubectl logs <acme-solver-pod> -n keycloak

# Test challenge URL manually
curl http://<your-domain>/.well-known/acme-challenge/<token>
```

**Common Causes**:
- DNS not configured correctly
- Domain not pointing to Kong
- Kong not routing to acme-solver
- Port 80 not accessible from internet

#### 4. Certificate Issued but Not Used

**Symptoms**:
- Certificate exists in Secret
- HTTPS still not working

**Checks**:
```bash
# Check Secret
kubectl get secret <your-domain>-tls -n keycloak
kubectl describe secret <your-domain>-tls -n keycloak

# Check Kong configuration
kubectl get ingress -n keycloak -o yaml

# Check Kong logs
kubectl logs -n ingress-controller <kong-pod-name>
```

**Common Causes**:
- Secret name mismatch in Ingress
- Kong not reloaded after certificate creation
- TLS not enabled in Ingress

### Debugging Commands

```bash
# View all cert-manager resources
kubectl get certificate,certificaterequest,order,challenge --all-namespaces

# Check cert-manager logs
kubectl logs -n cert-manager deployment/cert-manager -f

# Check cert-manager webhook logs
kubectl logs -n cert-manager deployment/cert-manager-webhook -f

# Check Kong logs
kubectl logs -n ingress-controller deployment/<kong-deployment> -f

# Describe Certificate for detailed status
kubectl describe certificate <cert-name> -n <namespace>

# Check events
kubectl get events -n <namespace> --sort-by='.lastTimestamp'
```

### Rate Limits

**Let's Encrypt Rate Limits**:
- **Certificates per Registered Domain**: 50 per week
- **Duplicate Certificate**: 5 per week
- **Failed Validation**: 5 failures per account, per hostname, per hour

**Recommendation**: Use staging environment for testing

---

## Summary

### Certificate Creation Process

1. **cert-manager** is deployed and watches for Ingress resources
2. **Let's Encrypt ClusterIssuer** is configured with ACME settings
3. **Kong Ingress Controller** is deployed and ready to route traffic
4. **Keycloak Ingress** is created with cert-manager annotation
5. **cert-manager** detects Ingress and creates Certificate resource
6. **ACME protocol** is initiated with Let's Encrypt
7. **HTTP-01 challenge** is created and validated via Kong
8. **Certificate** is issued and stored in Kubernetes Secret
9. **Kong** reads certificate from Secret and configures TLS
10. **HTTPS traffic** is now served with valid certificate

### Key Takeaways

- **Automated**: Entire process is automated, no manual intervention needed
- **Declarative**: Configuration is declarative via Kubernetes resources
- **Secure**: Private keys never leave the cluster
- **Renewable**: Certificates automatically renew before expiry
- **Scalable**: Works for multiple services and domains
- **Production-Ready**: Uses Let's Encrypt production environment

### Service Responsible for Certificate Creation

**Answer**: **cert-manager** is the service responsible for creating certificates. Specifically:
- **cert-manager controller** orchestrates the entire certificate lifecycle
- It communicates with **Let's Encrypt** (via ACME protocol) to obtain certificates
- It uses **Kong Ingress Controller** to complete domain validation challenges
- It stores the issued certificates in **Kubernetes Secrets**

The certificate in the Keycloak namespace is created automatically when the Keycloak Ingress resource is deployed with the appropriate cert-manager annotations.
