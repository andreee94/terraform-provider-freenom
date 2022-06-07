package freenom

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type datasourceFreenomDnsRecordsType struct{}

func (c datasourceFreenomDnsRecordsType) GetSchema(_ context.Context) (tfsdk.Schema,
	diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"domain": {
				Type:     types.StringType,
				Computed: false,
				Required: true,
			},
			"records": {
				// When Computed is true, the provider will set value --
				// the user cannot define the value
				Computed: true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Type:     types.StringType,
						Computed: true,
						Required: false,
					},
					"domain": {
						Type:     types.StringType,
						Computed: true,
						Required: false,
					},
					"type": {
						Type:     types.StringType,
						Computed: true,
						Required: false,
					},
					"name": {
						Type:     types.StringType,
						Computed: true,
						Required: false,
					},
					"value": {
						Type:     types.StringType,
						Computed: true,
						Required: false,
					},
					"priority": {
						Type:     types.Int64Type,
						Computed: true,
						Required: false,
					},
					"ttl": {
						Type:     types.Int64Type,
						Computed: true,
						Required: false,
					},
					"fqdn": {
						Type:     types.StringType,
						Computed: true,
						Required: false,
					},
				}, tfsdk.ListNestedAttributesOptions{}),
			},
		},
	}, nil
}

func (c datasourceFreenomDnsRecordsType) NewDataSource(_ context.Context,
	p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return datasourceFreenomDnsRecords{
		p: *(p.(*provider)),
	}, nil
}

type datasourceFreenomDnsRecords struct {
	p provider
}

func (r datasourceFreenomDnsRecords) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	if !r.p.configured {
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
