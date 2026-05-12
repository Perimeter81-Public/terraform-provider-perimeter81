data "checkpointsase_enhanced_regions" "all" {}

# Provisions an enhanced (SD-WAN-capable) network with one region.
resource "checkpointsase_enhanced_network" "example" {
  name   = "tfExampleEnh"
  subnet = "10.1.0.0/16"
  tags   = ["terraform", "example"]

  region {
    harmony_sase_region_id = data.checkpointsase_enhanced_regions.all.regions[0].id
    scale_units            = 1
    idle                   = false
  }
}
