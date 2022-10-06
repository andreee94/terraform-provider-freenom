package freenom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDnsRecordResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDnsRecordResourceConfig("one", "10.10.10.10"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freenom_dns_record.test", "domain", "terraform-provider-freenom.tk"),
					resource.TestCheckResourceAttr("freenom_dns_record.test", "type", "A"),
					resource.TestCheckResourceAttr("freenom_dns_record.test", "name", "one"),
					resource.TestCheckResourceAttr("freenom_dns_record.test", "value", "10.10.10.10"),
					resource.TestCheckResourceAttr("freenom_dns_record.test", "ttl", "3600"),
					resource.TestCheckResourceAttr("freenom_dns_record.test", "priority", "0"),
				),
			},
			// Update and Read testing
			{
				Config: testAccDnsRecordResourceConfig("one", "20.20.20.20"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("freenom_dns_record.test", "domain", "terraform-provider-freenom.tk"),
					resource.TestCheckResourceAttr("freenom_dns_record.test", "type", "A"),
					resource.TestCheckResourceAttr("freenom_dns_record.test", "name", "one"),
					resource.TestCheckResourceAttr("freenom_dns_record.test", "value", "20.20.20.20"),
					resource.TestCheckResourceAttr("freenom_dns_record.test", "ttl", "3600"),
					resource.TestCheckResourceAttr("freenom_dns_record.test", "priority", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDnsRecordResourceConfig(subdomain, ip string) string {
	return fmt.Sprintf(`
provider "freenom" {}
	  
resource "freenom_dns_record" "test" { 
    domain = "terraform-provider-freenom.tk"
    type = "A"
    name = "%s"
    value = "%s"
    ttl = 3600
    priority = 0
}
`, subdomain, ip)
}
