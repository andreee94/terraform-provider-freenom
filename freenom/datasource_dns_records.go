package freenom

import (
	"context"
	"fmt"
	"log"
	"terraform-provider-frenom/freenom/validators"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &dnsRecordListDataSource{}

type dnsRecordListDataSource struct {
	provider *freenomProvider
}

func NewDnsRecordListDataSource() datasource.DataSource {
	return &dnsRecordListDataSource{}
}

func (d *dnsRecordListDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_records" // TODO rename to _dns_record_list
}

func (r *dnsRecordListDataSource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	provider, ok := req.ProviderData.(*freenomProvider)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *freenomProvider, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	if !provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"Expected a configured provider but it wasn't. Please report this issue to the provider developers.",
		)

		return
	}

	r.provider = provider
}

func (r *dnsRecordListDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"domain": {
				Type:        types.StringType,
				Computed:    false,
				Required:    true,
				Description: "The domain name of the record",
				Validators: []tfsdk.AttributeValidator{
					validators.IsDomain(),
				},
			},
			"records": {
				// When Computed is true, the provider will set value --
				// the user cannot define the value
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Type:        types.StringType,
						Computed:    true,
						Required:    false,
						Description: "Unique identifier for this resource (<name>/<domain>)",
					},
					"domain": {
						Type:        types.StringType,
						Computed:    true,
						Required:    false,
						Description: "The domain name of the record",
					},
					"type": {
						Type:        types.StringType,
						Computed:    true,
						Required:    false,
						Description: "The DNS type of the record",
					},
					"name": {
						Type:        types.StringType,
						Computed:    true,
						Required:    false,
						Description: "The name of the record (Subdomain)",
					},
					"value": {
						Type:        types.StringType,
						Computed:    true,
						Required:    false,
						Description: "The value of the record (Ex. Ip Address)",
					},
					"priority": {
						Type:        types.Int64Type,
						Computed:    true,
						Required:    false,
						Description: "The priority of the record",
					},
					"ttl": {
						Type:        types.Int64Type,
						Computed:    true,
						Required:    false,
						Description: "The TTL of the record",
					},
					"fqdn": {
						Type:        types.StringType,
						Computed:    true,
						Required:    false,
						Description: "The fully qualified domain name of the record (<name>.<domain>)",
					},
				}),
			},
		},
	}, nil
}

func (d *dnsRecordListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if !d.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var resourceState struct {
		Domain  string             `tfsdk:"domain"`
		Records []FreenomDnsRecord `tfsdk:"records"`
	}

	diags := req.Config.Get(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	freenomRecords, err := getAllRecordsByDomainName(resourceState.Domain, &resp.Diagnostics)

	if err != nil {
		return
	}

	log.Printf("[INFO] Found %d records", len(freenomRecords))

	for _, freenomRecord := range freenomRecords {
		var datasourceRecord FreenomDnsRecord
		datasourceRecord.ID = types.String{Value: computeID(resourceState.Domain, freenomRecord.Name)}
		datasourceRecord.Domain = types.String{Value: resourceState.Domain}
		datasourceRecord.Type = types.String{Value: freenomRecord.Type}
		datasourceRecord.Name = types.String{Value: freenomRecord.Name}
		datasourceRecord.Value = types.String{Value: freenomRecord.Value}
		datasourceRecord.Priority = types.Int64{Value: int64(freenomRecord.Priority)}
		datasourceRecord.TTL = types.Int64{Value: int64(freenomRecord.TTL)}
		datasourceRecord.FQDN = types.String{Value: computeFQDN(resourceState.Domain, freenomRecord.Name)}
		resourceState.Records = append(resourceState.Records, datasourceRecord)
	}

	diags = resp.State.Set(ctx, &resourceState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
