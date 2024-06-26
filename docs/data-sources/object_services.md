---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "perimeter81_object_services Data Source - terraform-provider-perimeter81"
subcategory: ""
description: |-
  
---

# perimeter81_object_services (Data Source)

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `object_services` (List of Object) (see [below for nested schema](#nestedatt--object_services))

### Example

```terraform
 data "perimeter81_object_services" "all" {}
```

<a id="nestedatt--object_services"></a>
### Nested Schema for `object_services`

Read-Only:

- `id` (String) The ID of the datasource
- `name` (String) The name of the object services resource
- `protocols` (Block List) (see [below for nested schema](#nestedblock--protocols))
- `description` (String) The description of the added object services resource

<a id="nestedblock--protocols"></a>
### Nested Schema for `protocols`

Required:

- `protocol` (String) The protocol name
- `value_type` (String) The value type of the protocol
- `value` (List of ints) List of integer values
