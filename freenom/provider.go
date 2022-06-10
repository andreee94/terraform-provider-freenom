package freenom

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/tzwsoho/go-freenom/freenom"
)

var stderr = os.Stderr

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

type provider struct {
	configured bool
	version    string
}

// GetSchema
func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"username": {
				Type:     types.StringType,
				Optional: false,
				Computed: true,
				Required: true,
			},
			"password": {
				Type:      types.StringType,
				Optional:  false,
				Computed:  true,
				Sensitive: true,
				Required:  true,
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

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
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
}

func checkForUnknowsInConfig(config *providerData, resp *tfsdk.ConfigureProviderResponse) bool {
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

// GetResources - Defines provider resources
func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"freenom_dns_record": resourceFreenomDnsRecordType{},
	}, nil
}

// GetDataSources - Defines provider data sources
func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"freenom_dns_record":          datasourceFreenomDnsRecordType{},
		"freenom_dns_records":         datasourceFreenomDnsRecordsType{},
		"freenom_reverse_dns_records": datasourceFreenomReverseDnsRecordsType{},
	}, nil
}
