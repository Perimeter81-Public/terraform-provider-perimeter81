# An active/standby IPsec redundant tunnel pair for a standard network.
# Two tunnels (tunnel1 + tunnel2) terminate at distinct remote endpoints for
# failover; shared and advanced settings apply to both. Timing fields use
# duration-string syntax ("30s", "3600s", "8h"). `dh` values are integer
# Diffie-Hellman group numbers (e.g. 14 = MODP2048).
resource "checkpointsase_ipsec_redundant" "example" {
  network_id  = "ZwAeo5wqiF"
  region_id   = "K7tEfRm9vQ"
  tunnel_name = "ipsecRedundant01"

  shared_settings {
    p81_gateway_subnets    = ["10.99.0.0/24"]
    remote_gateway_subnets = ["192.168.30.0/24"]
  }

  tunnel1 {
    gateway_id           = "abc12345DE"
    passphrase           = "ChangeMe-tunnel1-secret"
    p81_gwinternal_ip    = "169.254.0.1"
    remote_gwinternal_ip = "169.254.0.2"
    remote_public_ip     = "203.0.113.30"
    remote_asn           = "65010"
  }

  tunnel2 {
    gateway_id           = "abc12345DF"
    passphrase           = "ChangeMe-tunnel2-secret"
    p81_gwinternal_ip    = "169.254.1.1"
    remote_gwinternal_ip = "169.254.1.2"
    remote_public_ip     = "203.0.113.31"
    remote_asn           = "65010"
  }

  advanced_settings {
    key_exchange  = "ikev2"
    ike_life_time = "28800s"
    lifetime      = "3600s"
    dpd_delay     = "30s"
    dpd_timeout   = "60s"

    phase1 {
      auth       = ["sha256"]
      encryption = ["aes-cbc-256"]
      dh         = [14]
    }

    phase2 {
      auth       = ["sha256"]
      encryption = ["aes-cbc-256"]
      dh         = [14]
    }
  }
}
