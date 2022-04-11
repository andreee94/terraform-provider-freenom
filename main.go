package main

import (
	"context"
	"terraform-provider-frenom/freenom"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func main() {
	tfsdk.Serve(context.Background(), freenom.New, tfsdk.ServeOpts{
		Name: "freenom",
	})
}
