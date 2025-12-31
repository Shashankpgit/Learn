![Image](https://digitalcloud.training/wp-content/uploads/2022/01/IAM-1.jpg)

![Image](https://miro.medium.com/v2/resize%3Afit%3A1200/1%2ACgfJAcRlay0O9amULak-fw.png)

![Image](https://miro.medium.com/1%2AJDZ5jC3Z4w2C4BJBcUI60w.gif)

![Image](https://docs.aws.amazon.com/images/prescriptive-guidance/latest/certificate-based-access-controls/images/iam-roles-anywhere-same-account-region.png)

Below are **all four deliverables**, written so you can **print, revise, practice, and pass interviews/certifications**. This is a **complete IAM mastery pack**.

All content applies to **Amazon Web Services**.

---

# 1️⃣ PRINTABLE AWS IAM CHEAT SHEET (1–2 Pages)

## IAM Core Components

| Component           | Purpose               | Key Rule                     |
| ------------------- | --------------------- | ---------------------------- |
| Root User           | Account owner         | Never use daily, enable MFA  |
| User                | Permanent identity    | For humans only              |
| Group               | Users collection      | Permissions assigned here    |
| Policy              | JSON permission rules | Allow / Deny actions         |
| Role                | Temporary identity    | For services & cross-account |
| STS                 | Temp credentials      | Used by roles                |
| SCP                 | Org-level guardrail   | Max allowed permissions      |
| Permission Boundary | Identity guardrail    | Limits max permissions       |

---

## IAM Hierarchy (Must Remember)

```
AWS Account
└── IAM
    ├── Policies
    ├── Groups → Users
    └── Roles → Services / External users
```

---

## Policy Evaluation Logic (EXAM GOLD)

1. Explicit **DENY**
2. Explicit **ALLOW**
3. Implicit deny (default)

---

## When to Use What (One Look)

| Use Case                     | IAM Feature         |
| ---------------------------- | ------------------- |
| Human login                  | User + Group        |
| EC2/Lambda access            | Role                |
| CI/CD                        | Role + OIDC         |
| Cross-account                | Role                |
| Temporary access             | Role                |
| Org-wide block               | SCP                 |
| Prevent privilege escalation | Permission Boundary |

---

## One-Line Memory Hook

> **Policy = rules**
> **User = who you are**
> **Group = where you belong**
> **Role = what you become temporarily**

---

# 2️⃣ INTERVIEW-ONLY AWS IAM Q&A (HIGH FREQUENCY)

### Q1. Is IAM regional or global?

**Answer:** IAM is a **global service**.

---

### Q2. Difference between IAM User and Role?

**Answer:**

* User → permanent identity
* Role → temporary identity (no credentials)

---

### Q3. Why are roles more secure than access keys?

**Answer:**

* Temporary credentials
* Auto-rotation
* No hard-coded secrets

---

### Q4. Can a role be added to a group?

**Answer:** No. Groups are **only for users**.

---

### Q5. What wins: Allow or Deny?

**Answer:** Explicit **DENY always wins**.

---

### Q6. Do SCPs grant permissions?

**Answer:** No. SCPs only **restrict** permissions.

---

### Q7. How does EC2 access S3 securely?

**Answer:** Using an **IAM Role**, not access keys.

---

### Q8. Difference: Managed vs Inline Policy?

**Answer:**

* Managed → reusable
* Inline → one identity only

---

### Q9. What is a Permission Boundary?

**Answer:** A **maximum permission limit** for a user or role.

---

### Q10. How do you audit IAM activity?

**Answer:** Using **CloudTrail** and **Access Analyzer**.

---

# 3️⃣ HANDS-ON AWS IAM LAB CHECKLIST (REAL PRACTICE)

## Beginner Labs

* [ ] Create IAM user
* [ ] Enable MFA for user
* [ ] Create group (Developers)
* [ ] Attach managed policy
* [ ] Login as IAM user

---

## Intermediate Labs

* [ ] Create custom IAM policy
* [ ] Restrict S3 bucket access
* [ ] Attach policy to group
* [ ] Test allowed vs denied actions

---

## Advanced Labs

* [ ] Create IAM role for EC2
* [ ] Attach role to EC2 instance
* [ ] Access S3 without keys
* [ ] Create Lambda execution role
* [ ] Restrict DynamoDB table access

---

## Enterprise-Level Labs

* [ ] Create cross-account role
* [ ] Configure trust policy
* [ ] Apply permission boundary
* [ ] Create SCP in AWS Organizations
* [ ] Analyze permissions using Access Analyzer

---

## Security Labs (Highly Recommended)

* [ ] Rotate access keys
* [ ] Generate credential report
* [ ] Detect overly permissive policies
* [ ] Enforce MFA via IAM condition

---

# 4️⃣ CERTIFICATION-FOCUSED AWS IAM NOTES (EXAM READY)

## Concepts AWS LOVES to Test

### 🔹 Least Privilege

* Never use `*` unless required
* Scope actions + resources tightly

---

### 🔹 Roles Over Users

* EC2, Lambda, ECS → ALWAYS roles
* CI/CD → Role + OIDC

---

### 🔹 Conditions

Know these:

* `aws:MultiFactorAuthPresent`
* `aws:SourceIp`
* `aws:RequestedRegion`
* Time-based access

---

### 🔹 Cross-Account Access

* Trust policy + permissions policy
* No access keys shared

---

### 🔹 Organizations & SCP

* SCPs override IAM permissions
* SCPs do not grant access
* Applied at OU or account level

---

### 🔹 IAM is Always Global

Even when resources are regional.

---

## Certification Trap Questions (Be Careful)

❌ “Attach policy directly to EC2”
✔ Use IAM Role

❌ “SCP grants permissions”
✔ SCP only limits

❌ “Use root user for admin work”
✔ Create admin IAM role/user

---

# FINAL 30-SECOND REVISION

* IAM = security foundation
* Policies define rules
* Users = humans
* Groups = users only
* Roles = services & temporary access
* Explicit deny always wins
* IAM is global

---
