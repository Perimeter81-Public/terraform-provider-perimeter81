package perimeter81

import (
	"context"
	"fmt"
	"testing"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var randNameIpsecSignle string = randStringBytesRmndr()

func TestAccIpsecSingle_basic(t *testing.T) {
	t.Parallel()
	var tunnel perimeter81Sdk.IpSecSingleTunnel
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIpsecSingleConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIpsecSingleExists("perimeter81_ipsec_single.ipss1", &tunnel),
					testAccCheckIpsecSingleAttributes(&tunnel, &testAccIpSecSingleExpectedAttributes{
						P81GatewaySubnets:    []string{"0.0.0.0/0"},
						RemoteGatewaySubnets: []string{"0.0.0.0/0"},
						KeyExchange:          "ikev1",
						IkeLifeTime:          "9h",
						Lifetime:             "2h",
						DpdDelay:             "20s",
						DpdTimeout:           "40s",
						Passphrase:           "tnEgVbTJE23",
						RemotePublicIP:       "198.51.100.41",
						Phase1: perimeter81Sdk.IpSecPhase{
							Auth:       []string{"3des"},
							Encryption: []string{"sha256"},
							Dh:         []int32{14},
						},
						Phase2: perimeter81Sdk.IpSecPhase{
							Auth:       []string{"3des"},
							Encryption: []string{"sha256"},
							Dh:         []int32{14},
						},
					}),
				),
			},
			{
				Config: testAccIpsecSingleUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIpsecSingleExists("perimeter81_ipsec_single.ipss1", &tunnel),
					testAccCheckIpsecSingleAttributes(&tunnel, &testAccIpSecSingleExpectedAttributes{
						P81GatewaySubnets:    []string{"0.0.0.0/0"},
						RemoteGatewaySubnets: []string{"0.0.0.0/0"},
						KeyExchange:          "ikev2",
						IkeLifeTime:          "10h",
						Lifetime:             "3h",
						DpdDelay:             "30s",
						DpdTimeout:           "50s",
						Passphrase:           "tnEgVbTJE23123",
						RemotePublicIP:       "198.51.100.42",
						Phase1: perimeter81Sdk.IpSecPhase{
							Auth:       []string{"blowfish256"},
							Encryption: []string{"md5"},
							Dh:         []int32{19},
						},
						Phase2: perimeter81Sdk.IpSecPhase{
							Auth:       []string{"blowfish256"},
							Encryption: []string{"md5"},
							Dh:         []int32{19},
						},
					}),
				),
			},
		},
	})
}

func testAccCheckIpsecSingleExists(n string, tunnel *perimeter81Sdk.IpSecSingleTunnel) resource.TestCheckFunc {
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
		gotIpsecSingle, _, err := conn.IPSecSingleApi.GetIPSecSingleTunnel(ctx, networkId, tunnelId)
		if err != nil {
			return err
		}

		*tunnel = gotIpsecSingle
		return nil
	}
}

func testAccCheckIpsecSingleAttributes(tunnel *perimeter81Sdk.IpSecSingleTunnel, want *testAccIpSecSingleExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !testComparableArraiesEq(tunnel.P81GatewaySubnets, want.P81GatewaySubnets) {
			return fmt.Errorf("got p81 gateway subnets %q; want %q", tunnel.P81GatewaySubnets, want.P81GatewaySubnets)
		}

		if !testComparableArraiesEq(tunnel.RemoteGatewaySubnets, want.RemoteGatewaySubnets) {
			return fmt.Errorf("got remote gateway subnets %q; want %q", tunnel.RemoteGatewaySubnets, want.RemoteGatewaySubnets)
		}
		if tunnel.KeyExchange != want.KeyExchange {
			return fmt.Errorf("got key exchange %q; want %q", tunnel.KeyExchange, want.KeyExchange)
		}
		if tunnel.IkeLifeTime != want.IkeLifeTime {
			return fmt.Errorf("got ike life time %q; want %q", tunnel.IkeLifeTime, want.IkeLifeTime)
		}
		if tunnel.DpdDelay != want.DpdDelay {
			return fmt.Errorf("got dpd delay %q; want %q", tunnel.DpdDelay, want.DpdDelay)
		}
		if tunnel.DpdTimeout != want.DpdTimeout {
			return fmt.Errorf("got dpd timeout %q; want %q", tunnel.DpdTimeout, want.DpdTimeout)
		}
		if tunnel.Passphrase != want.Passphrase {
			return fmt.Errorf("got passphrase %q; want %q", tunnel.Passphrase, want.Passphrase)
		}
		if tunnel.RemotePublicIP != want.RemotePublicIP {
			return fmt.Errorf("got remote public ip %q; want %q", tunnel.RemotePublicIP, want.RemotePublicIP)
		}
		if !testComparableArraiesEq(tunnel.Phase1.Auth, want.Phase1.Auth) {
			return fmt.Errorf("got phase1 auth %q; want %q", tunnel.Phase1.Auth, want.Phase1.Auth)
		}
		if !testComparableArraiesEq(tunnel.Phase1.Encryption, want.Phase1.Encryption) {
			return fmt.Errorf("got phase1 encryption %q; want %q", tunnel.Phase1.Encryption, want.Phase1.Encryption)
		}
		if !testComparableArraiesEq(tunnel.Phase1.Dh, want.Phase1.Dh) {
			return fmt.Errorf("got phase1 encryption %q; want %q", tunnel.Phase1.Dh, want.Phase1.Dh)
		}
		if !testComparableArraiesEq(tunnel.Phase2.Auth, want.Phase2.Auth) {
			return fmt.Errorf("got phase2 auth %q; want %q", tunnel.Phase2.Auth, want.Phase2.Auth)
		}
		if !testComparableArraiesEq(tunnel.Phase2.Encryption, want.Phase2.Encryption) {
			return fmt.Errorf("got phase2 encryption %q; want %q", tunnel.Phase2.Encryption, want.Phase2.Encryption)
		}
		if !testComparableArraiesEq(tunnel.Phase2.Dh, want.Phase2.Dh) {
			return fmt.Errorf("got phase2 dh %q; want %q", tunnel.Phase2.Dh, want.Phase2.Dh)
		}

		return nil
	}
}

type testAccIpSecSingleExpectedAttributes struct {
	P81GatewaySubnets    []string
	RemoteGatewaySubnets []string
	KeyExchange          string
	IkeLifeTime          string
	Lifetime             string
	DpdDelay             string
	DpdTimeout           string
	Passphrase           string
	RemotePublicIP       string
	Phase1               perimeter81Sdk.IpSecPhase
	Phase2               perimeter81Sdk.IpSecPhase
}

func testAccIpsecSingleConfig() string {
	config := `
resource "perimeter81_network" "n3" {
  network {
    name = "%s"
    tags = ["test"]
  }
  region {
    cpregion_id = "Xv3BREC4QI"
    idle = true
  }
}

data "perimeter81_networks" "all3" {
	depends_on = [
    	perimeter81_network.n3
  	]
}

resource "perimeter81_ipsec_single" "ipss1" { 
  network_id = perimeter81_network.n3.id
  region_id = {
    for network in data.perimeter81_networks.all3.networks :
    network.id => network.regions[0].id
    if network.id == perimeter81_network.n3.id
  }[perimeter81_network.n3.id]
  gateway_id = {
    for network in data.perimeter81_networks.all3.networks :
    network.id => network.regions[0].instances[0].id
    if network.id == perimeter81_network.n3.id
  }[perimeter81_network.n3.id]
  tunnel_name = "IpSecSingle"
  p81_gateway_subnets = ["0.0.0.0/0"]
  remote_gateway_subnets = ["0.0.0.0/0"]
  key_exchange = "ikev1"
  ike_life_time = "9h"
  lifetime = "2h"
  dpd_delay = "20s"
  dpd_timeout = "40s"
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
  passphrase = "tnEgVbTJE23"
  remote_public_ip = "198.51.100.41"
}
  `
	return fmt.Sprintf(config, randNameIpsecSignle)
}

func testAccIpsecSingleUpdateConfig() string {
	config := `
resource "perimeter81_network" "n3" {
  network {
    name = "%s"
    tags = ["test"]
  }
  region {
    cpregion_id = "Xv3BREC4QI"
    idle = true
  }
}

data "perimeter81_networks" "all3" {
	depends_on = [
    	perimeter81_network.n3
  	]
}

resource "perimeter81_ipsec_single" "ipss1" { 
  network_id = perimeter81_network.n3.id
  region_id = {
    for network in data.perimeter81_networks.all3.networks :
    network.id => network.regions[0].id
    if network.id == perimeter81_network.n3.id
  }[perimeter81_network.n3.id]
  gateway_id = {
    for network in data.perimeter81_networks.all3.networks :
    network.id => network.regions[0].instances[0].id
    if network.id == perimeter81_network.n3.id
  }[perimeter81_network.n3.id]
  tunnel_name = "IpSecSingle"
  p81_gateway_subnets = ["0.0.0.0/0"]
  remote_gateway_subnets = ["0.0.0.0/0"]
  key_exchange = "ikev2"
  ike_life_time = "10h"
  lifetime = "3h"
  dpd_delay = "30s"
  dpd_timeout = "50s"
  phase1 {
    auth = ["blowfish256"]
    encryption = ["md5"]
    dh = [19]
  }
  phase2 {
    auth = ["blowfish256"]
    encryption = ["md5"]
    dh = [19]
  }
  passphrase = "tnEgVbTJE23123"
  remote_public_ip = "198.51.100.42"
}
  `
	return fmt.Sprintf(config, randNameIpsecSignle)
}
