package perimeter81

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceServers_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworksConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworksExists("data.perimeter81_networks.all"),
				),
			},
		},
	})
}

func testAccCheckNetworksExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Data source not Found: %s", n)
		}
		return nil
	}
}

func testAccNetworksConfig() string {
	return `
	data "perimeter81_networks" "all" {}
  `
}
