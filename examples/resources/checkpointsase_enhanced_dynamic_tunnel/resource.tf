# A dynamic (BGP-routed) IPsec tunnel for an enhanced network.
resource "checkpointsase_enhanced_dynamic_tunnel" "example" {
  network_id             = "ZwAeo5wqiF"
  tunnel_name            = "dynamicTunnel01"
  description            = "Dynamic tunnel — managed by Terraform"
  remote_gateway_subnets = ["192.168.50.0/24"]

  tunnel {
    region_id        = "K7tEfRm9vQ"
    auth_type        = "passphrase"
    passphrase       = "ChangeMe-shared-secret"
    remote_public_ip = "203.0.113.50"
  }

  key_exchange  = "ikev2"
  ike_life_time = 86400
  lifetime      = 3600
  dpd_delay     = 30
  dpd_timeout   = 120
}
