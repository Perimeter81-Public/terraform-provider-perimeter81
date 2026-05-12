data "checkpointsase_regions" "all" {}

resource "checkpointsase_network" "example" {
  network {
    name   = "tfExampleNet"
    subnet = "10.0.0.0/16"
  }
  region {
    cpregion_id = data.checkpointsase_regions.all.regions[0].id
  }
}

# Adopt-style: a firewall policy is auto-created with the network.
# Importing it under terraform lets you manage its enabled / allowed flags.
resource "checkpointsase_firewall_policy" "example" {
  network_id = checkpointsase_network.example.id
  enabled    = true
  allowed    = true
}
