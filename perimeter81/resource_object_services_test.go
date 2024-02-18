package perimeter81

import (
	"context"
	"fmt"
	"testing"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccObjectServices_basic(t *testing.T) {
	t.Parallel()
	var objectServices perimeter81Sdk.ObjectsServicesResponseObj

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectServicesConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectServicesExists("perimeter81_object_services.os", &objectServices),
					testAccCheckObjectServicesAttributes(&objectServices, &testAccObjectServicesExpectedAttributes{
						Name:        "test-os",
						Description: "10.30.0.90/16",
						Protocols: []perimeter81Sdk.ObjectServiceProtocolTcpudp{
							{
								Protocol:  "tcp",
								ValueType: "single",
								Value:     []int32{22},
							},
						},
					}),
				),
			},
			{
				Config: testAccObjectServicesUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectServicesExists("perimeter81_object_services.os", &objectServices),
					testAccCheckObjectServicesAttributes(&objectServices, &testAccObjectServicesExpectedAttributes{
						Name:        "test-os-updated",
						Description: "10.30.0.91/16",
						Protocols: []perimeter81Sdk.ObjectServiceProtocolTcpudp{
							{
								Protocol:  "udp",
								ValueType: "list",
								Value:     []int32{23, 24},
							},
						},
					}),
				),
			},
		},
	})
}

func testAccCheckObjectServicesExists(n string, objectServices *perimeter81Sdk.ObjectsServicesResponseObj) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		ObjectServicesID := rs.Primary.ID
		if ObjectServicesID == "" {
			return fmt.Errorf("No ObjectServices ID is set")
		}
		conn := testAccProvider.Meta().(*perimeter81Sdk.APIClient)
		ctx := context.Background()
		objectsServices, _, err := conn.ObjectsServicesApi.GetObjectsServices(ctx)
		if err != nil {
			return fmt.Errorf("No ObjectServices found")
		}
		currentObjectServices := getCurrentObjectServicesInArray(&objectsServices, ObjectServicesID)

		*objectServices = *currentObjectServices
		return nil
	}
}

type testAccObjectServicesExpectedAttributes struct {
	Name        string
	Description string
	Protocols   []perimeter81Sdk.ObjectServiceProtocolTcpudp
}

func testAccCheckObjectServicesAttributes(objectServices *perimeter81Sdk.ObjectsServicesResponseObj, want *testAccObjectServicesExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if objectServices.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", objectServices.Name, want.Name)
		}

		if objectServices.Description != want.Description {
			return fmt.Errorf("got tags %q; want %q", objectServices.Description, want.Description)
		}

		if objectServices.Protocols[0].Protocol != want.Protocols[0].Protocol {
			return fmt.Errorf("got protocol %q; want %q", objectServices.Protocols[0].Protocol, want.Protocols[0].Protocol)
		}

		if objectServices.Protocols[0].ValueType != want.Protocols[0].ValueType {
			return fmt.Errorf("got value type %q; want %q", objectServices.Protocols[0].ValueType, want.Protocols[0].ValueType)
		}
		if !testComparableArraiesEq(objectServices.Protocols[0].Value, want.Protocols[0].Value) {
			return fmt.Errorf("got value %q; want %q", objectServices.Protocols[0].Value, want.Protocols[0].Value)
		}

		return nil
	}
}

func testAccObjectServicesConfig() string {
	config := `
resource "perimeter81_object_services" "os" {
  name = "test-os"
  description = "10.30.0.90/16"

  protocols {
    protocol = "tcp"
    value_type = "single"
    value = [22]
  }
}
  `
	return config
}

func testAccObjectServicesUpdateConfig() string {
	config := `
resource "perimeter81_object_services" "os" {
  name = "test-os-updated"
  description = "10.30.0.91/16"

  protocols {
    protocol = "udp"
    value_type = "list"
    value = [23, 24]
  }
}
  `
	return config
}
