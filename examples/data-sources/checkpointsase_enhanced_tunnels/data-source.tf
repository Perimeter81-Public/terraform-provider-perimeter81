data "checkpointsase_enhanced_networks" "all" {}

# Tunnels attached to a single enhanced network.
data "checkpointsase_enhanced_tunnels" "first" {
  network_id = data.checkpointsase_enhanced_networks.all.networks[0].id
}
