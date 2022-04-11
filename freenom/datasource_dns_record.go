package freenom

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	gofreenom "github.com/tzwsoho/go-freenom/freenom"
)

type datasourceFreenomDnsRecordType struct{}

func (c datasourceFreenomDnsRecordType) GetSchema(_ context.Context) (tfsdk.Schema,
	diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"domain": {
				Type: types.StringType,
				// Computed: false,
				Required: true,
			},
			"type": {
				Type:     types.StringType,
				Computed: true,
				Required: false,
			},
			"name": {
				Type: types.StringType,
				// Computed: false,
				Required: true,
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
		},
	}, nil
}

func (c datasourceFreenomDnsRecordType) NewDataSource(_ context.Context,
	p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return datasourceFreenomDnsRecord{
		p: *(p.(*provider)),
	}, nil
}

type datasourceFreenomDnsRecord struct {
	p provider
}

func (r datasourceFreenomDnsRecord) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	if !r.p.configured {
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

	diags = resp.State.Set(ctx, &datasourceRecord)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
