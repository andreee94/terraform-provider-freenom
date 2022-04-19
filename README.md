# terraform-provider-freenom
A terraform provider to configure Freenom DNS


# Example

An example of how to:
- import the freenom provider
- setup the provider
- get data from the subdomain
- create a resource

```
terraform {
  required_providers {
    freenom = {
        source = "andreee94/freenom"
        version = "~> 0.0.1"
    }
  }
}

provider "freenom" {
  username = "<freenom-email>"
  password = "<freenom-password>"
}

data "freenom_dns_record" "grafana" {
  domain = "example.com"
  name = "grafana" // subdomain
}

// take grafana.example.com and set terraform.example.com with the same ip
resource "freenom_dns_record" "test" {
  domain = "example.com"
  type = "A"
  name = "terraform" # subdomain
  value = data.freenom_dns_record.grafana.value # ip address
  ttl = 3600
  priority = 0
}
```