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
	gofreenom "github.com/tzwsoho/go-freenom/freenom"
)

var _ datasource.DataSource = &dnsRecordDataSource{}

type dnsRecordDataSource struct {
	provider *freenomProvider
}

func NewDnsRecordDataSource() datasource.DataSource {
	return &dnsRecordDataSource{}
}

func (d *dnsRecordDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

func (r *dnsRecordDataSource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dnsRecordDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Unique identifier for this resource (<name>/<domain>)",
			},
			"domain": {
				Type: types.StringType,
				// Computed: false,
				Required:    true,
				Description: "The domain name of the record",
				Validators: []tfsdk.AttributeValidator{
					validators.IsDomain(),
				},
			},
			"type": {
				Type:        types.StringType,
				Computed:    true,
				Required:    false,
				Description: "The DNS type of the record",
			},
			"name": {
				Type: types.StringType,
				// Computed: false,
				Required:    true,
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
		},
	}, nil
}

func (d *dnsRecordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if !d.provider.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var datasourceRecord FreenomDnsRecord
	var freenomRecord *gofreenom.DomainRecord

	diags := req.Config.Get(ctx, &datasourceRecord)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	log.Println("[INFO] Reading record", datasourceRecord.Domain, datasourceRecord.Name)

	freenomRecord, err := getRecordByName(datasourceRecord.Domain.Value, datasourceRecord.Name.Value, &resp.Diagnostics)

	if err != nil {
		return
	}

	datasourceRecord.ID = types.String{Value: computeID(datasourceRecord.Domain.Value, datasourceRecord.Name.Value)}
	datasourceRecord.Value = types.String{Value: freenomRecord.Value}
	datasourceRecord.Type = types.String{Value: freenomRecord.Type}
	datasourceRecord.TTL = types.Int64{Value: int64(freenomRecord.TTL)}
	datasourceRecord.Priority = types.Int64{Value: int64(freenomRecord.Priority)}
	datasourceRecord.FQDN = types.String{Value: computeFQDN(datasourceRecord.Domain.Value, datasourceRecord.Name.Value)}

	diags = resp.State.Set(ctx, &datasourceRecord)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
