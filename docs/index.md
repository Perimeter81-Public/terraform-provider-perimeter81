---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "perimeter81 Provider"
subcategory: ""
description: |-
  
---

# perimeter81 Provider

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `api_key` (String, Sensitive) The access key for API operations. You can retrieve this from the Perimeter81 Admin Console. Alternatively this can be specified using the PERIMETER81_API_KEY environment variable.
- `base_url` (String) The base url for the rest api.

### Example

```terraform
  provider "perimeter81" {
    api_key = "abcdedfgh123456789"
    base_url = "https://api.perimeter81.com/api/rest"
  }
```

