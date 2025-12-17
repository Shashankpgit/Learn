From the Terraform script, we are defining which Google APIs the node (VM) is allowed to request access tokens for.
Terraform is NOT granting actual permissions — it is only setting the maximum possible access.
<!-- `oauth_scopes = [
        "https://www.googleapis.com/auth/cloud-platform",
        "https://www.googleapis.com/auth/devstorage.read_only",
        "https://www.googleapis.com/auth/logging.write",
        "https://www.googleapis.com/auth/monitoring",
        "https://www.googleapis.com/auth/servicecontrol",
        "https://www.googleapis.com/auth/service.management.readonly",
        "https://www.googleapis.com/auth/trace.append",
      ]` -->

✅ 1. Your understanding of cloud-platform vs specific scopes
Your statement

if we keep https://www.googleapis.com/auth/cloud-platform it means it is allowed to request every api
if we keep specific scopes, node is allowed to request only those APIs

✅ This is 100% correct

Let me restate it in technical terms:

cloud-platform
https://www.googleapis.com/auth/cloud-platform


Means:

This VM/node can request OAuth tokens for any Google Cloud API

It does not automatically grant access — but it removes the API restriction layer.

So with cloud-platform:

API access is limited only by IAM

OAuth scope is no longer restrictive

Specific scopes (your list)
oauth_scopes = [
  "https://www.googleapis.com/auth/devstorage.read_only",
  "https://www.googleapis.com/auth/logging.write",
  "https://www.googleapis.com/auth/monitoring",
  "https://www.googleapis.com/auth/servicecontrol",
  "https://www.googleapis.com/auth/service.management.readonly",
  "https://www.googleapis.com/auth/trace.append",
]


Means:

The node can request tokens ONLY for these APIs

Even if IAM allows more, the node:

❌ Cannot call Pub/Sub

❌ Cannot call Secret Manager

❌ Cannot call Compute Engine API

❌ Cannot call Admin SDK, Apps Script, etc.

So yes — your mental model here is spot on.


Why Google shows these “categories”

Google’s OAuth scope documentation is structured like this:

API PRODUCT
 ├─ Scope 1
 ├─ Scope 2
 ├─ Scope 3


So you see:

Example
Apps Script API
  - https://www.googleapis.com/auth/script.projects
  - https://www.googleapis.com/auth/script.deployments


This does not mean:

Your node has access to Apps Script ❌

Your node can manage scripts ❌

It only means:

“If someone wants to use the Apps Script API, these are the scopes they would need.”