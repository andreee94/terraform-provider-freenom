package freenom

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FreenomDnsRecord struct {
	ID       types.String `tfsdk:"id"`
	Domain   types.String `tfsdk:"domain"`
	Type     types.String `tfsdk:"type"`
	Name     types.String `tfsdk:"name"`
	Value    types.String `tfsdk:"value"`
	Priority types.Int64  `tfsdk:"priority"`
	TTL      types.Int64  `tfsdk:"ttl"`
}
