# Different Modes of Vault

## Overview
1. **Dev Mode** - Development and testing
2. **Standalone Mode** - Single instance production
3. **HA Mode** - High availability production

---

## 1. Dev Mode

A single-node Vault that:
- Runs in memory
- Auto-unseals
- Auto-creates a root token
- Deletes all secrets on restart
- Good for local testing

### ✅ Advantages
- Very easy to start
- No unseal keys required
- No configuration needed

### ❌ Disadvantages
- Not persistent
- Not secure
- **NOT for production**
- Only 1 instance (no clustering)

---

## 2. Standalone Mode

A single Vault instance that:
- Stores data persistently
- Requires manual unsealing
- Suitable for small teams/projects

### ✅ Advantages
- Persistent storage
- Easy to configure
- Can be used for small teams/projects

### ❌ Disadvantages
- No HA (if Vault goes down → secrets unavailable)
- Single point of failure
- Not ideal for real production

---

## 3. HA Mode (Production Mode)

A cluster of 2 or more Vault nodes, usually using:
- **Raft Integrated Storage** (recommended)
- OR Consul (older method)

### ✅ Advantages
- Automatic failover
- Highly available
- Stable for production
- Scales better than standalone

### ❌ Disadvantages
- Slightly more complex setup
- Requires multiple servers
- Needs proper networking between nodes