package freenom

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/tzwsoho/go-freenom/freenom"
)

type resourceFreenomDnsRecordType struct{}

// Order Resource schema
func (r resourceFreenomDnsRecordType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
				Type: types.StringType,
				// Computed: false,
				Required: true,
			},
			"name": {
				Type: types.StringType,
				// Computed: false,
				Required: true,
			},
			"value": {
				Type: types.StringType,
				// Computed: false,
				Required: true,
			},
			"priority": {
				Type: types.Int64Type,
				// Computed: false,
				Required: true,
			},
			"ttl": {
				Type: types.Int64Type,
				// Computed: false,
				Required: true,
			},
			"fqdn": {
				Type:     types.StringType,
				Computed: true,
				Required: false,
			},
		},
	}, nil
}

// New resource instance
func (r resourceFreenomDnsRecordType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceFreenomDnsRecord{
		p: *(p.(*provider)),
	}, nil
}

type resourceFreenomDnsRecord struct {
	p provider
}

// Create a new resource
func (r resourceFreenomDnsRecord) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var plan FreenomDnsRecord
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// log.Println("[INFO] Creating record", plan.Name.Value, plan.Value.Value)

	err := freenom.AddRecord(plan.Domain.Value, []freenom.DomainRecord{
		{
			Type:     plan.Type.Value,
			Name:     strings.ToLower(plan.Name.Value),
			Value:    plan.Value.Value,
			Priority: int(plan.Priority.Value),
			TTL:      int(plan.TTL.Value),
		},
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating record",
			err.Error(),
		)
		return
	}

	plan.ID = types.String{Value: computeID(plan.Domain.Value, plan.Name.Value)}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r resourceFreenomDnsRecord) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var state FreenomDnsRecord
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain, name, err := parseID(state.ID.Value)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing id"+state.ID.Value,
			err.Error(),
		)
		return
	}

	log.Println("[INFO] Reading record", state.ID.Value, domain, name)

	record, err := getRecordByName(domain, name, &resp.Diagnostics)

	if err != nil {
		return
	}

	state.ID = types.String{Value: computeID(domain, name)}
	state.Type = types.String{Value: record.Type}
	state.Name = types.String{Value: strings.ToLower(record.Name)}
	state.Value = types.String{Value: record.Value}
	state.Priority = types.Int64{Value: int64(record.Priority)}
	state.TTL = types.Int64{Value: int64(record.TTL)}
	state.FQDN = types.String{Value: computeFQDN(domain, name)}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update resource
func (r resourceFreenomDnsRecord) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var state FreenomDnsRecord
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan FreenomDnsRecord
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Domain.Value != plan.Domain.Value {
		resp.Diagnostics.AddError(
			"Domain changed",
			"Domain cannot be changed",
		)
		return
	}

	domain := plan.Domain.Value

	oldRecord := &freenom.DomainRecord{
		Type:     state.Type.Value,
		Name:     strings.ToLower(state.Name.Value),
		Value:    state.Value.Value,
		Priority: int(state.Priority.Value),
		TTL:      int(state.TTL.Value),
	}

	newRecord := &freenom.DomainRecord{
		Type:     plan.Type.Value,
		Name:     strings.ToLower(plan.Name.Value),
		Value:    plan.Value.Value,
		Priority: int(plan.Priority.Value),
		TTL:      int(plan.TTL.Value),
	}

	log.Printf("[DEBUG] Old record: %v\n", *oldRecord)
	log.Printf("[DEBUG] New record: %v\n", *newRecord)

	log.Printf("[DEBUG] Domain: %v\n", domain)

	err := freenom.ModifyRecord(domain, oldRecord, newRecord)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating record "+state.ID.Value,
			err.Error(),
		)
		return
	}
}

// Delete resource
func (r resourceFreenomDnsRecord) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError(
			"Provider not configured",
			"The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		)
		return
	}

	var state FreenomDnsRecord
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain, name, err := parseID(state.ID.Value)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing id"+state.ID.Value,
			err.Error(),
		)
		return
	}

	log.Println("[INFO] Reading record", state.ID.Value, domain, name)

	record, err := getRecordByName(domain, name, &resp.Diagnostics)

	if err != nil {
		return
	}

	err = freenom.DeleteRecord(domain, record)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting record",
			err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

// Import resource
func (r resourceFreenomDnsRecord) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the id attribute
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
