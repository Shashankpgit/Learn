1. Enable GKE logging and monitoring for security and operational visibility. Currently both were set to none
2. Removed wide/generic roles and provided only specific roles to the nodes or vm.
3. Enabled private nodes. use `only RFC 1918 private addresses`
4. Disabled Public endpoint. So that the master Plane can only be accessed from the internal IP address
5. Update the `master_ipv4_cidr_block` with the pre-defined CIDR range.
6. We might need to create the KMS to keep the kubernetes secrets.
Note: if there is no KMS (or KMS=null or no secrete encryption kms key then throw the warning)
7. Updated the networking to allow the request from from the certain ranges instead of open to internet (0.0.0.0/0)
8. With respect to service account, we havent given the storage admin role and need to remove the configuration which creates/generates the service account key.

-------------------------------------------------------------------------------------------

-------------------------------------------------------------------------------------------

1. Keep the logging and monitoring none for now
2. check which roles are required and which are not and test them out.
3. enable private nodes
4. don't disable the Public endpoint for the master plane.
5. since we are keeping the public endpoint, read little bit more about the Point 5
`(master_ipv4_cidr_block` with the pre-defined CIDR range)`
6. KMS hold on for now.