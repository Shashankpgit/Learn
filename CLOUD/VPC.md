1. Look into VPC peering.

![Image](https://edrawcloudpublicus.s3.amazonaws.com/work/1905656/2022-4-24/1650795720/main.png)

![Image](https://d2908q01vomqb2.cloudfront.net/da4b9237bacccdf19c0760cab7aec4a8359010b0/2020/03/19/Slide1.png)

![Image](https://docs.aws.amazon.com/images/vpc/latest/userguide/images/vpc-example-private-subnets.png)

![Image](https://docs.aws.amazon.com/vpc/latest/userguide/images/nat-gateway-diagram.png)

![Image](https://miro.medium.com/1%2AkXh27DyJkCv-qdY2z9yNGw.png)

Below is a **complete, structured set of notes on Amazon VPC**, designed to take you from **beginner to advanced / production-grade understanding**. The content is aligned with **real AWS usage, certifications, and interviews**.

---

# AWS VPC (Virtual Private Cloud) — Complete Notes

## 1. What is a VPC?

A **Virtual Private Cloud (VPC)** is a **logically isolated virtual network** inside **Amazon Web Services**, where you define:

* IP address range (CIDR)
* Subnets
* Routing
* Security boundaries
* Connectivity to the internet and on-prem networks

**Think of a VPC as your own data center network inside AWS.**

---

## 2. Core Networking Concepts (Foundation)

### 2.1 CIDR (Classless Inter-Domain Routing)

* Defines the IP range of your VPC or subnet
* Example: `10.0.0.0/16` → 65,536 IPs

**Rules**

* VPC CIDR cannot overlap with:

  * Another peered VPC
  * On-prem network (for VPN / Direct Connect)
* Subnet CIDR must be **subset of VPC CIDR**

---

### 2.2 Regions & Availability Zones

* **Region**: Geographical location (e.g., Mumbai, Virginia)
* **Availability Zone (AZ)**: Physically isolated data centers inside a region

**Key Rule**

* A subnet belongs to **exactly one AZ**
* High availability = **multiple AZs**

---

## 3. VPC Core Components

### 3.1 Subnets

A **subnet** is a segment of your VPC IP range.

| Type            | Internet Access | Typical Use                 |
| --------------- | --------------- | --------------------------- |
| Public Subnet   | Yes             | ALB, Bastion, NAT Gateway   |
| Private Subnet  | No (direct)     | App servers, DB             |
| Isolated Subnet | No              | Secure DB, internal systems |

**Public vs Private**

* Public subnet → route to Internet Gateway
* Private subnet → route to NAT Gateway or none

---

### 3.2 Route Tables

A **route table** decides **where network traffic goes**.

Example:

```
10.0.0.0/16 → local
0.0.0.0/0 → igw-xxxx
```

Rules:

* Every subnet must be associated with **one route table**
* If not specified → uses **main route table**

---

### 3.3 Internet Gateway (IGW)

* Enables **internet access** for resources with **public IP**
* Attached **at VPC level**

Flow:

```
EC2 → Route Table → IGW → Internet
```

---

### 3.4 NAT Gateway / NAT Instance

Allows **outbound-only internet access** from private subnets.

| Feature     | NAT Gateway | NAT Instance |
| ----------- | ----------- | ------------ |
| Managed     | Yes         | No           |
| HA          | Built-in    | Manual       |
| Performance | High        | Limited      |
| Cost        | Higher      | Lower        |

**Best Practice**: Use **NAT Gateway**

---

## 4. IP Addressing in VPC

### 4.1 Public IP

* Assigned at instance launch
* Changes on stop/start

### 4.2 Elastic IP (EIP)

* Static public IP
* Must be explicitly allocated
* Charged when unused

### 4.3 Private IP

* Primary + optional secondary IPs
* Never changes

---

## 5. VPC Security (CRITICAL)

![Image](https://media.geeksforgeeks.org/wp-content/uploads/20240514172628/aws-security-groups.webp)

![Image](https://docs.aws.amazon.com/images/vpc/latest/userguide/images/security-group-referencing.png)

![Image](https://docs.aws.amazon.com/images/vpc/latest/userguide/images/network-acl.png)

### 5.1 Security Groups (SG)

* **Stateful**
* Applied to **ENI / EC2**
* Allow rules only

Example:

```
Inbound: TCP 22 from MyIP
Outbound: All traffic
```

---

### 5.2 Network ACLs (NACL)

* **Stateless**
* Applied at **subnet level**
* Allow + Deny rules
* Ordered by rule number

**Comparison**

| Feature    | SG       | NACL   |
| ---------- | -------- | ------ |
| Level      | Instance | Subnet |
| Stateful   | Yes      | No     |
| Deny Rules | No       | Yes    |

---

## 6. Elastic Network Interfaces (ENI)

* Virtual network card
* Has:

  * MAC address
  * Private IP
  * Security Groups
* Can be attached/detached from instances

**Use Case**

* Failover
* Multiple IPs
* Multi-homed instances

---

## 7. DNS in VPC

### 7.1 AmazonProvidedDNS

* Enabled by default
* Resolves internal AWS DNS names

### 7.2 Route 53 Integration

* Private Hosted Zones
* Internal service discovery

---

## 8. VPC Connectivity Options (Advanced)

![Image](https://docs.aws.amazon.com/images/prescriptive-guidance/latest/integrate-third-party-services/images/p2_vpc-peering.png)

![Image](https://docs.aws.amazon.com/images/whitepapers/latest/building-scalable-secure-multi-vpc-network-infrastructure/images/hub-and-spoke-design.png)

![Image](https://docs.aws.amazon.com/images/whitepapers/latest/aws-vpc-connectivity-options/images/aws-managed-vpn.png)

![Image](https://docs.aws.amazon.com/images/whitepapers/latest/aws-vpc-connectivity-options/images/aws-direct-connect.png)

### 8.1 VPC Peering

* One-to-one VPC connection
* No transitive routing
* CIDR must not overlap

---

### 8.2 Transit Gateway (TGW)

* Hub-and-spoke model
* Centralized routing
* Supports:

  * VPCs
  * VPN
  * Direct Connect

**Enterprise standard**

---

### 8.3 Site-to-Site VPN

* IPsec tunnel
* Over public internet
* Encrypted

---

### 8.4 Client VPN

* Remote user access
* OpenVPN-based

---

### 8.5 Direct Connect

* Dedicated private connection
* High bandwidth
* Low latency

---

## 9. VPC Endpoints (Private AWS Access)

| Type               | Example       |
| ------------------ | ------------- |
| Gateway Endpoint   | S3, DynamoDB  |
| Interface Endpoint | EC2, SSM, SNS |

**Benefit**

* No internet/NAT required
* Traffic stays within AWS network

---

## 10. Flow Logs & Monitoring

### 10.1 VPC Flow Logs

Capture:

* Source IP
* Destination IP
* Port
* Action (ACCEPT / REJECT)

Stored in:

* CloudWatch Logs
* S3

Use Cases:

* Security audits
* Troubleshooting
* Compliance

---

## 11. High Availability VPC Design (Real World)

**Best Practice Architecture**

* Multi-AZ
* Public subnet:

  * ALB
  * NAT Gateway (one per AZ)
* Private subnet:

  * App servers
  * DB replicas
* No single point of failure

---

## 12. VPC Limits & Quotas

* VPCs per region
* Subnets per VPC
* Route tables
* IGWs

(Some limits are soft → increase via support)

---

## 13. Common Interview Questions

1. Difference between SG and NACL?
2. Why NAT Gateway is in public subnet?
3. Can private subnet access internet?
4. What happens if IGW is removed?
5. Can VPC CIDR be changed?
6. Difference between VPC Peering and Transit Gateway?
7. Why VPC endpoints are preferred over NAT?

---

## 14. Common Production Mistakes

* Overlapping CIDRs
* Single AZ NAT Gateway
* Public DB instances
* Over-permissive SG rules (`0.0.0.0/0`)
* No VPC Flow Logs

---

## 15. Mental Model (Easy Analogy)

| Real World     | VPC              |
| -------------- | ---------------- |
| City           | AWS Region       |
| Apartment      | VPC              |
| Floor          | Subnet           |
| Security Guard | Security Group   |
| Main Gate      | Internet Gateway |
| Exit-only Door | NAT Gateway      |

---

## 16. Certification-Focused Summary

* VPC = isolation + control
* SG = instance firewall
* NACL = subnet firewall
* NAT = outbound internet
* IGW = inbound/outbound internet
* TGW = enterprise networking
* Endpoints = private AWS access

---
