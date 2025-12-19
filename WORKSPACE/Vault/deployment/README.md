## Vault Deployment Guide (Dev & Production)

This folder contains Helm assets for deploying HashiCorp Vault for the UNITS service:

- The **official Vault Helm chart** already vendored in this repo at `vault/`
- A **`vault-bootstrap` Helm chart** that configures secrets engines, JWT auth with Keycloak, policies, and AppRole
- **Override values** for dev and production deployments of the main Vault chart

The goal is:

- **Do not modify** the upstream `vault/` chart
- Use small, clear **override files** for dev/prod
- Use a separate, templated **bootstrap chart** for Vault runtime configuration

---

## Repo Structure (Vault-related)

- `vault/`  
  Official HashiCorp Vault Helm chart (vendored).

- `deployment/`  
  - `values/`
    - `vault-values-dev.yaml` – overrides for dev/local Vault
    - `vault-values-prod.yaml` – overrides for production Vault
  - `vault-bootstrap/`
    - Helm chart that:
      - Creates a ConfigMap with policies
      - Runs a Job that:
        - Enables Transit and KV v2 secrets engines
        - Creates `transit/keys/units-encryption-master`
        - Configures JWT auth with Keycloak
        - Applies policies
        - Configures an AppRole for the service

---

## Prerequisites

- A Kubernetes cluster (for dev: kind/minikube/etc., for prod: managed cluster)
- `kubectl` installed and pointing to the cluster
- `helm` v3.6+ installed
- For production:
  - Some access to store Vault unseal keys and bootstrap token securely
  - Keycloak reachable from the cluster for JWT auth

---

## Overview of Charts and Sequencing

- **Chart 1: `vault/` (main Vault server)**
  - Deployed with environment-specific values files:
    - `deployment/values/vault-values-dev.yaml`
    - `deployment/values/vault-values-prod.yaml`

- **Chart 2: `deployment/vault-bootstrap/`**
  - Deployed after Vault is running (and unsealed in non-dev)
  - Uses a Vault bootstrap token stored in a Kubernetes Secret
  - Configures secrets engines, JWT auth, policies, and AppRole

**Can we deploy both charts at the same time?**

- They are **separate Helm releases**, so they are installed with **two `helm install` commands**.
- In CI/CD, you typically:
  1. Install/upgrade `vault` chart.
  2. Wait for Vault pods Ready.
  3. Initialize & unseal Vault (prod only), or rely on dev mode.
  4. Install/upgrade `vault-bootstrap` chart.
- Helm will not "merge" them into a single release, but you can run the install commands back-to-back in a pipeline.

---

## Dev Deployment (Local / Non-HA)

This mode is similar to the existing Docker Compose dev setup: Vault dev server with an in-memory backend and a fixed root token.

### Step 1: Install Vault (dev mode)

From the repo root:

```bash
helm install vault ./vault \
  --namespace vault \
  --create-namespace \
  -f deployment/values/vault-values-dev.yaml
```

What this does (via `vault-values-dev.yaml`):

- Runs Vault in **dev mode** (`server.dev.enabled=true`)
- Sets `devRootToken` to `dev-token`
- Enables the UI inside the cluster
- Disables TLS (HTTP only) – for local/dev usage only

In dev mode:

- Vault is **auto-initialized and auto-unsealed**
- You do **not** run `vault operator init` or `vault operator unseal`
- The root token is fixed (`dev-token`)

### Step 2: Create a bootstrap token Secret

Even in dev, we use the same pattern as production: a Kubernetes Secret that holds the token which the bootstrap Job will use.

```bash
kubectl -n vault create secret generic vault-bootstrap-token \
  --from-literal=VAULT_TOKEN=dev-token
```

> In production you will use a different token value here (root or restricted admin token).

### Step 3: Install `vault-bootstrap` chart

From the repo root:

```bash
helm install vault-bootstrap ./deployment/vault-bootstrap \
  --namespace vault
```

What happens:

- A ConfigMap with three policies is created:
  - `user-transit-keys`
  - `user-pii-encryption`
  - `service-admin`
- A Job runs that:
  - Waits for Vault to be reachable at `VAULT_ADDR`
  - Enables Transit and KV v2 engines
  - Creates the `units-encryption-master` key
  - Configures JWT auth with Keycloak
  - Writes the three policies
  - Creates an AppRole `service-app` with `service-admin` policy

You can watch the Job:

```bash
kubectl -n vault logs job/vault-bootstrap
kubectl -n vault get job vault-bootstrap
```

Once the Job succeeds, the Vault configuration matches the behavior described in `workflow/docs/VAULT_JWT_SETUP_GUIDE.md`.

---

## Production Deployment (Kubernetes)

In production, Vault should run in HA mode with persistent storage (and usually TLS). The steps below follow the JWT setup guide and keep the same bootstrap pattern.

### Step 1: Install Vault (HA, non-dev)

From the repo root:

```bash
helm install vault ./vault \
  --namespace vault \
  --create-namespace \
  -f deployment/values/vault-values-prod.yaml
```

What `vault-values-prod.yaml` configures:

- Disables dev mode (`server.dev.enabled=false`)
- Enables HA with integrated Raft storage
- Enables persistent data storage via PVCs
- Enables the Vault UI service
- Optionally enables the agent injector (can be toggled per environment)

At this point:

- Vault pods are running but **sealed** and **not yet initialized**.

### Step 2: Initialize Vault

Run the following **once** (from your workstation, with `kubectl` pointing at the cluster):

```bash
kubectl exec -n vault vault-0 -- \
  vault operator init \
    -key-shares=5 \
    -key-threshold=3 \
    -format=json > vault-keys.json
```

This:

- Initializes Vault
- Outputs:
  - 5 unseal keys (in `unseal_keys_b64[]`)
  - 1 root token

> **Important:** Store `vault-keys.json` securely in a real secret store or HSM-backed system. Do not commit it to Git.

### Step 3: Unseal Vault pods

Extract three unseal keys:

```bash
UNSEAL_KEY_1=$(jq -r '.unseal_keys_b64[0]' vault-keys.json)
UNSEAL_KEY_2=$(jq -r '.unseal_keys_b64[1]' vault-keys.json)
UNSEAL_KEY_3=$(jq -r '.unseal_keys_b64[2]' vault-keys.json)
```

Unseal each Vault pod (e.g., `vault-0`, `vault-1`, `vault-2`):

```bash
for pod in vault-0 vault-1 vault-2; do
  kubectl exec -n vault $pod -- vault operator unseal "$UNSEAL_KEY_1"
  kubectl exec -n vault $pod -- vault operator unseal "$UNSEAL_KEY_2"
  kubectl exec -n vault $pod -- vault operator unseal "$UNSEAL_KEY_3"
done
```

After this, `vault status` should show `sealed=false` on each pod.

### Step 4: Create a bootstrap token Secret

From `vault-keys.json`, get the root token:

```bash
ROOT_TOKEN=$(jq -r '.root_token' vault-keys.json)
```

Create a Secret in the `vault` namespace:

```bash
kubectl -n vault create secret generic vault-bootstrap-token \
  --from-literal=VAULT_TOKEN="$ROOT_TOKEN"
```

> For stricter security, you can create a **dedicated admin token** with limited policies and use that instead of the full root token. The bootstrap chart simply needs any token with enough privileges to:
> - Enable secrets/auth engines
> - Write policies
> - Manage AppRole

### Step 5: Install `vault-bootstrap` chart

```bash
helm install vault-bootstrap ./deployment/vault-bootstrap \
  --namespace vault
```

The Job will:

- Wait until `vault status` succeeds (which requires Vault to be unsealed)
- Apply the same configuration as in dev:
  - Transit + KV engines
  - `units-encryption-master` key
  - JWT auth with Keycloak
  - Policies (`user-transit-keys`, `user-pii-encryption`, `service-admin`)
  - `service-app` AppRole
  - AppRole `role_id` and `secret_id` stored in `secret/approle/service-app` (KV)

You can re-run it in the future with:

```bash
helm upgrade --install vault-bootstrap ./deployment/vault-bootstrap \
  --namespace vault
```

---

## `vault-bootstrap` Chart Internals and Flow

The `deployment/vault-bootstrap` chart encapsulates all **runtime** Vault configuration that would normally be done with manual `vault` CLI commands.

### Components

- **`values.yaml`**
  - `vault.address`: where the Job contacts Vault (default `http://vault.vault.svc:8200`).
  - `vault.bootstrapTokenSecret`: reference to a Secret containing an admin token (`VAULT_TOKEN`).
  - `vault.keycloak.url`, `vault.keycloak.realm`, `vault.keycloak.clientId`: used to build the Keycloak JWKS URL.
  - `policies.*`: HCL bodies for the three policies.
- **`templates/policies-configmap.yaml`**
  - Creates a ConfigMap `<release>-vault-bootstrap-policies` with:
    - `user-transit-keys.hcl`
    - `user-pii-encryption.hcl`
    - `service-admin.hcl`
- **`templates/bootstrap-job.yaml`**
  - Creates a Job `vault-bootstrap-vault-bootstrap` that:
    - Runs a `hashicorp/vault` container with:
      - `VAULT_ADDR` set from `values.yaml`
      - `VAULT_TOKEN` read from `vault-bootstrap-token` Secret
      - `KEYCLOAK_URL`, `KEYCLOAK_REALM`, `KEYCLOAK_CLIENT_ID` from values
    - Mounts the policies ConfigMap at `/vault/policies`

### Execution Flow

On each install/upgrade of `vault-bootstrap`:

1. **Wait for Vault**
   - Loops on `vault status` until Vault is reachable and unsealed.

2. **Enable secrets engines**
   - `vault secrets enable -path=transit transit`
   - `vault secrets enable -path=secret kv-v2`
   - Both commands are idempotent; if the engines already exist, Vault returns an error but the script continues.

3. **Create master key**
   - `vault write -f transit/keys/units-encryption-master type=aes256-gcm96`
   - Also idempotent; an existing key is treated as a no-op.

4. **Configure JWT auth (Keycloak)**
   - If `KEYCLOAK_URL` and `KEYCLOAK_REALM` are non-empty:
     - `vault auth enable jwt`
     - `vault write auth/jwt/config jwks_url=<Keycloak JWKS URL> default_role="finternet-user"`
     - `vault write auth/jwt/role/finternet-user ... token_policies="user-transit-keys,user-pii-encryption" ...`
   - This wires Vault’s `jwt/` auth method to your Keycloak realm, matching `workflow/docs/VAULT_JWT_SETUP_GUIDE.md`.

5. **Apply policies**
   - `vault policy write user-transit-keys /vault/policies/user-transit-keys.hcl`
   - `vault policy write user-pii-encryption /vault/policies/user-pii-encryption.hcl`
   - `vault policy write service-admin /vault/policies/service-admin.hcl`

6. **Configure AppRole for the service**
   - `vault auth enable approle`
   - `vault write auth/approle/role/service-app token_policies="service-admin" ...`
   - Creates an AppRole that the UNITS service uses for service-level access to Vault.

7. **Store AppRole credentials in KV**
   - Reads the AppRole `role_id`:
     - `vault read -field=role_id auth/approle/role/service-app/role-id`
   - Generates a `secret_id`:
     - `vault write -field=secret_id -f auth/approle/role/service-app/secret-id`
   - Stores both under:
     - `vault kv put secret/approle/service-app role_id="<ROLE_ID>" secret_id="<SECRET_ID>"`
   - This mirrors the docker-compose `start-and-init-vault.sh` behavior so that the application can retrieve AppRole credentials from a fixed KV location.

The Job is designed to be **idempotent**. Running:

```bash
helm upgrade --install vault-bootstrap ./deployment/vault-bootstrap -n vault
```

multiple times will reconcile configuration without breaking existing engines or keys, as long as Vault is reachable and, for JWT, the Keycloak JWKS URL is valid.

---

## Verification (Dev and Prod)

After deploying both charts, you can use the following commands to verify the complete bootstrap in either dev or prod.

### Dev (server.dev.enabled=true)

In dev, Vault uses `dev-token` as the root/admin token:

```bash
kubectl exec -n vault vault-0 -- \
  env VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=dev-token \
  sh -c '
    echo "== auth methods ==";
    vault auth list;

    echo "== jwt config ==";
    vault read auth/jwt/config;

    echo "== jwt role finternet-user ==";
    vault read auth/jwt/role/finternet-user;

    echo "== policies ==";
    vault policy list;
    vault policy read service-admin;
    vault policy read user-transit-keys;
    vault policy read user-pii-encryption;

    echo "== approle role service-app ==";
    vault list auth/approle/role;
    vault read auth/approle/role/service-app;

    echo "== approle creds KV ==";
    vault kv get secret/approle/service-app;
  '
```

You should see:

- `jwt/`, `approle/`, `token/` in `vault auth list`
- A populated JWKS URL in `auth/jwt/config`
- `token_policies=[user-transit-keys user-pii-encryption]` for `auth/jwt/role/finternet-user`
- The three custom policies in `vault policy list`
- `auth/approle/role/service-app` with `token_policies=[service-admin]`
- `secret/approle/service-app` with `role_id` and `secret_id`

### Prod (HA, non-dev)

In prod, use the root/admin token from `vault-keys.json`:

```bash
ROOT_TOKEN=$(jq -r '.root_token' vault-keys.json)

kubectl exec -n vault vault-0 -- \
  env VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN="$ROOT_TOKEN" \
  sh -c '
    echo "== status ==";
    vault status;

    echo "== auth methods ==";
    vault auth list;

    echo "== jwt config ==";
    vault read auth/jwt/config;

    echo "== jwt role finternet-user ==";
    vault read auth/jwt/role/finternet-user;

    echo "== policies ==";
    vault policy list;
    vault policy read service-admin;
    vault policy read user-transit-keys;
    vault policy read user-pii-encryption;

    echo "== approle role service-app ==";
    vault list auth/approle/role;
    vault read auth/approle/role/service-app;

    echo "== approle creds KV ==";
    vault kv get secret/approle/service-app;
  '
```

These commands confirm that the production Vault cluster is:

- Initialized and unsealed
- Exposing the same auth methods, policies, AppRole, and KV layout as the dev setup
- Correctly integrated with Keycloak via the JWT auth method

---

## Auto-Unseal in Production (Optional)

The current flow uses **Shamir unseal keys** and manual unseal via `vault operator unseal`. For many production setups, you may prefer **auto-unseal** using a cloud KMS (GCP KMS, AWS KMS, Azure Key Vault, etc.).

The upstream `vault/values.yaml` already includes example `seal` stanzas (commented) for auto-unseal. To enable auto-unseal:

1. Choose a KMS provider (e.g., GCP KMS).
2. Create and authorize a service account / IAM principal that Vault pods will use.
3. Provide necessary credentials to Vault via `extraSecretEnvironmentVars` and/or mounted secrets.
4. Add a `seal` block to the Vault server config (for HA + raft, usually under `server.ha.raft.config`).

Example (GCP KMS, simplified, in an environment-specific values file):

```yaml
server:
  ha:
    enabled: true
    raft:
      enabled: true
      config: |
        ui = true

        listener "tcp" {
          tls_disable = 1
          address = "[::]:8200"
          cluster_address = "[::]:8201"
        }

        storage "raft" {
          path = "/vault/data"
        }

        service_registration "kubernetes" {}

        seal "gcpckms" {
          project    = "your-gcp-project"
          region     = "global"
          key_ring   = "vault-unseal-kr"
          crypto_key = "vault-unseal-key"
        }

  extraSecretEnvironmentVars:
    - envName: GOOGLE_APPLICATION_CREDENTIALS
      secretName: vault-gcp-creds
      secretKey: service-account.json
```

When auto-unseal is correctly configured:

- Vault will **auto-initialize** (depending on config) and **auto-unseal** using the KMS.
- You no longer need to manually run `vault operator unseal` on each pod.
- You **still need a bootstrap token** (root or admin) for the `vault-bootstrap` chart.

> Implementing auto-unseal requires cloud-specific setup. The snippet above is a starting point; follow the official Vault documentation for your cloud provider for a complete, secure configuration.

---

## Summary

- Use the **vendored `vault/` chart** as-is, with environment-specific override values under `deployment/values/`.
- Use the **`vault-bootstrap` chart** to encapsulate all runtime Vault configuration (engines, auth methods, policies, AppRole).
- In dev:
  - Dev mode (auto-init/unseal) + `dev-token` in `vault-bootstrap-token` Secret.
- In prod:
  - Install Vault (HA), manually init + unseal, create a bootstrap token Secret, then install `vault-bootstrap`.
- Optionally, configure **auto-unseal** via a KMS provider in the Vault server config for fully automated restarts.


