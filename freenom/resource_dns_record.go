package freenom

import (
	"context"
	"fmt"
	"log"
	"strings"
	"terraform-provider-frenom/freenom/validators"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tzwsoho/go-freenom/freenom"
)

// var _ provider.ResourceType = freenomDnsRecordResourceType{}
var _ resource.Resource = &dnsRecordResource{}
var _ resource.ResourceWithImportState = &dnsRecordResource{}

type dnsRecordResource struct {
	provider *freenomProvider
}

func NewDnsRecordResource() resource.Resource {
	return &dnsRecordResource{}
}

func (r *dnsRecordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

func (r *dnsRecordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dnsRecordResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:        types.StringType,
				Computed:    true,
				Description: "Unique identifier for this resource (<name>/<domain>)",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"domain": {
				Type:        types.StringType,
				Required:    true,
				Description: "The domain name of the record",
				Validators: []tfsdk.AttributeValidator{
					validators.IsDomain(),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"type": {
				Type: types.StringType,
				// Computed: false,
				Required:    true,
				Description: "The DNS type of the record",
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.OneOf(
						"A", "AAAA", "CNAME", "LOC", "MX", "NAPTR", "RP", "TXT",
					),
					// validators.StringIn{ValidValues: []string{"A", "AAAA", "CNAME", "LOC", "MX", "NAPTR", "RP", "TXT"}, IgnoreCase: false},
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"name": {
				Type:        types.StringType,
				Required:    true,
				Description: "The name of the record (Subdomain)",
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"value": {
				Type:        types.StringType,
				Required:    true,
				Description: "The value of the record (Ex. Ip Address)",
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"priority": {
				Type:        types.Int64Type,
				Required:    true,
				Description: "The priority of the record",
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"ttl": {
				Type:        types.Int64Type,
				Required:    true,
				Description: "The TTL of the record",
				Validators: []tfsdk.AttributeValidator{
					int64validator.AtLeast(1),
					// validators.IntGreaterThan(1),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"fqdn": {
				Type:        types.StringType,
				Computed:    true,
				Description: "The fully qualified domain name of the record (<name>.<domain>)",
			},
		},
	}, nil
}

// Create a new resource
func (r *dnsRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if !r.provider.configured {
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
	plan.FQDN = types.String{Value: computeFQDN(plan.Domain.Value, plan.Name.Value)}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *dnsRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if !r.provider.configured {
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

	log.Println("[INFO] Reading record ", state.ID.Value, domain, name)

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
func (r *dnsRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if !r.provider.configured {
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
func (r *dnsRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if !r.provider.configured {
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

	log.Println("[INFO] Reading record ", state.ID.Value, domain, name)

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
func (r *dnsRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
