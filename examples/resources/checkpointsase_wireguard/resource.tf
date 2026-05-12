# A WireGuard tunnel attached to a gateway in a standard network's region.
# region_id, network_id, and gateway_id come from the network's deployment.
resource "checkpointsase_wireguard" "example" {
  network_id      = "ZwAeo5wqiF"
  region_id       = "K7tEfRm9vQ"
  gateway_id      = "abc12345DE"
  tunnel_name     = "wgTunnel01"
  remote_endpoint = "203.0.113.10"
  remote_subnets  = ["192.168.10.0/24"]
}
