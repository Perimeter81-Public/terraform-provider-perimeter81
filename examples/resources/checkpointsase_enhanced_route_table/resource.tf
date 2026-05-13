# A static route attached to a single static-tunnel entry.
# `propagated` is a computed attribute and must not be set in configuration.
resource "checkpointsase_enhanced_route_table" "example" {
  network_id = "ZwAeo5wqiF"
  type       = "static"
  tunnel_id  = "tun-abc12345"
  subnets    = ["10.50.0.0/16"]
}
