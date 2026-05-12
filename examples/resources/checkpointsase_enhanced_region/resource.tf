data "checkpointsase_enhanced_networks" "all" {}
data "checkpointsase_enhanced_regions" "all" {}

# Add a secondary region to an existing enhanced network.
resource "checkpointsase_enhanced_region" "second_region" {
  network_id             = data.checkpointsase_enhanced_networks.all.networks[0].id
  harmony_sase_region_id = data.checkpointsase_enhanced_regions.all.regions[0].id
  scale_units            = 1
  idle                   = false
}
