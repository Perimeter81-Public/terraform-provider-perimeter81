# Manage the gateway pool of a standard network's existing region.
# region_id is the ID of the network's region (not the cloud region ID).
resource "checkpointsase_gateway" "example" {
  network_id = "ZwAeo5wqiF"
  region_id  = "K7tEfRm9vQ"

  gateways {
    name = "gw1"
    idle = false
  }
}
