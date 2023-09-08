package perimeter81

import (
	"fmt"
	"testing"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var randNameRegion string = randStringBytesRmndr()

func TestAccRegion_basic(t *testing.T) {
	t.Parallel()
	var network perimeter81Sdk.Network

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkExists("perimeter81_network.n5", &network),
					testAccCheckRegionsCount(&network, 1),
				),
			},
			{
				Config: testAccRegionsUpdate1Config(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkExists("perimeter81_network.n5", &network),
					testAccCheckRegionsCount(&network, 2),
				),
			},
			{
				Config: testAccRegionsUpdate2Config(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkExists("perimeter81_network.n5", &network),
					testAccCheckRegionsCount(&network, 1),
				),
			},
		},
	})
}

func testAccCheckRegionsCount(network *perimeter81Sdk.Network, want int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(network.Regions) != want {
			return fmt.Errorf("got region count %q; want %q", len(network.Regions), want)
		}

		return nil
	}
}

func testAccRegionsConfig() string {
	config := `
resource "perimeter81_network" "n5" {
	network {
		name = "%s"
		tags = ["test"]
	}
	region {
		cpregion_id = "r2Epw6OJsx"
		idle = true
	}
}
  `
	return fmt.Sprintf(config, randNameNetwork)
}

func testAccRegionsUpdate1Config() string {
	config := `
resource "perimeter81_network" "n5" {
	network {
		name = "%s"
		tags = ["test"]
	}
	region {
		cpregion_id = "r2Epw6OJsx"
		idle = true
	}
	region {
    	cpregion_id = "F2w4QTggWt"
    	idle = true
  	}
}
  `
	return fmt.Sprintf(config, randNameRegion)
}
func testAccRegionsUpdate2Config() string {
	config := `
resource "perimeter81_network" "n5" {
	network {
		name = "%s"
		tags = ["test"]
	}
	region {
		cpregion_id = "r2Epw6OJsx"
		idle = true
	}
}
  `
	return fmt.Sprintf(config, randNameRegion)
}
