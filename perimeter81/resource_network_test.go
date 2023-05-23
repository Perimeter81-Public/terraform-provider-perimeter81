package perimeter81

import (
	"context"
	"fmt"
	"testing"

	perimeter81Sdk "terraform-provider-perimeter81/perimeter81sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetwork_basic(t *testing.T) {
	var network perimeter81Sdk.Network

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkExists("perimeter81_network.n1", &network),
					testAccCheckNetworkAttributes(&network, &testAccNetworkExpectedAttributes{
						Name: "Network test",
						Tags: []string{"test"},
					}),
				),
			},
			{
				Config: testAccNetworkUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkExists("perimeter81_network.n1", &network),
					testAccCheckNetworkAttributes(&network, &testAccNetworkExpectedAttributes{
						Name: "Network test updated",
						Tags: []string{"test", "updated"},
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
	Name string
	Tags []string
}

func testTagsEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func testAccCheckNetworkAttributes(network *perimeter81Sdk.Network, want *testAccNetworkExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if network.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", network.Name, want.Name)
		}

		if !testTagsEq(network.Tags, want.Tags) {
			return fmt.Errorf("got tags %q; want %q", network.Tags, want.Tags)
		}

		return nil
	}
}

func testAccNetworkConfig() string {
	return `
resource "perimeter81_network" "n1" {
	network {
		name = "Network test"
		tags = ["test"]
	}
	region {
		cpregionid = "r2Epw6OJsx"
		instancecount = 1
		idle = true
	}
}
  `
}

func testAccNetworkUpdateConfig() string {
	return `
resource "perimeter81_network" "n1" {
	network {
		name = "Network test updated"
		tags = ["test", "updated"]
	}
	region {
		cpregionid = "r2Epw6OJsx"
		instancecount = 1
		idle = true
	}
}
  `
}
