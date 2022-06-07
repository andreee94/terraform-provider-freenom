# terraform-provider-freenom
A terraform provider to configure Freenom DNS


# Example

An example of how to:
- import the freenom provider
- setup the provider
- get data from the subdomain
- create a resource
- get all the subdomains
- get all the subdomains with a specific value (Ex. ip address)

```hcl
terraform {
  required_providers {
    freenom = {
        source = "andreee94/freenom"
        version = "~> 0.1.0"
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

// get all the subdomains of example.com and export them as an output
data "freenom_dns_records" "example"{
  domain = "example.com"
}

output "example" {
    value = data.freenom_dns_records.example
}

// get all the subdomains of example.com with 192.168.100.100 as value and export them as an output
data "freenom_reverse_dns_records" "example"{
  domain = "example.com"
  value = "192.168.100.100"
}

output "example_reverse" {
    value = data.freenom_reverse_dns_records.example
}
```

# NOTES

When creating more than one dns record, the creation may not succeed for every one due to race conditions of the freenom website. 

The suggestion is to run `terraform` command with `-parallelism=1` to avoid this issue.

For example:

```bash
terraform apply -parallelism=1
```

```bash
terraform destroy -parallelism=1
```

Unfortunately it is not possible to set the `parallelism` flag at resource or provider level yet ([Terraform Issue](https://github.com/hashicorp/terraform/issues/14258)). 

