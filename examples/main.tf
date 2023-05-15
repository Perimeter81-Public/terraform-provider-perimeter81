terraform {
  required_providers {
    perimeter81 = {
      version = "0.3.1"
      source  = "hashicorp.com/edu/perimeter81"
    }
  }
}

provider "perimeter81" {}

module "psl" {
  source = "./network"

  network_name = "Network 1"
}

# output "psl" {
#   value = module.psl.network_name
# }
