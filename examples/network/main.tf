terraform {
  required_providers {
    perimeter81 = {
      version = "0.3.1"
      source  = "hashicorp.com/edu/perimeter81"
    }
  }
}

variable "network_name" {
  type    = string
  default = "Test network"
}

resource "perimeter81_network" "psl" {
  network {
    name = var.network_name
    tags = ["test", "testo"]
    subnet = "alo"
  }
  region {
    cpregionid = "fshulsBA14"
    instancecount = 1
    idle = true
  }
}
# data "perimeter81_networks" "all" {}

# Returns all networks
# output "all_networks" {
#   value = data.perimeter81_networks.all.networks
# }

# output "network" {
#   value = {
#     for network in data.perimeter81_networks.all.networks :
#     network.id => network
#     if network.name == var.network_name
#   }
# }
