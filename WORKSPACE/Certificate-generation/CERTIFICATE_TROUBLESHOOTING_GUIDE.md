# Certificate Management Troubleshooting Guide

## Table of Contents
1. [Overview](#overview)
2. [How Certificates Get Created](#how-certificates-get-created)
3. [Working Configuration Example](#working-configuration-example)
4. [Common Issues and Solutions](#common-issues-and-solutions)
5. [Configuration Comparison](#configuration-comparison)
6. [Best Practices](#best-practices)

---

## Overview

This guide explains how TLS certificates are automatically created in the Finternet Kubernetes infrastructure using cert-manager and Let's Encrypt. It covers the complete flow, common issues, and solutions based on real troubleshooting scenarios.

### Key Components
- **cert-manager**: Kubernetes certificate management controller
- **Let's Encrypt**: Free automated Certificate Authority
- **Kong Ingress Controller**: API Gateway handling TLS termination
- **ClusterIssuer**: Configuration for obtaining certificates from Let's Encrypt

---

## How Certificates Get Created

### Method 1: Ingress-Driven (Automatic)

cert-manager has an **Ingress-Shim controller** that automatically creates Certificate resources when it detects Ingress resources with specific annotations.

#### Requirements for Automatic Certificate Creation

**Three conditions must be met:**

1. **Ingress must exist** (`ingress.enabled: true`)
2. **cert-manager annotation must be present**
3. **TLS section must be configured with hosts and secretName**

#### The Ingress-Shim Controller Logic

cert-manager controller (running as a pod) continuously watches all Ingress resources:

```go
// Pseudo-code of cert-manager's Ingress-Shim controller
if ingress.HasAnnotation("cert-manager.io/cluster-issuer") {
    if ingress.Spec.TLS != nil && len(ingress.Spec.TLS) > 0 {
        for each tlsConfig in ingress.Spec.TLS {
            // Automatically create Certificate resource
            certificate := CreateCertificate{
                Name: tlsConfig.SecretName,
                Namespace: ingress.Namespace,
                SecretName: tlsConfig.SecretName,
                IssuerRef: ingress.Annotations["cert-manager.io/cluster-issuer"],
                DNSNames: tlsConfig.Hosts,
            }
            kubernetes.Create(certificate)
        }
    }
}
```

**Key Point:** This is **runtime code logic** in the cert-manager binary, not a configuration file or template.

#### What Triggers Certificate Creation

When you deploy an Ingress like this:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-app
  namespace: my-namespace
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-issuer-letsencrypt-clusterissuer"
spec:
  ingressClassName: kong
  rules:
    - host: "example.com"
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: my-app
                port:
                  number: 80
  tls:  # ← This section triggers cert-manager
    - hosts:
        - "example.com"
      secretName: "example-com-tls"
```

cert-manager **automatically creates** this Certificate resource:

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: example-com-tls
  namespace: my-namespace
  ownerReferences:  # Links to parent Ingress
    - apiVersion: networking.k8s.io/v1
      kind: Ingress
      name: my-app
spec:
  secretName: example-com-tls
  issuerRef:
    name: letsencrypt-issuer-letsencrypt-clusterissuer
    kind: ClusterIssuer
    group: cert-manager.io
  dnsNames:
    - example.com
```

### Method 2: Standalone Certificate (Manual)

Create a Certificate resource directly without depending on Ingress:

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: my-certificate
  namespace: default
spec:
  secretName: my-domain-tls
  issuerRef:
    name: letsencrypt-issuer-letsencrypt-clusterissuer
    kind: ClusterIssuer
  dnsNames:
    - example.com
    - www.example.com
```

**Advantages:**
- Independent of Ingress state
- Can be shared across multiple Ingresses
- Centralized certificate management
- Works even when Ingress is disabled

---

## Working Configuration Example

### Keycloak Configuration (Working Setup)

#### ArgoCD Override File
**Location:** `/fresh/units-automation/helmchart/argocd-application/values/applications/keycloak.yaml`

```yaml
{{- $globals := .Values.global | default dict -}}
{{- $auth := $globals.auth | default dict -}}
{{- $keycloak := $auth.keycloak | default dict -}}
ingress:
  enabled: true  # ✓ Ingress is enabled
  hostname: {{ $globals.domain | default "foundry.finternetlab.io" }}  # ✓ Hostname set
  tls: true  # ✓ TLS enabled (boolean flag)
  annotations:
    cert-manager.io/cluster-issuer: {{ $globals.clusterIssuer | default "letsencrypt-issuer-letsencrypt-clusterissuer" }}  # ✓ cert-manager annotation

auth:
  adminUser: {{ $keycloak.adminUser | default "user" }}
  adminPassword: {{ $keycloak.adminPassword | default "admin" | quote }}

postgresql:
  enabled: false

externalDatabase:
  host: "{{ .Values.global.database.yugabyte.host }}"
  port: "{{ .Values.global.database.yugabyte.port }}"
  user: "{{ .Values.global.database.keycloak.user }}"
  database: "{{ .Values.global.database.keycloak.database }}"
  password: "{{ .Values.global.database.keycloak.password }}"
```

#### Keycloak Ingress Template Logic
**Location:** `helmcharts/charts/keycloak/templates/ingress.yaml` (Lines 44-52)

```yaml
{{- if or (and .Values.ingress.tls 
           (or (include "common.ingress.certManagerRequest" (dict "annotations" .Values.ingress.annotations)) 
               .Values.ingress.selfSigned 
               .Values.ingress.secrets)) 
       .Values.ingress.extraTls }}
  tls:
  {{- if and .Values.ingress.tls (or (include "common.ingress.certManagerRequest" (dict "annotations" .Values.ingress.annotations)) .Values.ingress.secrets .Values.ingress.selfSigned) }}
    - hosts:
        - {{ (tpl .Values.ingress.hostname .) | quote }}
      secretName: {{ printf "%s-tls" (tpl .Values.ingress.hostname .) }}  # Auto-generates: hostname + "-tls"
  {{- end }}
{{- end }}
```

**Template Logic:**
1. Checks if `tls: true` (boolean)
2. Checks if cert-manager annotation exists
3. Auto-generates TLS section with:
   - `hosts`: From `ingress.hostname`
   - `secretName`: Automatically generated as `<hostname>-tls`

#### Resulting Kubernetes Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: keycloak
  namespace: keycloak
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-issuer-letsencrypt-clusterissuer"
    konghq.com/strip-path: "true"
    konghq.com/preserve-host: "true"
spec:
  ingressClassName: kong
  rules:
    - host: "foundry.finternetlab.io"
      http:
        paths:
          - path: /auth
            pathType: Prefix
            backend:
              service:
                name: keycloak
                port:
                  number: 80
  tls:  # ← Generated by template, triggers cert-manager
    - hosts:
        - "foundry.finternetlab.io"
      secretName: "foundry.finternetlab.io-tls"
```

#### Certificate Creation Flow

1. **Ingress deployed** with cert-manager annotation and TLS section
2. **cert-manager Ingress-Shim controller detects** the Ingress
3. **Certificate resource automatically created**:
   ```yaml
   apiVersion: cert-manager.io/v1
   kind: Certificate
   metadata:
     name: foundry.finternetlab.io-tls
     namespace: keycloak
   spec:
     secretName: foundry.finternetlab.io-tls
     issuerRef:
       name: letsencrypt-issuer-letsencrypt-clusterissuer
       kind: ClusterIssuer
     dnsNames:
       - foundry.finternetlab.io
   ```
4. **ACME challenge initiated** with Let's Encrypt
5. **HTTP-01 challenge** validated via Kong
6. **Certificate issued** and stored in Secret `foundry.finternetlab.io-tls`
7. **Kong reads certificate** from Secret and configures TLS
8. **HTTPS traffic** now works

---

## Common Issues and Solutions

### Issue 1: Certificate Not Created - Ingress Disabled

#### Symptoms
- No Certificate resource created
- No cert-manager activity in logs
- Ingress doesn't exist

#### Root Cause
```yaml
# ArgoCD override
ingress:
  enabled: false  # ← Ingress not created
```

**Why it fails:**
- No Ingress resource exists in Kubernetes
- cert-manager has nothing to watch
- No Certificate created

#### Solution
Enable Ingress with proper configuration:
```yaml
ingress:
  enabled: true
  hostname: "your-domain.com"
  tls: true
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-issuer-letsencrypt-clusterissuer"
```

---

### Issue 2: Certificate Not Created - Empty TLS Section

#### Symptoms
- Ingress exists
- cert-manager annotation present
- But no Certificate created

#### Root Cause (Finternet-app Example)

**ArgoCD Override:**
```yaml
ingress:
  enabled: true
  hostname: "foundry.finternetlab.io"
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-issuer-letsencrypt-clusterissuer"
  hosts:
    - host: "foundry.finternetlab.io"
      paths:
        - path: /
          pathType: Prefix
  tls: []  # ← Empty array! No TLS configuration
```

**Resulting Ingress:**
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-issuer-letsencrypt-clusterissuer"
spec:
  rules:
    - host: "foundry.finternetlab.io"
      http:
        paths:
          - path: /
  # No tls section! ← cert-manager ignores this Ingress
```

**Why it fails:**
- Ingress exists ✓
- Annotation exists ✓
- But `spec.tls` is missing ✗
- cert-manager requires TLS section to create Certificate

#### Solution
Provide proper TLS configuration:
```yaml
{{- $domain := .Values.global.domain | default "foundry.finternetlab.io" -}}
ingress:
  enabled: true
  hostname: {{ $domain }}
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-issuer-letsencrypt-clusterissuer"
  hosts:
    - host: {{ $domain | quote }}
      paths:
        - path: /
          pathType: Prefix
  tls:  # ← Proper TLS configuration
    - hosts:
        - {{ $domain | quote }}
      secretName: {{ printf "%s-tls" $domain }}
```

---

### Issue 3: Certificate Not Created - Missing Annotation

#### Symptoms
- Ingress exists with TLS section
- But no Certificate created

#### Root Cause
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-app
  # Missing cert-manager annotation!
spec:
  tls:
    - hosts:
        - "example.com"
      secretName: "example-com-tls"
```

**Why it fails:**
- cert-manager only watches Ingresses with its annotation
- Without annotation, cert-manager ignores the Ingress

#### Solution
Add cert-manager annotation:
```yaml
metadata:
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-issuer-letsencrypt-clusterissuer"
```

---

### Issue 4: Certificate Pending - Challenge Failed

#### Symptoms
- Certificate resource created
- Status shows "Pending"
- Challenge fails

#### Common Causes

**1. DNS Not Configured**
```bash
# Check DNS resolution
nslookup your-domain.com
dig your-domain.com

# Should point to your cluster's external IP
```

**2. Port 80 Not Accessible**
```bash
# Test HTTP access
curl http://your-domain.com/.well-known/acme-challenge/test

# Should reach Kong (even if 404)
```

**3. Kong Not Routing Challenge**
```bash
# Check acme-solver pod
kubectl get pods -n <namespace> | grep acme-solver

# Check challenge Ingress
kubectl get ingress -n <namespace> | grep cm-acme-http-solver
```

#### Solution Steps

1. **Verify DNS:**
   ```bash
   # DNS must point to Kong's external IP
   kubectl get svc -n ingress-controller
   ```

2. **Check Challenge Resources:**
   ```bash
   kubectl get challenge -n <namespace>
   kubectl describe challenge <challenge-name> -n <namespace>
   ```

3. **Check cert-manager Logs:**
   ```bash
   kubectl logs -n cert-manager deployment/cert-manager -f
   ```

4. **Test Challenge URL:**
   ```bash
   # Get challenge token from Challenge resource
   kubectl get challenge <challenge-name> -n <namespace> -o yaml
   
   # Test URL
   curl http://your-domain.com/.well-known/acme-challenge/<token>
   ```

---

## Configuration Comparison

### Keycloak (Working) vs Finternet-app (Not Working)

| Aspect | Keycloak (✓ Working) | Finternet-app (✗ Not Working) |
|--------|---------------------|-------------------------------|
| **Ingress Enabled** | `enabled: true` | `enabled: true` |
| **Hostname** | `hostname: "foundry.finternetlab.io"` | `hostname: "foundry.finternetlab.io"` |
| **Annotation** | `cert-manager.io/cluster-issuer: "..."` | `cert-manager.io/cluster-issuer: "..."` |
| **TLS Config** | `tls: true` (boolean) | `tls: []` (empty array) |
| **Template Logic** | Smart: Auto-generates TLS section | Simple: Passthrough only |
| **Result** | TLS section created ✓ | No TLS section ✗ |
| **Certificate** | Created automatically ✓ | Not created ✗ |

### Key Difference

**Keycloak Template (Smart):**
```yaml
# Input: tls: true (boolean)
# Template auto-generates:
tls:
  - hosts:
      - "foundry.finternetlab.io"
    secretName: "foundry.finternetlab.io-tls"
```

**Finternet-app Template (Simple):**
```yaml
# Input: tls: [] (empty array)
# Template outputs: (nothing, empty)
# No TLS section in Ingress
```

---

## Best Practices

### 1. Use Standalone Certificates for Production

Instead of relying on Ingress-driven certificates, create standalone Certificate resources:

**Create Certificate Template:**
```yaml
# charts/letsencrypt-issuer/templates/certificate.yaml
{{- if .Values.certificate.enabled }}
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: {{ .Values.certificate.name }}
  namespace: {{ .Values.certificate.namespace }}
spec:
  secretName: {{ .Values.certificate.secretName }}
  issuerRef:
    name: {{ include "letsencrypt-issuer.fullname" . }}-letsencrypt-clusterissuer
    kind: ClusterIssuer
  dnsNames:
    {{- range .Values.certificate.dnsNames }}
    - {{ . | quote }}
    {{- end }}
{{- end }}
```

**Configure in Values:**
```yaml
# values/applications/letsencrypt-issuer.yaml
certificate:
  enabled: true
  name: "finternet-certificate"
  namespace: "default"
  secretName: "foundry.finternetlab.io-tls"
  dnsNames:
    - "foundry.finternetlab.io"
    - "*.foundry.finternetlab.io"  # Wildcard
```

**Benefits:**
- Certificate exists regardless of Ingress state
- Single certificate for multiple services
- Centralized management
- Can use wildcard certificates

### 2. Reference Certificate in Ingress

Any Ingress can reference the standalone certificate:

```yaml
ingress:
  enabled: true
  tls:
    - hosts:
        - "foundry.finternetlab.io"
      secretName: "foundry.finternetlab.io-tls"  # References standalone cert
```

### 3. Use Staging for Testing

```yaml
# letsencrypt-issuer values
environment: staging  # Use staging for testing
issuerUrl: "https://acme-staging-v02.api.letsencrypt.org/directory"
```

**Rate Limits:**
- **Production:** 50 certs/week per domain
- **Staging:** Much higher limits, untrusted certs

### 4. Monitor Certificate Expiry

```bash
# Check certificate status
kubectl get certificate --all-namespaces

# Check expiry date
kubectl get secret <cert-secret> -n <namespace> -o jsonpath='{.data.tls\.crt}' | base64 -d | openssl x509 -noout -dates
```

### 5. Enable cert-manager Logging

```yaml
# cert-manager values
global:
  logLevel: 2  # 0-6, higher = more verbose
```

---

## Verification Commands

### Check Certificate Creation

```bash
# 1. Check if Ingress has correct annotation and TLS section
kubectl get ingress <ingress-name> -n <namespace> -o yaml

# 2. Watch for Certificate creation
kubectl get certificate -n <namespace> -w

# 3. Check Certificate status
kubectl describe certificate <cert-name> -n <namespace>

# 4. Check CertificateRequest
kubectl get certificaterequest -n <namespace>

# 5. Check ACME Order
kubectl get order -n <namespace>

# 6. Check Challenge
kubectl get challenge -n <namespace>
kubectl describe challenge <challenge-name> -n <namespace>

# 7. Check Secret created
kubectl get secret <cert-secret> -n <namespace>

# 8. Verify certificate content
kubectl get secret <cert-secret> -n <namespace> -o jsonpath='{.data.tls\.crt}' | base64 -d | openssl x509 -text -noout
```

### Check cert-manager Health

```bash
# Check cert-manager pods
kubectl get pods -n cert-manager

# Check cert-manager logs
kubectl logs -n cert-manager deployment/cert-manager -f

# Check webhook logs
kubectl logs -n cert-manager deployment/cert-manager-webhook -f

# Check ClusterIssuer
kubectl get clusterissuer
kubectl describe clusterissuer letsencrypt-issuer-letsencrypt-clusterissuer
```

---

## Summary

### Certificate Creation Requirements

For cert-manager to automatically create a certificate from an Ingress:

1. ✓ **Ingress must exist** (`ingress.enabled: true`)
2. ✓ **cert-manager annotation** must be present
3. ✓ **TLS section** must have hosts and secretName
4. ✓ **cert-manager** must be running
5. ✓ **ClusterIssuer** must be configured
6. ✓ **Domain DNS** must point to cluster
7. ✓ **Port 80** must be accessible for HTTP-01 challenge

### Why Certificates Weren't Created

**Keycloak (Initially):**
- ArgoCD override had `ingress.enabled: false`
- No Ingress → No certificate

**Finternet-app:**
- ArgoCD override had `tls: []` (empty)
- Ingress created but no TLS section → No certificate

### Solution

**Option 1: Fix Ingress Configuration**
```yaml
ingress:
  enabled: true
  hostname: "your-domain.com"
  tls: true  # or proper array
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-issuer-letsencrypt-clusterissuer"
```

**Option 2: Use Standalone Certificate (Recommended)**
```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: my-cert
spec:
  secretName: domain-tls
  issuerRef:
    name: letsencrypt-clusterissuer
    kind: ClusterIssuer
  dnsNames:
    - your-domain.com
```

### Key Takeaway

**cert-manager does NOT have a configuration file that creates certificates.** It uses **runtime controller logic** that watches Ingress resources and automatically creates Certificate resources when it detects:
- Ingress with cert-manager annotation
- TLS section with hosts and secretName

The certificate creation is **automatic behavior** of the cert-manager controller, not a template or configuration file.
