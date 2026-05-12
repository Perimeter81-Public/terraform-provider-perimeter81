# A redundant (active/standby) IPsec tunnel pair for a standard network.
resource "checkpointsase_ipsec_redundant" "example" {
  network_id  = "ZwAeo5wqiF"
  region_id   = "K7tEfRm9vQ"
  tunnel_name = "ipsecRedundant01"

  shared_settings {
    passphrase             = "ChangeMe-shared-secret"
    remote_gateway_subnets = ["192.168.30.0/24"]
  }

  advanced_settings {
    key_exchange  = "ikev2"
    ike_life_time = 86400
    lifetime      = 3600
    dpd_delay     = 30
    dpd_timeout   = 120

    phase_one_proposals {
      auth       = "sha256"
      encryption = "aes-cbc-256"
      dh         = "modp2048"
    }

    phase_two_proposals {
      auth       = "sha256"
      encryption = "aes-cbc-256"
      dh         = "modp2048"
    }
  }
}
