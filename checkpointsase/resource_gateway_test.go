package checkpointsase

import (
	"fmt"
	"testing"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var randNameGateway string = randStringBytesRmndr()

func TestAccGateway_basic(t *testing.T) {
	t.Parallel()
	var network perimeter81Sdk.Network

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGatewaysConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkExists("checkpointsase_network.n6", &network),
					testAccCheckGatewaysCount(&network, 2),
				),
			},
			{
				Config: testAccGatewaysUpdate1Config(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkExists("checkpointsase_network.n6", &network),
					testAccCheckGatewaysCount(&network, 1),
				),
			},
			{
				Config: testAccGatewaysUpdate2Config(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkExists("checkpointsase_network.n6", &network),
					testAccCheckGatewaysCount(&network, 2),
				),
			},
		},
	})
}

func testAccCheckGatewaysCount(network *perimeter81Sdk.Network, want int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(network.Regions[0].Instances) != want {
			return fmt.Errorf("got gateway count %d; want %d", len(network.Regions[0].Instances), want)
		}

		return nil
	}
}

func testAccGatewaysConfig() string {
	config := `
resource "checkpointsase_network" "n6" {
	network {
		name = "%s"
		tags = ["test"]
	}
	region {
		cpregion_id = "r2Epw6OJsx"
		idle = true
	}
}

data "checkpointsase_networks" "all5" {
	depends_on = [
    	checkpointsase_network.n6
  	]
}

resource "checkpointsase_gateway"  "g1"{

  gateways {
      name = "%s"
      idle = true
  }

  network_id = checkpointsase_network.n6.id
  region_id = {
    for network in data.checkpointsase_networks.all5.networks :
    network.id => network.regions[0].id
    if network.id == checkpointsase_network.n6.id
  }[checkpointsase_network.n6.id]
}
  `
	return fmt.Sprintf(config, randNameGateway, randNameGateway)
}

func testAccGatewaysUpdate1Config() string {
	config := `
resource "checkpointsase_network" "n6" {
	network {
		name = "%s"
		tags = ["test"]
	}
	region {
		cpregion_id = "r2Epw6OJsx"
		idle = true
	}
}

data "checkpointsase_networks" "all5" {
	depends_on = [
    	checkpointsase_network.n6
  	]
}

resource "checkpointsase_gateway"  "g1"{

  network_id = checkpointsase_network.n6.id
  region_id = {
    for network in data.checkpointsase_networks.all5.networks :
    network.id => network.regions[0].id
    if network.id == checkpointsase_network.n6.id
  }[checkpointsase_network.n6.id]
}
  `
	return fmt.Sprintf(config, randNameGateway)
}
func testAccGatewaysUpdate2Config() string {
	config := `
resource "checkpointsase_network" "n6" {
	network {
		name = "%s"
		tags = ["test"]
	}
	region {
		cpregion_id = "r2Epw6OJsx"
		idle = true
	}
}

data "checkpointsase_networks" "all5" {
	depends_on = [
    	checkpointsase_network.n6
  	]
}

resource "checkpointsase_gateway"  "g1"{

 gateways {
      name = "%s"
      idle = true
  }

  network_id = checkpointsase_network.n6.id
  region_id = {
    for network in data.checkpointsase_networks.all5.networks :
    network.id => network.regions[0].id
    if network.id == checkpointsase_network.n6.id
  }[checkpointsase_network.n6.id]
}
  `
	return fmt.Sprintf(config, randNameGateway, randNameGateway)
}
