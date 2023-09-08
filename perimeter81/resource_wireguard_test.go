package perimeter81

import (
	"context"
	"fmt"
	"testing"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var randNameWireguard string = randStringBytesRmndr()

func TestAccWireguard_basic(t *testing.T) {
	t.Parallel()
	var tunnel perimeter81Sdk.WireguardTunnel

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWireguardConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWireguardExists("perimeter81_wireguard.wgd1", &tunnel),
					testAccCheckWireguardAttributes(&tunnel, &testAccWireguardExpectedAttributes{
						RemoteEndpoint: "192.177.100.42",
						RemoteSubnets:  []string{"192.177.255.255/32"},
					}),
				),
			},
			{
				Config: testAccWireguardUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWireguardExists("perimeter81_wireguard.wgd1", &tunnel),
					testAccCheckWireguardAttributes(&tunnel, &testAccWireguardExpectedAttributes{
						RemoteEndpoint: "192.178.100.42",
						RemoteSubnets:  []string{"192.178.255.255/32"},
					}),
				),
			},
		},
	})
}

func testAccCheckWireguardExists(n string, tunnel *perimeter81Sdk.WireguardTunnel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		tunnelId := rs.Primary.ID
		if tunnelId == "" {
			return fmt.Errorf("No tunnel id is set")
		}
		conn := testAccProvider.Meta().(*perimeter81Sdk.APIClient)
		ctx := context.Background()
		networkId := rs.Primary.Attributes["network_id"]
		gotWireguard, _, err := conn.WireguardApi.GetWireguardTunnel(ctx, networkId, tunnelId)
		if err != nil {
			return err
		}
		*tunnel = gotWireguard
		return nil
	}
}

type testAccWireguardExpectedAttributes struct {
	RemoteEndpoint string
	RemoteSubnets  []string
}

func testAccCheckWireguardAttributes(tunnel *perimeter81Sdk.WireguardTunnel, want *testAccWireguardExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if tunnel.RemoteEndpoint != want.RemoteEndpoint {
			return fmt.Errorf("got remote endpoint %q; want %q", tunnel.RemoteEndpoint, want.RemoteEndpoint)
		}

		if !testComparableArraiesEq(tunnel.RemoteSubnets, want.RemoteSubnets) {
			return fmt.Errorf("got remote subnets %q; want %q", tunnel.RemoteSubnets, want.RemoteSubnets)
		}

		return nil
	}
}

func testAccWireguardConfig() string {
	config := `
resource "perimeter81_network" "n1" {
  network {
    name = "%s"
    tags = ["test"]
  }
  region {
    cpregion_id = "Xv3BREC4QI"
    idle = true
  }
}

data "perimeter81_networks" "all" {
	depends_on = [
    	perimeter81_network.n1
  	]
}

resource "perimeter81_wireguard" "wgd1" { 
  network_id = perimeter81_network.n1.id
  remote_endpoint = "192.177.100.42"
  region_id = {
    for network in data.perimeter81_networks.all.networks :
    network.id => network.regions[0].id
    if network.id == perimeter81_network.n1.id
  }[perimeter81_network.n1.id]
  gateway_id = {
    for network in data.perimeter81_networks.all.networks :
    network.id => network.regions[0].instances[0].id
    if network.id == perimeter81_network.n1.id
  }[perimeter81_network.n1.id]
  tunnel_name = "Wireguard1"
  remote_subnets = ["192.177.255.255/32"]
}
  `
	return fmt.Sprintf(config, randNameWireguard)
}

func testAccWireguardUpdateConfig() string {
	config := `
resource "perimeter81_network" "n1" {
  network {
    name = "%s"
    tags = ["test"]
  }
  region {
    cpregion_id = "Xv3BREC4QI"
    idle = true
  }
}

data "perimeter81_networks" "all1" {
	depends_on = [
    	perimeter81_network.n1
  	]
}

resource "perimeter81_wireguard" "wgd1" { 
  network_id = perimeter81_network.n1.id
  remote_endpoint = "192.178.100.42"
  region_id = {
    for network in data.perimeter81_networks.all1.networks :
    network.id => network.regions[0].id
    if network.id == perimeter81_network.n1.id
  }[perimeter81_network.n1.id]
  gateway_id = {
    for network in data.perimeter81_networks.all1.networks :
    network.id => network.regions[0].instances[0].id
    if network.id == perimeter81_network.n1.id
  }[perimeter81_network.n1.id]
  tunnel_name = "Wireguard1"
  remote_subnets = ["192.178.255.255/32"]
}
  `
	return fmt.Sprintf(config, randNameWireguard)
}
