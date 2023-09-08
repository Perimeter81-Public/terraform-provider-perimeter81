package perimeter81

import (
	"context"
	"fmt"
	"testing"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIpsecRedundant_basic(t *testing.T) {
	t.Parallel()
	var tunnel perimeter81Sdk.IpSecRedundantTunnels
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIpsecRedundantConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIpsecRedundantExists("perimeter81_ipsec_redundant.ipsr1", &tunnel),
					testAccCheckIpsecRedundantAttributes(&tunnel, &perimeter81Sdk.IpSecRedundantTunnels{
						TunnelName: "ipseed",
						SharedSettings: &perimeter81Sdk.IpSecSharedSettings{
							P81GatewaySubnets:    []string{"0.0.0.0/0"},
							RemoteGatewaySubnets: []string{"0.0.0.0/0"},
						},
						AdvancedSettings: &perimeter81Sdk.IpSecAdvancedSettings{
							KeyExchange: "ikev2",
							IkeLifeTime: "8h",
							Lifetime:    "1h",
							DpdDelay:    "10s",
							DpdTimeout:  "30s",
							Phase1: &perimeter81Sdk.IpSecPhase{
								Auth:       []string{"3des"},
								Encryption: []string{"sha256"},
								Dh:         []int32{14},
							},
							Phase2: &perimeter81Sdk.IpSecPhase{
								Auth:       []string{"3des"},
								Encryption: []string{"sha256"},
								Dh:         []int32{14},
							},
						},
					}),
				),
			},
		},
	})
}

func testAccCheckIpsecRedundantExists(n string, tunnel *perimeter81Sdk.IpSecRedundantTunnels) resource.TestCheckFunc {
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
		gotIpsecRedundant, _, err := conn.IPSecRedundantApi.GetIPSecRedundantTunnel(ctx, networkId, tunnelId)
		if err != nil {
			return err
		}

		*tunnel = gotIpsecRedundant
		return nil
	}
}

func testAccCheckIpsecRedundantAttributes(tunnel *perimeter81Sdk.IpSecRedundantTunnels, want *perimeter81Sdk.IpSecRedundantTunnels) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !testComparableArraiesEq(tunnel.SharedSettings.P81GatewaySubnets, want.SharedSettings.P81GatewaySubnets) {
			return fmt.Errorf("got p81 gateway subnets %q; want %q", tunnel.SharedSettings.P81GatewaySubnets, want.SharedSettings.P81GatewaySubnets)
		}
		if !testComparableArraiesEq(tunnel.SharedSettings.RemoteGatewaySubnets, want.SharedSettings.RemoteGatewaySubnets) {
			return fmt.Errorf("got remote gateway subnets %q; want %q", tunnel.SharedSettings.RemoteGatewaySubnets, want.SharedSettings.RemoteGatewaySubnets)
		}
		if tunnel.TunnelName != want.TunnelName {
			return fmt.Errorf("got tunnel name %q; want %q", tunnel.TunnelName, want.TunnelName)
		}
		if tunnel.AdvancedSettings.IkeLifeTime != want.AdvancedSettings.IkeLifeTime {
			return fmt.Errorf("got ike life time %q; want %q", tunnel.AdvancedSettings.IkeLifeTime, want.AdvancedSettings.IkeLifeTime)
		}
		if tunnel.AdvancedSettings.DpdDelay != want.AdvancedSettings.DpdDelay {
			return fmt.Errorf("got dpd delay %q; want %q", tunnel.AdvancedSettings.DpdDelay, want.AdvancedSettings.DpdDelay)
		}
		if tunnel.AdvancedSettings.DpdTimeout != want.AdvancedSettings.DpdTimeout {
			return fmt.Errorf("got dpd timeout %q; want %q", tunnel.AdvancedSettings.DpdTimeout, want.AdvancedSettings.DpdTimeout)
		}
		if tunnel.Tunnel1.GatewayID == "" {
			return fmt.Errorf("got Gateway id empty for tunnel 1")
		}
		if tunnel.Tunnel2.GatewayID == "" {
			return fmt.Errorf("got Gateway id empty for tunnel 2")
		}
		if tunnel.Tunnel1.TunnelID == "" {
			return fmt.Errorf("got Tunnel id empty for tunnel 1")
		}
		if tunnel.Tunnel2.GatewayID == "" {
			return fmt.Errorf("got Tunnel id empty for tunnel 2")
		}
		if !testComparableArraiesEq(tunnel.AdvancedSettings.Phase1.Auth, want.AdvancedSettings.Phase1.Auth) {
			return fmt.Errorf("got phase1 auth %q; want %q", tunnel.AdvancedSettings.Phase1.Auth, want.AdvancedSettings.Phase1.Auth)
		}
		if !testComparableArraiesEq(tunnel.AdvancedSettings.Phase1.Encryption, want.AdvancedSettings.Phase1.Encryption) {
			return fmt.Errorf("got phase1 encryption %q; want %q", tunnel.AdvancedSettings.Phase1.Encryption, want.AdvancedSettings.Phase1.Encryption)
		}
		if !testComparableArraiesEq(tunnel.AdvancedSettings.Phase1.Dh, want.AdvancedSettings.Phase1.Dh) {
			return fmt.Errorf("got phase1 encryption %q; want %q", tunnel.AdvancedSettings.Phase1.Dh, want.AdvancedSettings.Phase1.Dh)
		}
		if !testComparableArraiesEq(tunnel.AdvancedSettings.Phase2.Auth, want.AdvancedSettings.Phase2.Auth) {
			return fmt.Errorf("got phase2 auth %q; want %q", tunnel.AdvancedSettings.Phase2.Auth, want.AdvancedSettings.Phase2.Auth)
		}
		if !testComparableArraiesEq(tunnel.AdvancedSettings.Phase2.Encryption, want.AdvancedSettings.Phase2.Encryption) {
			return fmt.Errorf("got phase2 encryption %q; want %q", tunnel.AdvancedSettings.Phase2.Encryption, want.AdvancedSettings.Phase2.Encryption)
		}
		if !testComparableArraiesEq(tunnel.AdvancedSettings.Phase2.Dh, want.AdvancedSettings.Phase2.Dh) {
			return fmt.Errorf("got phase2 dh %q; want %q", tunnel.AdvancedSettings.Phase2.Dh, want.AdvancedSettings.Phase2.Dh)
		}

		return nil
	}
}

func testAccIpsecRedundantConfig() string {
	config := `

resource "perimeter81_network" "n4" {
  network {
    name = "%s"
    tags = ["test"]
  }
  region {
    cpregion_id = "Xv3BREC4QI"
    idle = true
  }
}
resource "perimeter81_gateway"  "g2"{

  network_id = perimeter81_network.n4.id
  region_id = perimeter81_network.n4.region[0].region_id
   gateways {
      name = "perimeter81"
      idle = true
  }
  	depends_on = [
    	perimeter81_network.n4
  	]
}

data "perimeter81_networks" "all4" {
	depends_on = [
    	perimeter81_network.n4,
		perimeter81_gateway.g2
  	]
}
resource "perimeter81_ipsec_redundant" "ipsr1" {
  region_id = perimeter81_network.n4.region[0].region_id
  network_id = perimeter81_network.n4.id
  tunnel_name = "ipseed"
  tunnel1 {
      passphrase = "aXvgHEYt"
      p81_gwinternal_ip = "169.254.100.19"
      remote_gwinternal_ip = "169.254.100.5"
      remote_public_ip = "169.254.100.7"
      remote_asn = "65323"
      gateway_id = {
		for network in data.perimeter81_networks.all4.networks :
		network.id => network.regions[0].instances[0].id
		if network.id == perimeter81_network.n4.id
	  }[perimeter81_network.n4.id]
  }
  tunnel2 {
      passphrase = "Sg4gKHtT"
      p81_gwinternal_ip = "169.254.100.10"
      remote_gwinternal_ip = "169.254.100.14"
      remote_public_ip = "169.254.100.16"
      remote_asn = "65324"
      gateway_id = {
		for network in data.perimeter81_networks.all4.networks :
		network.id => network.regions[0].instances[1].id
		if network.id == perimeter81_network.n4.id
	  }[perimeter81_network.n4.id]
  }
  shared_settings {
    p81_gateway_subnets = ["0.0.0.0/0"]
    remote_gateway_subnets = ["0.0.0.0/0"]
  }
  advanced_settings {
    key_exchange = "ikev2"
    ike_life_time = "8h"
    lifetime = "1h"
    dpd_delay = "10s"
    dpd_timeout = "30s"
    phase1 {
      auth = ["3des"]
      encryption = ["sha256"]
      dh = [14]
    }
    phase2 {
      auth = ["3des"]
      encryption = ["sha256"]
      dh = [14]
    }
  }
}
  `
	return fmt.Sprintf(config, randStringBytesRmndr())
}
