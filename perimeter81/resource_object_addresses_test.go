package perimeter81

import (
	"context"
	"fmt"
	"testing"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccObjectAddresses_basic(t *testing.T) {
	t.Parallel()
	var objectAddress perimeter81Sdk.ObjectsAddressObj

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectAddressConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectAddressExists("perimeter81_object_addresses.os", &objectAddress),
					testAccCheckObjectAddressesAttributes(&objectAddress, &testAccObjectAddressExpectedAttributes{
						Name:        "test-os",
						Description: "10.30.0.90/16",
						ValueType:   "single",
						Value:       []string{"193.168.3.1"},
					}),
				),
			},
			{
				Config: testAccObjectAddressUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectAddressExists("perimeter81_object_services.os", &objectAddress),
					testAccCheckObjectAddressesAttributes(&objectAddress, &testAccObjectAddressExpectedAttributes{
						Name:        "test-os-updated",
						Description: "10.30.0.91/16",
						ValueType:   "list",
						Value:       []string{"193.168.3.2"},
					}),
				),
			},
		},
	})
}

func testAccCheckObjectAddressExists(n string, objectAddress *perimeter81Sdk.ObjectsAddressObj) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		ObjectAddressID := rs.Primary.ID
		if ObjectAddressID == "" {
			return fmt.Errorf("No ObjectAddresses ID is set")
		}
		conn := testAccProvider.Meta().(*perimeter81Sdk.APIClient)
		ctx := context.Background()
		objectsAddresses, _, err := conn.ObjectsAddressesApi.GetObjectsAddresses(ctx)
		if err != nil {
			return fmt.Errorf("No ObjectServices found")
		}
		currentObjectAddress := getCurrentObjectAddressesInArray(&objectsAddresses, ObjectAddressID)

		*objectAddress = *currentObjectAddress
		return nil
	}
}

type testAccObjectAddressExpectedAttributes struct {
	Name        string
	Description string
	ValueType   string
	Value       []string
}

func testAccCheckObjectAddressesAttributes(objectAddress *perimeter81Sdk.ObjectsAddressObj, want *testAccObjectAddressExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if objectAddress.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", objectAddress.Name, want.Name)
		}

		if objectAddress.Description != want.Description {
			return fmt.Errorf("got description %q; want %q", objectAddress.Description, want.Description)
		}

		if objectAddress.ValueType != want.ValueType {
			return fmt.Errorf("got value type %q; want %q", objectAddress.ValueType, want.ValueType)
		}

		if !testComparableArraiesEq(objectAddress.Value, want.Value) {
			return fmt.Errorf("got value %q; want %q", objectAddress.Value, want.Value)
		}

		return nil
	}
}

func testAccObjectAddressConfig() string {
	config := `
resource "perimeter81_object_addresses" "os" {
  name = "test-os"
  description = "10.30.0.90/16"
  value_type = "single"
  value = ["193.168.3.1""]
}
  `
	return config
}

func testAccObjectAddressUpdateConfig() string {
	config := `
resource "perimeter81_object_addresses" "os" {
  name = "test-os-updated"
  description = "10.30.0.91/16"
  value_type = "list"
  value = ["193.168.3.2""]
}
  `
	return config
}
