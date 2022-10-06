package main

import (
	"context"
	"flag"
	"log"
	"terraform-provider-frenom/freenom"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// func main() {
// 	tfsdk.Serve(context.Background(), freenom.New, tfsdk.ServeOpts{
// 		Name: "freenom",
// 	})
// }

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "0.3.0"

	// goreleaser can also pass the specific commit if you want
	// commit  string = ""
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/andreee94/freenom",
		// Address: "local/tr/freenom",
		Debug: debug,
	}

	err := providerserver.Serve(context.Background(), freenom.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
