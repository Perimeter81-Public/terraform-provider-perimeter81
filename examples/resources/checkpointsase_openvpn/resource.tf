# An OpenVPN tunnel attached to a gateway in a standard network's region.
# region_id, network_id, and gateway_id come from the network's deployment.
resource "checkpointsase_openvpn" "example" {
  network_id  = "ZwAeo5wqiF"
  region_id   = "K7tEfRm9vQ"
  gateway_id  = "abc12345DE"
  tunnel_name = "ovpnTunnel01"
  version     = 1
}
