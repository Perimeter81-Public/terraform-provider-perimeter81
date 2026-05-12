data "checkpointsase_standard_networks" "all" {}

# Health check for a single standard network.
data "checkpointsase_network_health" "first" {
  network_id = data.checkpointsase_standard_networks.all.networks[0].id
}
