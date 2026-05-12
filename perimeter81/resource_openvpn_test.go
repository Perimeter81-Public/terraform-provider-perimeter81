package perimeter81

import (
	"context"
	"fmt"
	"testing"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var randNameOpenVpn string = randStringBytesRmndr()

func TestAccOpenvpn_basic(t *testing.T) {
	t.Parallel()
	var tunnel perimeter81Sdk.OpenVPNTunnel
	var access_key string
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenvpnConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOpenvpnExists("sase_openvpn.ovpn2", &tunnel, &access_key),
				),
			},
			{
				Config: testAccOpenvpnUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOpenvpnExists("sase_openvpn.ovpn2", &tunnel, &access_key),
					testAccCheckOpenvpnAttributes(&tunnel, access_key),
				),
			},
		},
	})
}

func testAccCheckOpenvpnExists(n string, tunnel *perimeter81Sdk.OpenVPNTunnel, access_key *string) resource.TestCheckFunc {
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
		version := rs.Primary.Attributes["version"]
		gotOpenvpn, _, err := conn.OpenVPNAPI.StandardGetOpenVPNTunnel(ctx, networkId, tunnelId).Execute()
		if err != nil {
			return err
		}
		if version == "1" {
			*access_key = gotOpenvpn.GetSecretAccessKey()
		}
		*tunnel = *gotOpenvpn
		return nil
	}
}

func testAccCheckOpenvpnAttributes(tunnel *perimeter81Sdk.OpenVPNTunnel, access_key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if tunnel.GetSecretAccessKey() == access_key {
			return fmt.Errorf("Access key was %q; and didn't change", access_key)
		}

		return nil
	}
}

func testAccOpenvpnConfig() string {
	config := `
resource "sase_network" "n2" {
  network {
    name = "%s"
    tags = ["test"]
  }
  region {
    cpregion_id = "Xv3BREC4QI"
    idle = true
  }
}

data "sase_networks" "all2" {
	depends_on = [
    	sase_network.n2
  	]
}

resource "sase_openvpn" "ovpn2" {
  network_id = sase_network.n2.id
  region_id = {
    for network in data.sase_networks.all2.networks :
    network.id => network.regions[0].id
    if network.id == sase_network.n2.id
  }[sase_network.n2.id]
  gateway_id = {
    for network in data.sase_networks.all2.networks :
    network.id => network.regions[0].instances[0].id
    if network.id == sase_network.n2.id
  }[sase_network.n2.id]
  tunnel_name = "OpenVPNTunnel"
  version = 1
}
  `
	return fmt.Sprintf(config, randNameOpenVpn)
}

func testAccOpenvpnUpdateConfig() string {
	config := `
resource "sase_network" "n2" {
  network {
    name = "%s"
    tags = ["test"]
  }
  region {
    cpregion_id = "Xv3BREC4QI"
    idle = true
  }
}

data "sase_networks" "all2" {
	depends_on = [
    	sase_network.n2
  	]
}

resource "sase_openvpn" "ovpn2" {
  network_id = sase_network.n2.id
  region_id = {
    for network in data.sase_networks.all2.networks :
    network.id => network.regions[0].id
    if network.id == sase_network.n2.id
  }[sase_network.n2.id]
  gateway_id = {
    for network in data.sase_networks.all2.networks :
    network.id => network.regions[0].instances[0].id
    if network.id == sase_network.n2.id
  }[sase_network.n2.id]
  tunnel_name = "OpenVPNTunnel"
  version = 2
}
  `
	return fmt.Sprintf(config, randNameOpenVpn)
}
