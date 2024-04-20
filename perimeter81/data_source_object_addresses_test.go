package perimeter81

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceObjectAddresses_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectAddressesConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectAddressesExists("data.perimeter81_object_addresses.all"),
				),
			},
		},
	})
}

func testAccCheckObjectAddressesExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Data source not Found: %s", n)
		}
		return nil
	}
}

func testAccObjectAddressesConfig() string {
	return `
	data "perimeter81_object_addresses" "all" {}
  `
}
