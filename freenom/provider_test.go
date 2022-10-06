package freenom

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"freenom": providerserver.NewProtocol6WithError(New("0.3.0")()),
}

func testAccPreCheck(t *testing.T) {
	checkEnvNotNull(t, "FREENOM_USERNAME")
	checkEnvNotNull(t, "FREENOM_PASSWORD")
}

func checkEnvNotNull(t *testing.T, env string) {
	if v := os.Getenv(env); v == "" {
		t.Fatalf("%s must be set for acceptance tests", env)
	}
}
