# A static IPsec tunnel for an enhanced network's region.
resource "checkpointsase_enhanced_static_tunnel" "example" {
  network_id             = "ZwAeo5wqiF"
  region_id              = "K7tEfRm9vQ"
  tunnel_name            = "staticTunnel01"
  remote_public_ip       = "203.0.113.40"
  remote_gateway_subnets = ["192.168.40.0/24"]
  passphrase             = "ChangeMe-shared-secret"
  auth_type              = "passphrase"

  key_exchange  = "ikev2"
  ike_life_time = 86400
  lifetime      = 3600
  dpd_delay     = 30
  dpd_timeout   = 120

  proposals {
    auth       = "sha256"
    encryption = "aes-cbc-256"
    dh         = "modp2048"
  }
}
