| Pipeline Stage (Order)                | Primary Tool (Common Choice) | What It Does                            | Alternative Tools                     |
| ------------------------------------- | ---------------------------- | --------------------------------------- | ------------------------------------- |
| **1. Source Code Management**         | **GitHub**                   | Hosts source code, PRs, version control | GitLab, Bitbucket                     |
| **2. CI Orchestration**               | **GitHub Actions**           | Triggers pipelines on commits, PRs      | GitLab CI/CD, Jenkins, CircleCI       |
| **3. Build**                          | **Docker**                   | Builds application artifacts / images   | Podman, Buildah                       |
| **4. Test**                           | **JUnit**                    | Runs unit & integration tests           | pytest, TestNG                        |
| **5. Artifact Registry**              | **Docker Hub**               | Stores built images / artifacts         | Amazon ECR, GitHub Container Registry |
| **6. Infrastructure Provisioning**    | **Terraform**                | Creates cloud infra (VPC, EKS, DB)      | Pulumi, AWS CloudFormation            |
| **7. Configuration Management**       | **Ansible**                  | Configures VMs, installs dependencies   | Chef, Puppet                          |
| **8. Kubernetes Packaging**           | **Helm**                     | Templates Kubernetes manifests          | Kustomize                             |
| **9. Continuous Deployment (GitOps)** | **Argo CD**                  | Syncs Git → Kubernetes                  | Flux CD                               |
| **10. Runtime Platform**              | **Kubernetes**               | Runs workloads in production            | Amazon ECS, Nomad                     |
| **11. Release Strategy**              | **Argo Rollouts**            | Canary / blue-green deployments         | Flagger, Spinnaker                    |
| **12. Monitoring & Alerts**           | **Prometheus**               | Metrics & alerting                      | Datadog, New Relic                    |
| **13. Logging**                       | **ELK Stack**                | Centralized logs                        | Loki, Splunk                          |
--------------------------------------------------------------------------------------------------------------------------

| Layer    | Tool                    |
| -------- | ----------------------- |
| SCM + CI | GitHub + GitHub Actions |
| Build    | Docker                  |
| Infra    | Terraform               |
| Deploy   | Argo CD                 |
| Runtime  | Kubernetes              |
| Observe  | Prometheus + Grafana    |
