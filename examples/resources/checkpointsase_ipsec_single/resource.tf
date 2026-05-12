# A single-gateway IPsec tunnel attached to a standard network.
resource "checkpointsase_ipsec_single" "example" {
  network_id             = "ZwAeo5wqiF"
  region_id              = "K7tEfRm9vQ"
  gateway_id             = "abc12345DE"
  tunnel_name            = "ipsecSingle01"
  remote_public_ip       = "203.0.113.20"
  remote_gateway_subnets = ["192.168.20.0/24"]
  passphrase             = "ChangeMe-shared-secret"

  key_exchange = "ikev2"
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
