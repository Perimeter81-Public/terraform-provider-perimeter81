# A dynamic (BGP-routed) IPsec tunnel attached to an enhanced network.
# Timing fields use duration-string syntax ("30s", "60m", "8h"), not raw integers.
# Encryption values must match the public-api enum (PhaseEncryptionV2_1) — use
# "aes256", "aes128", etc. Passphrases must satisfy the IsPassphrase regex
# (letters/digits/`.`/`_`, 8-64 chars — no hyphens).
resource "checkpointsase_enhanced_dynamic_tunnel" "example" {
  network_id             = "ZwAeo5wqiF"
  tunnel_name            = "dynamicTunnel01"
  description            = "Dynamic tunnel — managed by Terraform"
  p81_gateway_subnets    = ["10.99.0.0/24"]
  remote_gateway_subnets = ["192.168.50.0/24"]
  peak_bandwidth         = 1000
  # BGP autonomous-system number for the Check Point SASE side (required).
  left_asn = 65000

  tunnel {
    region_id            = "K7tEfRm9vQ"
    auth_type            = "psk"
    passphrase           = "ChangeMeSharedSecret"
    remote_public_ip     = "203.0.113.50"
    remote_asn           = 65010
    p81_gw_internal_ip   = "169.254.0.1"
    remote_gw_internal_ip = "169.254.0.2"
  }

  key_exchange  = "ikev2"
  ike_life_time = "28800s"
  lifetime      = "3600s"
  dpd_delay     = "30s"
  dpd_timeout   = "60s"

  phase1 {
    auth                = ["sha256"]
    encryption          = ["aes256"]
    key_exchange_method = ["modp2048"]
  }

  phase2 {
    auth                = ["sha256"]
    encryption          = ["aes256"]
    key_exchange_method = ["modp2048"]
  }
}
