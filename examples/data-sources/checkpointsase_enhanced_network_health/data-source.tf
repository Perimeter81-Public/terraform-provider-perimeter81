data "checkpointsase_enhanced_networks" "all" {}

# Health check for a single enhanced network.
data "checkpointsase_enhanced_network_health" "first" {
  network_id = data.checkpointsase_enhanced_networks.all.networks[0].id
}
