# Lists customer-uploaded TLS certificates available to enhanced networks.
# Returns an empty list (no error) when the tenant has no enhanced networks.
data "checkpointsase_customer_certificates" "all" {}
