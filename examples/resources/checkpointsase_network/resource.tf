data "checkpointsase_regions" "all" {}

# Provisions a standard network with one cloud gateway in the first available region.
resource "checkpointsase_network" "example" {
  network {
    name   = "tfExampleNet"
    subnet = "10.0.0.0/16"
    tags   = ["terraform", "example"]
  }

  region {
    cpregion_id = data.checkpointsase_regions.all.regions[0].id
    idle        = false
  }
}
