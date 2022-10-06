package freenom

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tzwsoho/go-freenom/freenom"
)

// var stderr = os.Stderr

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &freenomProvider{
			version: version,
		}
	}
}

type freenomProvider struct {
	configured bool
	version    string
}

func (p *freenomProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "freenom"
	resp.Version = p.version
}

// GetSchema
func (p *freenomProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"username": {
				Type:     types.StringType,
				Optional: true,
				Required: false,
				Computed: false,
			},
			"password": {
				Type:      types.StringType,
				Optional:  true,
				Required:  false,
				Computed:  false,
				Sensitive: true,
			},
		},
	}, nil
}

// Provider schema struct
type providerData struct {
	// Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *freenomProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var username string
	var password string

	if !checkForUnknowsInConfig(&config, resp) {
		return
	}

	if config.Username.Null {
		username = os.Getenv("FREENOM_USERNAME")
	} else {
		username = config.Username.Value
	}

	if config.Password.Null {
		password = os.Getenv("FREENOM_PASSWORD")
	} else {
		password = config.Password.Value
	}

	if username == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find username",
			"Username cannot be an empty string",
		)
		return
	}

	if password == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find password",
			"password cannot be an empty string",
		)
		return
	}

	// Login to freenom
	err := freenom.Login(username, password)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Freenom login failed. Error: "+err.Error(),
		)
		return
	}

	p.configured = true

	resp.DataSourceData = p
	resp.ResourceData = p
}

func checkForUnknowsInConfig(config *providerData, resp *provider.ConfigureResponse) bool {

	if config.Username.Unknown {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as username",
		)
		return false
	}

	if config.Password.Unknown {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as password",
		)
		return false
	}
	return true
}

func (p *freenomProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDnsRecordResource,
	}
}

func (p *freenomProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDnsRecordDataSource,
		NewDnsRecordListDataSource,
		NewReverseDnsRecordListDataSource,
	}
}
