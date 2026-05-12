package perimeter81

import (
	"context"
	"fmt"
	"testing"

	perimeter81Sdk "github.com/Perimeter81-Public/perimeter-81-client-sdk/v2"

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
					testAccCheckObjectServicesExists("checkpointsase_object_services.os", &objectServices),
					testAccCheckObjectServicesAttributes(&objectServices, &testAccObjectServicesExpectedAttributes{
						Name:        "test-os",
						Description: "10.30.0.90/16",
						ValueType:   "single",
						Value:       []int32{22},
					}),
				),
			},
			{
				Config: testAccObjectServicesUpdateConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObjectServicesExists("checkpointsase_object_services.os", &objectServices),
					testAccCheckObjectServicesAttributes(&objectServices, &testAccObjectServicesExpectedAttributes{
						Name:        "test-os-updated",
						Description: "10.30.0.91/16",
						ValueType:   "list",
						Value:       []int32{23, 24},
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
		objectsServices, _, err := conn.ObjectsServicesAPI.GetObjectsServices(ctx).Execute()
		if err != nil {
			return fmt.Errorf("No ObjectServices found")
		}
		// Match by name since the list API does not return IDs
		name := rs.Primary.Attributes["name"]
		currentObjectServices := getCurrentObjectServicesInArray(objectsServices, name)
		if currentObjectServices == nil {
			return fmt.Errorf("ObjectServices with name %q not found", name)
		}

		*objectServices = *currentObjectServices
		return nil
	}
}

type testAccObjectServicesExpectedAttributes struct {
	Name        string
	Description string
	ValueType   string
	Value       []int32
}

func testAccCheckObjectServicesAttributes(objectServices *perimeter81Sdk.ObjectsServicesResponseObj, want *testAccObjectServicesExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if objectServices.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", objectServices.Name, want.Name)
		}

		if objectServices.GetDescription() != want.Description {
			return fmt.Errorf("got description %q; want %q", objectServices.GetDescription(), want.Description)
		}

		if len(objectServices.Protocols) == 0 {
			return fmt.Errorf("got no protocols; want at least one")
		}

		// Extract value_type and value from the first protocol's union type
		proto := objectServices.Protocols[0]
		var gotValueType string
		var gotValue []int32
		if proto.ObjectServiceProtocolTCPUDP != nil {
			tcpudp := proto.ObjectServiceProtocolTCPUDP
			if tcpudp.ObjectServiceProtocolList != nil {
				gotValueType = tcpudp.ObjectServiceProtocolList.ValueType
				gotValue = tcpudp.ObjectServiceProtocolList.Value
			} else if tcpudp.ObjectServiceProtocolRange != nil {
				gotValueType = tcpudp.ObjectServiceProtocolRange.ValueType
				gotValue = tcpudp.ObjectServiceProtocolRange.Value
			} else if tcpudp.ObjectServiceProtocolSingle != nil {
				gotValueType = tcpudp.ObjectServiceProtocolSingle.ValueType
				gotValue = tcpudp.ObjectServiceProtocolSingle.Value
			}
		}

		if gotValueType != want.ValueType {
			return fmt.Errorf("got value_type %q; want %q", gotValueType, want.ValueType)
		}
		if !testComparableArraiesEq(gotValue, want.Value) {
			return fmt.Errorf("got value %v; want %v", gotValue, want.Value)
		}

		return nil
	}
}

func testAccObjectServicesConfig() string {
	config := `
resource "checkpointsase_object_services" "os" {
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
resource "checkpointsase_object_services" "os" {
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
