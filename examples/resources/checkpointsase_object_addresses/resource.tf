resource "checkpointsase_object_addresses" "internal_api" {
  name        = "internalApi"
  description = "Internal API host"
  value_type  = "ip"
  value       = ["10.0.0.1"]
}
