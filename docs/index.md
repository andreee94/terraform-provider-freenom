---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "freenom Provider"
subcategory: ""
description: |-
  
---

# freenom Provider

## Example 


```hcl
terraform {
  required_providers {
    freenom = {
        source = "andreee94/freenom"
        version = "~> 0.2.1"
    }
  }
}

provider "freenom" {
  username = "<freenom-email>"
  password = "<freenom-password>"
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `username` (String)
- `password` (String, Sensitive)