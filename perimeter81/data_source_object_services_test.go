package perimeter81

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceObjectServices_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectServicessConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectServicessExists("data.perimeter81_object_services.all"),
				),
			},
		},
	})
}

func testAccCheckObjectServicessExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Data source not Found: %s", n)
		}
		return nil
	}
}

func testAccObjectServicessConfig() string {
	return `
	data "perimeter81_object_services" "all" {}
  `
}
