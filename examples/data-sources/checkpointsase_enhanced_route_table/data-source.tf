data "checkpointsase_enhanced_networks" "all" {}

# Route table for a single enhanced network.
data "checkpointsase_enhanced_route_table" "first" {
  network_id = data.checkpointsase_enhanced_networks.all.networks[0].id
}
