data "checkpointsase_standard_networks" "all" {}

# Route table for a single standard network.
data "checkpointsase_route_table" "first" {
  network_id = data.checkpointsase_standard_networks.all.networks[0].id
}
