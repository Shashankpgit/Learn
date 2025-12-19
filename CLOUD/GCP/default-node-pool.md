This is a **very common GKE + Terraform behavior**, and it’s not a bug.
I’ll explain **why it happens**, **what Terraform is doing**, and **how you can avoid it or keep the default node pool**.

---

## 1️⃣ Why you see “default node pool created → deleted”

### What GKE does by default

When you create a GKE cluster **without explicitly disabling it**, GKE **always**:

* Creates a **default node pool**
* Name: `default-pool`
* Size: usually **3 nodes**

This is **hard-coded GKE behavior**.

---

## 2️⃣ What Terraform does on top of that

In Terraform, when you write something like:

```hcl
resource "google_container_cluster" "gke" {
  name = "my-cluster"

  remove_default_node_pool = true
  initial_node_count       = 1
}
```

Terraform does this **sequence**:

1️⃣ Create cluster
→ GKE **must** create a default node pool (requirement)

2️⃣ Terraform immediately deletes `default-pool`
→ Because you said `remove_default_node_pool = true`

3️⃣ Terraform creates **your custom node pools**
→ With your machine type, autoscaling, taints, labels, etc.

This is expected and correct.

---

## 3️⃣ Why Terraform MUST do this

Because:

* GKE **does not allow a cluster with zero node pools**
* Terraform must give GKE:

  ```hcl
  initial_node_count = 1
  ```
* Then it removes the default pool once the cluster exists

There is **no API to create a cluster without a default pool**.

---

## 4️⃣ Can you keep the default node pool instead?

### ✅ YES — but with limitations

To keep the default node pool:

```hcl
resource "google_container_cluster" "gke" {
  name = "my-cluster"

  remove_default_node_pool = false
  initial_node_count       = 3
}
```

### What you lose if you keep it

❌ Cannot manage autoscaling cleanly
❌ Cannot customize machine types per pool
❌ Harder to manage upgrades
❌ Mixing infra & app workloads

That’s why **most production clusters remove it**.

---

## 5️⃣ Best practice (why people delete default pool)

Production GKE clusters usually have:

* `system-pool` (infra)
* `app-pool`
* `batch-pool`
* `gpu-pool`

Each with:

* Different machine types
* Different autoscaling rules
* Different taints

Default pool:

* Can’t be renamed
* Harder to control
* Becomes a liability

---

## 6️⃣ Can you “revert back” and use default pool?

### Yes, two ways:

### Option A: Keep it from the beginning (recommended if you want it)

```hcl
remove_default_node_pool = false
```

### Option B: Recreate the cluster

Once deleted, the default pool **cannot be restored**.

You must:

* Destroy cluster
* Recreate cluster
* Keep default pool

---

## 7️⃣ Simple timeline (visual)

```
Terraform apply
   |
   |--> GKE creates cluster + default-pool
   |
   |--> Terraform deletes default-pool
   |
   |--> Terraform creates custom node pools
```

This is normal.

---

## 8️⃣ One-line takeaway (memorize this)

> **GKE always creates a default node pool.
> Terraform deletes it if you ask, then creates custom pools.
> You cannot create a cluster without a default pool at creation time.**

---

## 9️⃣ Recommendation (honest advice)

* ✅ For learning → keep default pool
* ✅ For production → remove default pool
* ❌ Don’t mix both approaches

---
Short, clear answer first:

👉 **No — you cannot stop GKE from creating the default node pool at cluster creation time.**

That behavior is **mandatory in GKE**, not a Terraform choice.

Now let me explain **what you *can* control**, **what you cannot**, and **the correct patterns**.

---

## 1️⃣ What you CANNOT do (important)

❌ You **cannot** tell GKE:

> “Create the cluster without any node pool.”

GKE **always** creates a `default-pool` during cluster creation.

This is a **hard GKE API requirement**.

---

## 2️⃣ What Terraform is actually doing

When you write:

```hcl
remove_default_node_pool = true
initial_node_count       = 1
```

Terraform is doing this:

1️⃣ Ask GKE to create the cluster
→ GKE **must** create `default-pool`

2️⃣ After cluster is ready
→ Terraform deletes `default-pool`

So Terraform is not “creating” the default pool — **GKE is**.

---

## 3️⃣ What you CAN control

### ✅ You can remove it immediately

```hcl
remove_default_node_pool = true
```

This is the **only supported way**.

---

### ✅ You can keep it

```hcl
remove_default_node_pool = false
```

But you accept its limitations.

---

### ❌ You cannot disable its creation

There is **no Terraform flag** and **no GKE API** for this.

---

## 4️⃣ Why GKE enforces this (design reason)

* Cluster needs:

  * kubelet
  * system pods
  * CNI
* GKE requires at least one node pool at creation
* Default pool is the bootstrap mechanism

---

## 5️⃣ Best practice (production)

The standard pattern is:

```text
Create cluster
↓
Delete default-pool
↓
Create custom node pools
```

Yes, it looks weird
Yes, it’s intentional
Yes, everyone does this

---

## 6️⃣ One-line final answer (memorize)

> **You cannot prevent GKE from creating the default node pool.
> You can only delete it immediately using `remove_default_node_pool = true`.**

---


