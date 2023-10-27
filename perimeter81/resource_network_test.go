package perimeter81

import (
	"context"
	"fmt"
	"testing"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var randNameNetwork string = randStringBytesRmndr()
var randNameNetworkUpdated string = randStringBytesRmndr()

func TestAccNetwork_basic(t *testing.T) {
	t.Parallel()
	var network perimeter81Sdk.Network

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkExists("perimeter81_network.n", &network),
					testAccCheckNetworkAttributes(&network, &testAccNetworkExpectedAttributes{
						Name:   randNameNetwork,
						Tags:   []string{"test"},
						Subnet: "10.30.0.90/16",
					}),
				),
			},
			{
				Config: testAccNetworkUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkExists("perimeter81_network.n", &network),
					testAccCheckNetworkAttributes(&network, &testAccNetworkExpectedAttributes{
						Name:   randNameNetworkUpdated,
						Tags:   []string{"test", "updated"},
						Subnet: "10.30.0.90/16",
					}),
				),
			},
		},
	})
}

func testAccCheckNetworkExists(n string, network *perimeter81Sdk.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		networkID := rs.Primary.ID
		if networkID == "" {
			return fmt.Errorf("No network ID is set")
		}
		conn := testAccProvider.Meta().(*perimeter81Sdk.APIClient)
		ctx := context.Background()
		gotNetwork, _, err := conn.NetworksApi.NetworksControllerV2NetworkFind(ctx, networkID)
		if err != nil {
			return err
		}
		*network = gotNetwork
		return nil
	}
}

type testAccNetworkExpectedAttributes struct {
	Name   string
	Tags   []string
	Subnet string
}

func testAccCheckNetworkAttributes(network *perimeter81Sdk.Network, want *testAccNetworkExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if network.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", network.Name, want.Name)
		}

		if !testComparableArraiesEq(network.Tags, want.Tags) {
			return fmt.Errorf("got tags %q; want %q", network.Tags, want.Tags)
		}

		return nil
	}
}

func testAccNetworkConfig() string {
	config := `
resource "perimeter81_network" "n" {
	network {
		name = "%s"
		tags = ["test"]
		subnet = "10.30.0.90/16"
	}
	region {
		cpregion_id = "r2Epw6OJsx"
		idle = true
	}
}
  `
	return fmt.Sprintf(config, randNameNetwork)
}

func testAccNetworkUpdateConfig() string {
	config := `
resource "perimeter81_network" "n" {
	network {
		name = "%s"
		tags = ["test", "updated"]
		subnet = "10.30.0.90/16"
	}
	region {
		cpregion_id = "r2Epw6OJsx"
		idle = true
	}
}
  `
	return fmt.Sprintf(config, randNameNetworkUpdated)
}
