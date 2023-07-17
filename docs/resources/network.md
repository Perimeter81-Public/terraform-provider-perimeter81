---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "perimeter81_network Resource - terraform-provider-perimeter81"
subcategory: ""
description: |-
  
---

# perimeter81_network (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `region` (Block List, Min: 1) (see [below for nested schema](#nestedblock--region))

### Optional

- `network` (Block List) (see [below for nested schema](#nestedblock--network))

### Read-Only

- `id` (String) The ID of this resource.

### Example

```terraform
 resource "perimeter81_network" "n1" {
   network {
     name = "network-test",
     tags = ["test"]
   }
   region {
     cpregion_id = "v2cRwzGRua"
     instance_count = 1
     idle = true
   }
   region {
     cpregion_id = "F2w4QTggWt"
     instance_count = 1
     idle = true
   }
 }
```

<a id="nestedblock--region"></a>
### Nested Schema for `region`

Required:

- `cpregion_id` (String) General ID that uniquely identifies the region across all regions can be found from the regions datasource.
- `instance_count` (Number) Number of instances that will be created for a specified region inside the network.

Optional:

- `idle` (Boolean) either the gateway is idel or not


<a id="nestedblock--network"></a>
### Nested Schema for `network`

Required:

- `name` (String) The name of the Resource

Optional:

- `subnet` (String) Subnet to associate
- `tags` (List of String) List of tags for the network Resource