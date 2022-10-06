---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "freenom_dns_records Data Source - terraform-provider-freenom"
subcategory: ""
description: |-
  
---

# freenom_dns_records (Data Source)

## Example

```hcl
data "freenom_dns_records" "all" {
  domain = "example.com"
}

output "test1" {
    value = data.freenom_dns_record // extract all the records for example.com domain
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `domain` (String) The domain name of the record

### Read-Only

- `records` (Attributes List) (see [below for nested schema](#nestedatt--records))

<a id="nestedatt--records"></a>
### Nested Schema for `records`

Read-Only:

- `domain` (String) The domain name of the record
- `fqdn` (String) The fully qualified domain name of the record (<name>.<domain>)
- `id` (String) Unique identifier for this resource (<name>/<domain>)
- `name` (String) The name of the record (Subdomain)
- `priority` (Number) The priority of the record
- `ttl` (Number) The TTL of the record
- `type` (String) The DNS type of the record
- `value` (String) The value of the record (Ex. Ip Address)

